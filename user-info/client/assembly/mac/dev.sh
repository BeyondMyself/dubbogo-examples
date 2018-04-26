#!/usr/bin/env bash
# ******************************************************
# DESC    : build script for dev env
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : LGPL V3
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2017-10-18 13:24
# FILE    : dev.sh
# ******************************************************


set -e

export GOOS=darwin
export GOARCH=amd64

PROFILE=dev

PROJECT_HOME=`pwd`

if [ -f "${PROJECT_HOME}/assembly/common/app.properties" ]; then
. ${PROJECT_HOME}/assembly/common/app.properties
fi


if [ -f "${PROJECT_HOME}/assembly/common/build.sh" ]; then
. ${PROJECT_HOME}/assembly/common/build.sh
fi
