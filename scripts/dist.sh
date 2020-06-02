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
project_path=$(pwd)
project_name="${project_path##*/}"
echo "==> Building $project_name $VERSION..."

cd scripts

if [[ -f "download.sh" ]];then
  source download.sh
fi

if [[ -f "build.sh" ]];then
  source build.sh
fi

if [[ -f "pkg.sh" ]];then
  source pkg.sh
fi

echo
echo "======> BUILD SUCCESS!"

exit 0