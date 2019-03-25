#!/bin/sh

TMP_PATH="$(mktemp -d)"
mkdir "${TMP_PATH}/webasset_package_name"
TMP_FILE="${TMP_PATH}/webasset_package_name/compression.go"

${CMD} -f "${TMP_FILE}" "${TEST_PATH}"

if [ -f "${TMP_FILE}" ]; then
    grep "package webasset_package_name" "${TMP_FILE}"
    if [ $? == 0 ]; then
        echo "package name auto-generated correctly"
        rm -rf "${TMP_PATH}"
        exit 0
    fi
    rm -rf "${TMP_PATH}"
    exit 1
fi

rm -rf "${TMP_PATH}"
exit 1