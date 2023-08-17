#!/usr/bin/env bash
for APP in server client
do
  for ARCH in amd64 arm64
  do
    for OS in linux windows darwin/
    do
      echo $APP $OS $ARCH
      if [ $OS == windows ]; then
        GOOS=$OS GOARCH=$ARCH go build -o ./bin/${APP}_${OS}_${ARCH}.exe ./cmd/$APP
      else
        GOOS=$OS GOARCH=$ARCH go build -o ./bin/${APP}_${OS}_${ARCH} ./cmd/$APP
      fi
    done
  done
done