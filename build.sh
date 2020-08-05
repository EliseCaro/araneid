#!/bin/bash
cd $GOPATH/src/github.com/beatrice950201/araneid/
go get github.com/beego/bee
bee pack -be GOOS=linux