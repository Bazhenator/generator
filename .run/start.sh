#!/use/bin/bash

set -a
. .env
set +a

go run ../cmd/generator/main.go