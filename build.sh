#!/bin/bash
cd $GOPATH/src/github.com/beatrice950201/araneid/
ls
rm -rf $GOPATH/src/github.com/qiniu/iconv
go get github.com/beego/bee
go get github.com/qiniu/iconv
bee pack -be GOOS=linux