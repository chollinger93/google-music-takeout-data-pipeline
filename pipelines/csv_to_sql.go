package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/databaseio"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	_ "github.com/go-sql-driver/mysql"
)

var (
	trackCsvDir      = flag.String("track_csv_dir", "", "Directory containing all track CSVs")
	playlistsCsvDir  = flag.String("playlist_csv_dir", "", "Directory containing all playlist CSVs")
	radioCsvDir      = flag.String("radio_csv_dir", "", "Directory containing all radio CSVs")
	databaseHost     = flag.String("database_host", "localhost", "Database host")
	databaseUser     = flag.String("database_user", "", "Database user")
	databasePassword = flag.String("database_password", "", "Database password")
	database         = flag.String("database", "googleplay_raw", "Database name")
)

type Playlist struct {
	Title         string `column:"title"`
	Album         string `column:"album"`
	Artist        string `column:"artist"`
	Duration_ms   int    `column:"duration_ms"`
	Rating        int    `column:"rating"`
	Play_count    int    `column:"play_count"`
	Removed       bool   `column:"removed"`
	PlaylistIndex int    `column:"playlist_ix"`
}

func (s Playlist) getCols() []string {
	return []string{"title", "album", "artist", "duration_ms", "rating", "play_count", "removed", "playlist_ix"}
}

func (f *Playlist) ProcessElement(w string, emit func(Playlist)) {
	data, err := PrepCsv(w)
	if err != nil {
		return
	}
	duration_ms, _ := strconv.Atoi(GetOrDefault(data, 3))
	rating, _ := strconv.Atoi(GetOrDefault(data, 4))
	playCount, _ := strconv.Atoi(GetOrDefault(data, 5))
	playlist_ix, _ := strconv.Atoi(GetOrDefault(data, 7))

	s := &Playlist{
		Title:         GetOrDefault(data, 0),
		Album:         GetOrDefault(data, 1),
		Artist:        GetOrDefault(data, 2),
		Duration_ms:   duration_ms,
		Rating:        rating,
		Play_count:    playCount,
		Removed:       ParseRemoved(data, 6),
		PlaylistIndex: playlist_ix,
	}

	fmt.Printf("Playlist: %v\n", s)
	emit(*s)
}

type Song struct {
	Title       string `column:"title"`
	Album       string `column:"album"`
	Artist      string `column:"artist"`
	Duration_ms int    `column:"duration_ms"`
	Rating      int    `column:"rating"`
	Play_count  int    `column:"play_count"`
	Removed     bool   `column:"removed"`
	FileName    string `column:"file_name"`
}

func (s Song) getCols() []string {
	return []string{"title", "album", "artist", "duration_ms", "rating", "play_count", "removed", "file_name"}
}

func (f *Song) ProcessElement(w string, emit func(Song)) {
	data, err := PrepCsv(w)
	if err != nil {
		return
	}
	duration_ms, _ := strconv.Atoi(GetOrDefault(data, 3))
	rating, _ := strconv.Atoi(GetOrDefault(data, 4))
	playCount, _ := strconv.Atoi(GetOrDefault(data, 5))

	s := &Song{
		Title:       GetOrDefault(data, 0),
		Album:       GetOrDefault(data, 1),
		Artist:      GetOrDefault(data, 2),
		Duration_ms: duration_ms,
		Rating:      rating,
		Play_count:  playCount,
		Removed:     ParseRemoved(data, 6),
		FileName:    GetOrDefault(data, 7),
	}

	fmt.Printf("Song: %v\n", s)
	emit(*s)
}

type Radio struct {
	Title                   string `column:"title"`
	Artist                  string `column:"artist"`
	Description             string `column:"description"`
	Removed                 bool   `column:"removed"`
	Similar_stations        string `column:"similar_stations"`
	Artists_on_this_station string `column:"artists_on_this_station"`
}

func (s Radio) getCols() []string {
	return []string{"title", "artist", "description", "removed", "similar_stations", "artists_on_this_station"}
}

func (f *Radio) ProcessElement(w string, emit func(Radio)) {
	data, err := PrepCsv(w)
	if err != nil {
		return
	}

	s := &Radio{
		Title:                   GetOrDefault(data, 0),
		Artist:                  GetOrDefault(data, 1),
		Description:             GetOrDefault(data, 2),
		Removed:                 ParseRemoved(data, 3),
		Similar_stations:        GetOrDefault(data, 4),
		Artists_on_this_station: GetOrDefault(data, 5),
	}

	fmt.Printf("Radio: %v\n", s)
	emit(*s)
}

func process(s beam.Scope, t interface{}, cols []string, inDir, dsn, table string) {
	// Read
	lines := textio.Read(s, inDir)
	data := beam.ParDo(s, t, lines)
	// Write to DB
	databaseio.Write(s, "mysql", dsn, table, cols, data)
}

func main() {
	// Check flags
	flag.Parse()
	if *trackCsvDir == "" || *playlistsCsvDir == "" || *radioCsvDir == "" || *databasePassword == "" {
		log.Fatalf("Usage: index_match --track-csv_dir=$TRACK_DIR --playlist_csv_dir=$PLAYLIST_DIR --radio_csv_dir=$RADIO_DIR--database_password=$PW [--database_host=$HOST --database_user=$USER]")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", *databaseUser, *databasePassword, *databaseHost, *database)
	fmt.Printf("dsn: %v\n", dsn)
	fmt.Printf("cols: %v\n", Song{}.getCols())
	// Initialize Beam
	beam.Init()
	p := beam.NewPipeline()
	s := p.Root()
	// Songs
	process(s, &Song{}, Song{}.getCols(), *trackCsvDir, dsn, "song")
	// Radio
	process(s, &Radio{}, Radio{}.getCols(), *radioCsvDir, dsn, "radio")
	// Playlists
	process(s, &Playlist{}, Playlist{}.getCols(), *playlistsCsvDir, dsn, "playlists")
	// Run until we find errors
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
