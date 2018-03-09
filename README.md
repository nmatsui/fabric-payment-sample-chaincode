# fabric-payment-sample-chaincode
## Description
A sample chaincode for [Hyperledger/fabric](https://github.com/hyperledger/fabric) version 1.1.
This chaincode implements some features like below:
- list accounts.
- retrieve, create, update, delete an account.
- deposit to an account.
- remit from an account to another account.
- withdraw from an account.
- show the histories of an account.

## See also
[fabric-payment-sample-api](https://github.com/nmatsui/fabric-payment-sample-api)  
[fabric-payment-sample-docker](https://github.com/nmatsui/fabric-payment-sample-docker)


## Requirement
||version|
|:--|:--|
|go|1.10|
|Hyperledger/fabric|1.1.0-rc1|

## How to build
### get the development libraries of Hyperledger/fabric 1.1
```bash
$ go get -d github.com/hyperledger/fabric/protos/peer
$ go get -u --tags nopkcs11 github.com/hyperledger/fabric/core/chaincode/shim
```

To avoid build failure, you have to get `fabric/core/chaincode/shim` from nopkcs11 tag.

### get source code to your $GOPATH
```bash
$ go get -u github.com/nmatsui/fabric-payment-sample-chaincode
```

### build this chaincode
```bash
$ go build --tags nopkcs11 fabric-payment.go
```

## Contribution
1. Fork this project ( https://github.com/nmatsui/fabric-payment-sample-chaincode )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## License
[Apache License, Version 2.0](/LICENSE)

## Copyright
Copyright (c) 2018 Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
