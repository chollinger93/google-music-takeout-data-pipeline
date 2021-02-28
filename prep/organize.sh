#!/bin/bash

# Exit on error
set -e

POSITIONAL=()
while [[ $# > 0 ]]; do
    case "$1" in
        -b|--base-dir)
        BASE_DIR="${2}"
        shift 2 
        ;;
        -t|--target-dir)
        TARGET_DIR="${2}"
        shift 2 
        ;;
        *) 
        POSITIONAL+=("$1")
        shift
        ;;
    esac
done

set -- "${POSITIONAL[@]}" 

# Check args
if [[ -z "${BASE_DIR}" || -z "${TARGET_DIR}" ]]; then 
    echo "Usage: organize.sh --base-dir \$PATH_TO_EXPORT --target_dir \$PATH_TO_OUTPUT"
    exit 1
fi 

# Set base folders
TRACK_FOLDER="${BASE_DIR}/Tracks"
PLAYLISTS_FOLDER="${BASE_DIR}/Playlists"
RADIO_FOLDER="${BASE_DIR}/Radio Stations"

# Check folder
if [[ ! -d "${TRACK_FOLDER}" || ! -d "${PLAYLISTS_FOLDER}" || ! -d "${RADIO_FOLDER}" ]]; then 
    echo "--base-dir must point to the folder containing /Tracks, /Playlists, and /Radio\ Stations"
    exit 1
fi

# Define targets
echo "Creating target dirs in $TARGET_DIR..."
mkdir -p "${TARGET_DIR}"
CSV_TARGET="${TARGET_DIR}/csv/"
MP3_TARGET="${TARGET_DIR}/mp3/"
PLAYLIST_TARGET="${TARGET_DIR}/playlists/"
RADIO_TARGET="${TARGET_DIR}/radios/"

# Create structure for CSVs and MP3s
mkdir -p "${CSV_TARGET}"
mkdir -p "${MP3_TARGET}"
mkdir -p "${PLAYLIST_TARGET}"
mkdir -p "${RADIO_TARGET}"

# Copy flat files files
echo "Copying files from $TRACK_FOLDER..."
find "${TRACK_FOLDER}" -type f -iname "*.csv" -exec cp -- "{}" "${CSV_TARGET}" \;
find "${TRACK_FOLDER}" -type f -iname "*.mp3" -exec cp -- "{}" "${MP3_TARGET}" \;

# Playlists need to be concatenated
echo "Copying files from $PLAYLISTS_FOLDER..."
find "${PLAYLISTS_FOLDER}" -type f -mmin +1 -print0 | while IFS="" read -r -d "" f; do
    # Ignore summaries 
    if [[ ! -z $(echo "${f}" | grep "Tracks/" ) ]]; then 
        # Group
        playlistName="$(basename "$(dirname "$(dirname "${f}")")")"
        targetCsv="${PLAYLIST_TARGET}/$playlistName.csv"
        touch "${targetCsv}"
        echo "Working on playlist $playlistName to $targetCsv"
        # Add header
        if [[ -z $(cat "${targetCsv}" | grep "Title,Album,Artist") ]]; then 
            echo "Title,Album,Artist,Duration (ms),Rating,Play Count,Removed,Playlist Index" > "${targetCsv}"
        fi
        echo "$(tail -n +2 "${f}")" >> "${targetCsv}"
    fi
done

# Just flat file copy 
echo "Copying files from $RADIO_TARGET..."
find "${RADIO_FOLDER}" -type f -exec cp -- "{}" "${RADIO_TARGET}" \;