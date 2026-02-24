set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]

default:
    @just --list

build:
    go build ./cmd/importer
    go build ./cmd/moodtea

import config="config.json" out="data":
    go run ./cmd/importer -config {{config}} -out {{out}}

cli config="config.json":
    go run ./cmd/moodtea -config {{config}}

run config="config.json" out="data":
    go run ./cmd/importer -config {{config}} -out {{out}}
    go run ./cmd/moodtea -config {{config}}
