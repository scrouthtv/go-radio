#!/bin/sh

ROOT=$(realpath $(dirname $0))
GOPATH="$ROOT" go build -o "$ROOT/bin/go-radio" ...main 
strip "$ROOT/bin/go-radio"
md5sum "$ROOT/bin/go-radio"
