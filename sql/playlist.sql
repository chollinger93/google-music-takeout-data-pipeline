CREATE DATABASE IF NOT EXISTS `googleplay_raw`;

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