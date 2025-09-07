#!/bin/bash
set -o pipefail

# gocolourize - add color to go test output
# usage:
#   ./scripts/gocolourize.sh go test ./...
#

$@ | awk 'BEGIN {
    RED="\033[31m"
    GREEN="\033[32m"
    CYAN="\033[36m"
    BRRED="\033[91m"
    BRGREEN="\033[92m"
    BRCYAN="\033[96m"
    NORMAL="\033[0m"
}
         { color=NORMAL }
/^ok /   { color=BRGREEN }
/^FAIL/  { color=BRRED }
/^SKIP/  { color=BRCYAN }
/PASS:/  { color=GREEN }
/FAIL:/  { color=RED }
/SKIP:/  { color=CYAN }
         { print color $0 NORMAL }
'
