#!/bin/bash

function GetDirs() {
  IDIRS=()
  i=1
  for _ in $(cat ./scripts/.distignore)
  do
    NUM=$i
    IDIR=$(awk 'NR=='$NUM' {print $1}' ./scripts/.distignore)
    if [[ -n "$IDIR" ]]; then
      IDIRS[$i]=$IDIR
    fi

    : $(( i++ ))
  done

  for f in `ls -l $PWD`
  do
    if [[ -d "$f" ]];then
      b=0
      for id in "${IDIRS[@]}"
      do
        if [[ "$id" == "$f" ]];then
          b=1
        fi
      done

      if [[ $b == 0 ]];then
        TDIRS[${#TDIRS[*]}]=$f
      fi
    fi
  done
  return 0
}

DIST_DIR=./build/dist/
TDIRS=()
cd ..

echo "==> Removing old directory..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

echo "==> Tar files..."
GetDirs

SRC_DIRS=""
for d in "${TDIRS[@]}"
do
  if [[ -z "$SRC_DIRS" ]];then
    SRC_DIRS="$d"
  else
    SRC_DIRS="$SRC_DIRS"" $d"
  fi

done
tar -zcf "$DIST_DIR$project_name""_$VERSION.tar.gz" $SRC_DIRS

# Make the checksums.
pushd "$DIST_DIR" > /dev/null
shasum -a256 ./* > "$project_name"_SHA256SUMS
popd >/dev/null

echo ""
echo "==> Build results:"
echo "==> Path: $DIST_DIR"
echo "==> Files: "
ls -hl "$DIST_DIR"