SELECT 
	s.artist AS csv_arist
	,i.artist AS mp3_artist 
	,s.album  AS csv_album
	,i.album AS mp3_album
	,s.title AS csv_title
	,i.title  AS mp3_title
	,s.file_name 
FROM googleplay_raw.song AS s
LEFT OUTER JOIN googleplay_raw.id3v2 AS i
ON (
	TRIM(LOWER(i.artist)) = TRIM(LOWER(s.artist))
	AND TRIM(LOWER(i.album)) = TRIM(LOWER(s.album))
	AND TRIM(LOWER(i.title)) = TRIM(LOWER(s.title))
	)
WHERE i.artist IS NOT NULL
ORDER BY i.artist ASC
