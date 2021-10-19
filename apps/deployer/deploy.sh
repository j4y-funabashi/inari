#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[${1:-0}]}")" && pwd)"


S3_MEDIASTORE_STACK_NAME="${PROJECT_NAME}-s3-mediastore-${CURRENT_ENV}"
DYNAMODB_STACK_NAME="${PROJECT_NAME}-dynamodb-${CURRENT_ENV}"
CLOUDFRONT_STACK_NAME="${PROJECT_NAME}-cloudfront-${CURRENT_ENV}"
COGNITO_STACK_NAME="${PROJECT_NAME}-cognito-${CURRENT_ENV}"

. "$__dir/functions.sh"

deploy_infra
deploy_apps
