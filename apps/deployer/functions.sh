#!/usr/bin/env bash

CF_TEMPLATES_DIR="$(cd "${__dir}/../infra" && pwd)"
APPS_DIR="$(cd "${__dir}/.." && pwd)"

function __log_info () {
	local log_msg="${1}"
	echo -e "$(date -u +"%Y-%m-%dT%H:%M:%S") ::: ${log_msg}" 1>&2
}

deploy_infra() {
    __log_info "${PROJECT_NAME} DEPLOYING INFRA from ${CF_TEMPLATES_DIR} to ${CURRENT_ENV}"
    deploy_ui_s3
}

deploy_apps() {
    __log_info "${PROJECT_NAME} DEPLOYING APPS from ${APPS_DIR} to ${CURRENT_ENV}"
    deploy_ui
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

deploy_ui() {
    __log_info "Deploying UI"
    S3_UI_BUCKET_NAME="inari-ui-dev"
    yarn --cwd "${APPS_DIR}/ui" build
    aws s3 cp "${APPS_DIR}/ui/build" s3://$S3_UI_BUCKET_NAME --recursive

}
