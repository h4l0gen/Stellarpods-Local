#!/bin/bash

# Build dbscanner

# Exit if any steps fail
set -e

CWD=${PWD}

GO_CMD=${GO_CMD:-"build"}
BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
VERBOSE=${VERBOSE:-}

repo_path="github.com/j4ck4l-24/StellarPods"
main_package="github.com/j4ck4l-24/StellarPods/cmd"

# Extract the go version
go_version=$(go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')

# go 1.4 requires ldflags format to be "-X key value", not "-X key=value"
# ldseparator here is for cross compatibility between go versions

ldseparator="="
if [ "${go_version:0:3}" = "1.4" ]; then
	ldseparator=" "
fi

ldflags="
  -X ${repo_path}/version.Version${ldseparator}${version}
  -X ${repo_path}/version.Revision${ldseparator}${revision}
  -X ${repo_path}/version.Branch${ldseparator}${branch}
  -X ${repo_path}/version.BuildUser${ldseparator}${BUILD_USER}
  -X ${repo_path}/version.BuildDate${ldseparator}${BUILD_DATE}
  -X ${repo_path}/version.GoVersion${ldseparator}${go_version}"

echo ">>> Building Stellar Pods binary..."

if [ -n "$VERBOSE" ]; then
  echo "Building with -ldflags $ldflags"
fi

GOBIN=$PWD CGO_ENABLED=0 go "${GO_CMD}" -o "${CWD}/bin/stellarpods" ${GO_FLAGS} -ldflags "${ldflags}" "${main_package}"

echo "[*] Build Complete."
exit 0