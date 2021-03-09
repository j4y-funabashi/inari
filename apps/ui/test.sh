#!/usr/bin/env bash

# Exit on error. Append "|| true" if you expect an error.
set -o errexit
# Exit on error inside any functions or subshells.
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch the error in case mysqldump fails (but gzip succeeds) in `mysqldump |gzip`
set -o pipefail
# Turn on traces, useful while debugging but commented out by default
# set -o xtrace

function __log_info () {
	local log_msg="${1}"
	echo -e "$(date -u +"%Y-%m-%dT%H:%M:%S") ::: ${log_msg}" 1>&2
}

__log_info "TESTING"

yarn test

__log_info "END"
