#!/bin/sh

rm -rf owndns
rm -rf owndns-linux
rm -rf owndns-arm

BT=`date +"%Y-%m-%dT%H:%M:%S"`

go-bindata -fs -prefix "ui/dist/" ui/dist/*

# Mikrotik hAP ac2
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o owndns-arm -ldflags "-X main.BuildTime=$BT"

GOOS=linux GOARCH=amd64 go build -o owndns-linux -ldflags "-X main.BuildTime=$BT"

go build -o owndns -ldflags "-X main.BuildTime=$BT"