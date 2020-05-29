#!/usr/bin/env bash
set -e

VERSION=$1
# Get the version from the environment, or try to figure it out.
if [ -z "$VERSION" ]; then
	VERSION=$(awk -F\" '/version =/ { print $2; exit }' < version.go)
fi
if [ -z "$VERSION" ];then
  echo "Please specify a version."
  exit 1
fi

VERSION="v$VERSION"
SRC_DIRS="common sdk sdkimpl utest"
DIST_DIR="build/dist"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

echo "===> disting sdk..."
tar -zcf "$DIST_DIR/sdk_$VERSION.tar.gz" $SRC_DIRS

exit 0