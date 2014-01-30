#!/usr/bin/env bash

function log {
	[ "$DEBUG" == "1" ] && echo "$@"
}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cmd="sqlite3 $DIR/app.db < $DIR/schema.sql"
log "cmd=$cmd"
eval $cmd
