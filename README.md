# google-music-takeout-data-pipeline
This repository is based on my [Blog](https://chollinger.com/blog), where I built an over-engineer Data Pipeline using a fairly standard Data Engineering toolset (`bash`, `go`, `Python`, `Apache Beam`/`Dataflow`, `SQL`) to clean up the Google Play Music Takeout data form early 2021.

*This is a tech-demo*

## Requirements
- `go`
- `mysql`/`mariadb`
- `python`
- `bash` (Tested on `Debian 10`)

## Run (Airflow)
Naturally, I suggest to read the blog article linked above.

If you want to run it via `Airflow`, just create a variable called `google-takeout-path` and clone this repo into it. Make sure to `cp conf/config.ini.sample conf/config.ini` and to customize it.

Then move `airflow/dag.py` to your `Airflow` DAG directory.

## Run (Manually)

### CSV-to-SQL
```
PW=PASSWORD
DB_USER=USER
DB_HOST=URI

cd pipelines/

go run csv_to_sql.go utils.go \
	--track_csv_dir="$(pwd)/test_data/clean/csv/*.csv" \
	--playlist_csv_dir="$(pwd)/test_data/clean/playlists/*.csv" \
	--radio_csv_dir="$(pwd)/test_data/clean/radios/*.csv" \
	--database_host=$DB_HOST \
	--database_user=$DB_USER \
	--database_password="${PW}"
```

### Extract Id3v2
```
PW=PASSWORD
DB_USER=USER
DB_HOST=URI

cd pipelines/

find $(pwd)/test_data/clean/mp3 -type f -name "*.mp3" > all_mp3.txt

go run extract_id3v2.go utils.go \
	--mp3_list="$(pwd)/all_mp3.txt" \
	--database_host=$DB_HOST \
	--database_user=$DB_USER \
	--database_password="${PW}"
```