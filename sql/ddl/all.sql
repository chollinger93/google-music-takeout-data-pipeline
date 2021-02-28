CREATE DATABASE IF NOT EXISTS `googleplay_raw`;

CREATE OR REPLACE TABLE `googleplay_raw`.`id3v2` (
    `artist`	    TEXT,
    `title`	        TEXT,
    `album`	        TEXT,
    `year`          INTEGER,
    `genre`         TEXT,
    `trackPos`      INTEGER,
    `maxTrackPos`   INTEGER,
    `partOfSet`     TEXT,
    `file_name`     TEXT,
    `full_path`     TEXT
);

CREATE OR REPLACE TABLE `googleplay_raw`.`playlists` (
    `title`			TEXT,
    `album`			TEXT,
    `artist`		TEXT,
    `duration_ms`	INTEGER,
    `rating`		INTEGER,
    `play_count`	INTEGER,
    `removed`		BOOLEAN,
    `playlist_ix`	INTEGER
);

CREATE OR REPLACE TABLE `googleplay_raw`.`radio` (
    `title`	                        TEXT,
    `artist`	                    TEXT,
    `description`	                TEXT,
    `removed`	                    BOOLEAN,
    `similar_stations`	            TEXT,
    `artists_on_this_station`	    TEXT
);

CREATE OR REPLACE TABLE `googleplay_raw`.`song` (
    `title`	        TEXT,
    `album`	        TEXT,
    `artist`	    TEXT,
    `duration_ms`	INTEGER,
    `rating`	    INTEGER,
    `play_count`	INTEGER,
    `removed`	    BOOLEAN,
    `file_name`     TEXT
);

CREATE DATABASE IF NOT EXISTS `googleplay_master`;
CREATE OR REPLACE TABLE `googleplay_master`.`id3v2` (
    `artist`	    TEXT,
    `title`	        TEXT,
    `album`	        TEXT,
    `year`          INTEGER,
    `genre`         TEXT,
    `trackPos`      INTEGER,
    `maxTrackPos`   INTEGER,
    `partOfSet`     TEXT,
    `file_name`     TEXT,
    `full_path`     TEXT
);