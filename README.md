# BlockchainTraversal

This project was developed on Ubuntu 16.04.6 LTS

To run first clone then set the GOPATH
export GOPATH={Absoulte Path of clone}/go

Then you will be able to run the project with
go run main.go

Usage of HTTP Web Server

Runs unsecure on port 8080

GET /blockchain
- Returns entire blockchain
POST /blockchain/dump/block
- Post as a form:
-- block "number"
POST /blockchain/search
- Post as a form:
-- key "string"
-- value "string"
