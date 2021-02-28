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