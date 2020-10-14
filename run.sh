#!/bin/sh

echo "generating jx3 pipeline catalog"

mkdir -p ../jx3-pipeline-catalog/packs
./bin/jx-v2-tekton-converter --out-dir ../jx3-pipeline-catalog/packs
