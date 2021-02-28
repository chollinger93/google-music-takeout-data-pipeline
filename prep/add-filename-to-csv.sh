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
        *) 
        POSITIONAL+=("$1")
        shift
        ;;
    esac
done


set -- "${POSITIONAL[@]}" 

# Check args
if [[ -z "${BASE_DIR}" ]]; then 
    echo "Usage: add_filename_to_csv.sh --base-dir \$PATH_TO_EXPORT"
    exit 1
fi 

find "${BASE_DIR}" -type f -name "*.csv" -print0 | while IFS= read -r -d '' f; do
    rm -f /tmp/out.csv
    bname=$(basename "${f}")

    ix=0
    while read l; do
        if [[ $ix -eq 0 ]]; then 
            echo "${l},file_name" >> /tmp/out.csv
        else 
            echo "${l},\"${bname}\"" >> /tmp/out.csv
        fi 
        ix=$(($ix + 1))
    done < "$f"
    IFS=$OLDIFS
    if [[ $ix -ne 0 ]]; then 
        echo "Adding $bname to $f"
        cat "/tmp/out.csv" > "${f}"
    else 
        echo "File $f is empty"
    fi
done

