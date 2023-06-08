#!/bin/bash

# Set the directory to check
dir="/opt/vigo360/assets/images"

find $dir -type f -name "*.webp" -mmin -1 -exec bash -c '/usr/local/bin/magick convert "$0" "${0%.webp}.avif"' {} \;


# Loop through all webp files in the directory
#for file in "$dir"/*.webp; do
  # Check if the corresponding avif file exists
#  if [ ! -f "${file%.webp}.avif" ]; then
    # If it doesn't exist, convert the webp file to avif using magick
#    /usr/local/bin/magick convert "$file" "${file%.webp}.avif"
#  fi
#done

