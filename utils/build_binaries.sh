#!/bin/sh

MAIN_FILES=(dao/main.go)
BIN_NAMES=(stpdaov2_test)
ROOT_DIR="$(git rev-parse --show-toplevel)"
BIN_DIR=$ROOT_DIR/bin

set +e

cd $ROOT_DIR

# make package binary
for (( i=0; i<${#MAIN_FILES[@]}; i++ ))
do
	echo "env GOOS=linux GOARCH=amd64 go build -v -o $BIN_DIR/${BIN_NAMES[i]} ${MAIN_FILES[i]}"
	env GOOS=linux GOARCH=amd64 go build -v -o $BIN_DIR/${BIN_NAMES[i]} ${MAIN_FILES[i]}
	if [ $? -ne 0 ]; then
		echo "make binary(${MAIN_FILES[i]}) failed"
		exit 1
	fi
done

set -e


