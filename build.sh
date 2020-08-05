#!/bin/bash
cd $GOPATH/src/tmp
go get github.com/beego/bee
bee pack -be GOOS=linux