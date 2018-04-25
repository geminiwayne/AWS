#!/usr/bin/env bash

dep ensure

echo "Compiling functions to bin/ ..."

rm -rf bin/

cd src/main/
for f in *.go; do
  filename="${f%.go}"
  if GOOS=linux go build -o "../../bin/$filename" ${f}; then
    echo "✓ Compiled $filename"
  else
    echo "✕ Failed to compile $filename!"
    exit 1
  fi
done

echo "Done."
