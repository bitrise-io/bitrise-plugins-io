#!/usr/bin/env bash
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
go build -C "$THIS_SCRIPT_DIR" -o "./_tmp/plugin"
"$THIS_SCRIPT_DIR/_tmp/plugin" $*
