#!/bin/bash
set -e

FILES="$(ls *[!_test].go | tr '\n' ' ')"
codecgen -d 100 -o structs.generated.go ${FILES}
