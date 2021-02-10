#!/bin/sh

go-bindata -fs -prefix "ui/dist/" ui/dist/*

go run .