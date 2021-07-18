#!/bin/zsh

# go install すれば $GOPATH/bin に置かれるのでこのスクリプトは不要
echo dmainfinder build
go build -o domainfinder

echo synonyms build
cd ../synonyms
go build -o ../domainfinder/lib/synonyms

echo available build
cd ../available
go build -o ../domainfinder/lib/available

echo sprinkle build
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle

echo coolify build
cd ../coolify
go build -o ../domainfinder/lib/coolify

echo domainify
cd ../domainify
go build -o ../domainfinder/lib/domainify

echo Done
