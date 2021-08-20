#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[${1:-0}]}")" && pwd)"


S3_STACK_NAME="${PROJECT_NAME}-s3-ui-${CURRENT_ENV}"
DYNAMODB_STACK_NAME="${PROJECT_NAME}-dynamodb-${CURRENT_ENV}"

. "$__dir/functions.sh"

deploy_infra
deploy_apps
