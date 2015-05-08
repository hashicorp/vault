#!/bin/bash
set -e


MARATHON_URL=http://dev.banno.com:8080 MESOS_URL=http://dev.banno.com:5051,http://dev.banno.com:5050 make testacc TEST=./builtin/credential/marathon
