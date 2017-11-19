#!/bin/sh

cd tests
go test -tags=integration -profile=local -path=../config
cd ../
