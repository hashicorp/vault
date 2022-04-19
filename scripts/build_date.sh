#!/bin/bash
# it's tricky to do an RFC3339 format in a cross platform way, so we hardcode UTC
DATE_FORMAT="%Y-%m-%dT%H:%M:%SZ"

# we're using this for build date because it's stable across platform builds
git show -s --format=%cd --date=format:"$DATE_FORMAT" HEAD
