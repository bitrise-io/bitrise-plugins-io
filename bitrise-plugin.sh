#!/usr/bin/env bash
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
go build -o "$THIS_SCRIPT_DIR/_tmp/plugin"
"$THIS_SCRIPT_DIR/_tmp/plugin" $*
