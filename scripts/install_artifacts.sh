#!/usr/bin/env bash
set -euo pipefail

# this script will install 2rm assets built on the local machine
# you will have to build 2rm before running this script
#
# you should be running this script through
# $ task build && sudo task install

if [ "$EUID" -ne 0 ]; then
    echo "Installer requires root privileges."
    exit 1
fi

if [ $# -lt 1 ]; then
    echo "Not enough arguments provided"
    echo "Usage: scripts/install.sh <build_output>"
    exit 2
fi

build_output=$1

if [ ! -d $build_output ]; then
    echo "could not find build output in $build_output"
    echo "try building with the 'task build' command"
    exit 3
fi

cp "./$build_output/2rm" "/usr/local/bin/2rm"
cp "./$build_output/2rm.1" "/usr/share/man/man1/2rm.1"

