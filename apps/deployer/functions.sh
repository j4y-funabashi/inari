#!/usr/bin/env bash

CF_TEMPLATES_DIR="$(cd "${__dir}/../infra" && pwd)"
APPS_DIR="$(cd "${__dir}/.." && pwd)"
BASE_DOMAIN="funabashi.co.uk"
if [[ "${CURRENT_ENV}" == "prod" ]]; then
    CF_BASE_DOMAIN="photos.${BASE_DOMAIN}"
else
    CF_BASE_DOMAIN="photos-dev.${BASE_DOMAIN}"
fi


function __log_info () {
	local log_msg="${1}"
	echo -e "$(date -u +"%Y-%m-%dT%H:%M:%S") ::: ${log_msg}" 1>&2
}

deploy_infra() {
    __log_info "${PROJECT_NAME} DEPLOYING INFRA from ${CF_TEMPLATES_DIR} to ${CURRENT_ENV}"
    deploy_s3_mediastore
    deploy_dynamodb
    deploy_cloudfront
    deploy_cognito
}

deploy_apps() {
    __log_info "${PROJECT_NAME} DEPLOYING APPS from ${APPS_DIR} to ${CURRENT_ENV}"
    deploy_ui
    deploy_api
}

deploy_s3_mediastore() {
    __log_info "Deploying ${S3_MEDIASTORE_STACK_NAME}"
    aws cloudformation deploy \
        --template-file "${CF_TEMPLATES_DIR}/s3-mediastore.yml" \
        --stack-name "${S3_MEDIASTORE_STACK_NAME}" \
        --no-fail-on-empty-changeset \
        --parameter-overrides \
        "EnvironmentName=${CURRENT_ENV}"
}

deploy_dynamodb() {
    __log_info "Deploying ${DYNAMODB_STACK_NAME}"
    aws cloudformation deploy \
        --template-file "${CF_TEMPLATES_DIR}/dynamodb.yml" \
        --stack-name "${DYNAMODB_STACK_NAME}" \
        --no-fail-on-empty-changeset \
        --parameter-overrides \
        "EnvironmentName=${CURRENT_ENV}"
}

deploy_cloudfront() {
    CERTIFICATE_ARN=`aws --region us-east-1 acm list-certificates \
        --query "CertificateSummaryList[?DomainName=='*.${BASE_DOMAIN}'].CertificateArn" \
        --output text`

    __log_info "Deploying ${CLOUDFRONT_STACK_NAME}"
    __log_info "BASE_DOMAIN: ${BASE_DOMAIN}"
    __log_info "CF_BASE_DOMAIN: ${CF_BASE_DOMAIN}"
    __log_info "CERTIFICATE_ARN ${CERTIFICATE_ARN}"

    aws cloudformation deploy \
        --template-file "${CF_TEMPLATES_DIR}/cloudfront.yml" \
        --stack-name "${CLOUDFRONT_STACK_NAME}" \
        --no-fail-on-empty-changeset \
        --parameter-overrides \
        "EnvironmentName=${CURRENT_ENV}" \
        "CloudFrontCertificate=${CERTIFICATE_ARN}" \
        "CloudFrontBaseDomain=${CF_BASE_DOMAIN}"
}

deploy_cognito() {
    __log_info "Deploying ${COGNITO_STACK_NAME}"
    aws cloudformation deploy \
        --template-file "${CF_TEMPLATES_DIR}/cognito.yml" \
        --stack-name "${COGNITO_STACK_NAME}" \
        --no-fail-on-empty-changeset \
        --parameter-overrides \
        "EnvironmentName=${CURRENT_ENV}" \
        "CallbackURL=https://${CF_BASE_DOMAIN}"
}

deploy_ui() {
    __log_info "Deploying UI"
    S3_UI_BUCKET_NAME=`aws cloudformation describe-stacks \
        --query "Stacks[0].Outputs[?OutputKey=='S3InariUIBucketName'].OutputValue" --output text \
        --stack-name ${CLOUDFRONT_STACK_NAME}`
    CLOUDFRONT_DISTRIBUTION_ID=`aws cloudformation describe-stacks \
        --query "Stacks[0].Outputs[?OutputKey=='CloudFrontDistributionID'].OutputValue" --output text \
        --stack-name ${CLOUDFRONT_STACK_NAME}`

    __log_info "Bucket name: ${S3_UI_BUCKET_NAME}"
    __log_info "Distribution ID: ${CLOUDFRONT_DISTRIBUTION_ID}"

    yarn --cwd "${APPS_DIR}/ui" build
    aws s3 cp "${APPS_DIR}/ui/build" s3://$S3_UI_BUCKET_NAME --recursive

    aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID \
        --paths /index.html /error.html
}

deploy_api() {
    __log_info "Deploying API"
    APIGATEWAYID=`aws apigatewayv2 get-apis \
        --query "Items[?Name=='inari-photos-api-${CURRENT_ENV}'].ApiId" --output text`

    __log_info "GATEWAY_ID: ${APIGATEWAYID}"
    cd "${APPS_DIR}/api/lambda_functions"
    APIGATEWAYID=${APIGATEWAYID} make deploy
}
