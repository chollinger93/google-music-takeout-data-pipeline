CREATE DATABASE IF NOT EXISTS `googleplay_raw`;

CREATE OR REPLACE TABLE `googleplay_raw`.`radio` (
    `title`	                        TEXT,
    `artist`	                    TEXT,
    `description`	                TEXT,
    `removed`	                    BOOLEAN,
    `similar_stations`	            TEXT,
    `artists_on_this_station`	    TEXT
);