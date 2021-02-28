SELECT iv2.artist, iv2.album, iv2.title, iv2.full_path 
FROM  googleplay_master.id3v2 AS iv2 
WHERE NOT EXISTS (
	SELECT iv.artist, iv.album, iv.title FROM googleplay_raw.id3v2 iv 
) 