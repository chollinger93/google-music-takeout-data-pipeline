package main

// go get -u github.com/bogem/id3v2

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/databaseio"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"github.com/bogem/id3v2"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mp3File          = flag.String("mp3_list", "", "File containing all Mp3 files")
	databaseHost     = flag.String("database_host", "localhost", "Database host")
	databaseUser     = flag.String("database_user", "", "Database user")
	databasePassword = flag.String("database_password", "", "Database password")
	database         = flag.String("database", "googleplay_raw", "Database name")
)

type Id3Data struct {
	Artist    string `column:"artist"`   // TPE1
	Title     string `column:"title"`    // TIT2
	Album     string `column:"album"`    // TALB
	Year      int    `column:"year"`     // TYER
	Genre     string `column:"genre"`    // TCON
	TrackPos  int    `column:"trackPos"` // TRCK
	TrackMax  int    `column:"maxTrackPos"`
	PartOfSet string `column:"partOfSet"` // TPOS
	FileName  string `column:"file_name"`
	FullPath  string `column:"full_path"`
	//Other     string
}

func (s Id3Data) getCols() []string {
	return []string{"artist", "title", "album", "year", "genre", "trackPos", "maxTrackPos", "partOfSet", "file_name", "full_path"}
}

func convTrackpos(pos string) (int, int) {
	if pos != "" {
		data := strings.Split(pos, "/")
		fmt.Println(data)
		if len(data) > 1 {
			trck, err1 := strconv.Atoi(data[0])
			mtrck, err2 := strconv.Atoi(data[1])
			if err1 == nil && err2 == nil {
				return trck, mtrck
			}
		} else if len(data) == 1 {
			trck, err := strconv.Atoi(data[0])
			if err == nil {
				return trck, -1
			}
		}
	}
	return -1, -1
}

func getAllTags(tag id3v2.Tag) []string {
	all := tag.AllFrames()
	keys := make([]string, len(all))

	i := 0
	for k := range all {
		// Skip APIC, PRIV (album image binary)
		if k != "APIC" {
			keys[i] = k
		}

		i++
	}

	for i := range keys {
		fmt.Printf("%v: %v\n", keys[i], tag.GetTextFrame(keys[i]))
	}
	return keys
}

func parseMp3ToTags(path string) (Id3Data, error) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return Id3Data{}, err
	}
	defer tag.Close()
	//getAllTags(*tag)
	year, _ := strconv.Atoi(tag.GetTextFrame("TYER").Text)
	tpos, mpos := convTrackpos(tag.GetTextFrame("TRCK").Text)
	return Id3Data{
		Artist:    tag.GetTextFrame("TPE1").Text,
		Title:     tag.GetTextFrame("TIT2").Text,
		Album:     tag.GetTextFrame("TALB").Text,
		Genre:     tag.GetTextFrame("TCON").Text,
		Year:      year,
		TrackPos:  tpos,
		TrackMax:  mpos,
		PartOfSet: tag.GetTextFrame("TPOS").Text,
		FileName:  filepath.Base(path),
		FullPath:  path,
		//Other:     getAllTags(tag),
	}, nil
}

func (f *Id3Data) ProcessElement(path string, emit func(Id3Data)) {
	data, err := parseMp3ToTags(path)
	if err != nil {
		fmt.Printf("Error: %\n", err)
		return
	}
	emit(data)
}

func process(s beam.Scope, t interface{}, cols []string, inDir, dsn, table string) {
	// Read
	lines := textio.Read(s, inDir)
	data := beam.ParDo(s, t, lines)
	// Write to DB
	databaseio.Write(s, "mysql", dsn, table, cols, data)
}

func main() {
	flag.Parse()
	if *mp3File == "" || *databasePassword == "" {
		log.Fatalf("Usage: extract_id3v2 --mp3_list=$TRACK_DIR --database_password=$PW [--database_host=$HOST --database_user=$USER]")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", *databaseUser, *databasePassword, *databaseHost, *database)

	// Initialize Beam
	beam.Init()
	p := beam.NewPipeline()
	s := p.Root()
	process(s, &Id3Data{}, Id3Data{}.getCols(), *mp3File, dsn, "id3v2")
	// Run until we find errors
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
