#!/usr/bin/env bash

function __log_info () {
	local log_msg="${1}"
	echo -e "$(date -u +"%Y-%m-%dT%H:%M:%S") ::: ${log_msg}" 1>&2
}


deploy_infra() {
    CF_TEMPLATES_DIR="apps/infra"

    __log_info "${PROJECT_NAME} DEPLOYING INFRA to ${CURRENT_ENV}"
    deploy_ui_s3
}

deploy_ui_s3() {
    __log_info "Deploying ${S3_STACK_NAME}"
    aws cloudformation deploy \
        --template-file "${CF_TEMPLATES_DIR}/s3.yml" \
        --stack-name "${S3_STACK_NAME}" \
        --no-fail-on-empty-changeset \
        --parameter-overrides \
        "EnvironmentName=${CURRENT_ENV}"
}
