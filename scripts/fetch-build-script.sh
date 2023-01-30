#!/bin/bash
# set -exo pipefail

cd $(cd `dirname "$0"`; cd ..; pwd)

# Defaults
VERSION=2.3.0
INSTALL_LOCATION="bin"
BUILD_SCRIPT="${INSTALL_LOCATION}/build"

if [[ -f "${BUILD_SCRIPT}" ]]; then
  exit 0
fi

mkdir ${INSTALL_LOCATION}

curl -LO https://raw.githubusercontent.com/Jordan-Cartwright/go-build-script/v${VERSION}/bin/build

chmod +x build
mv build ${BUILD_SCRIPT}
