# go-iap-onestore

![](https://img.shields.io/badge/golang-1.19-blue.svg?style=flat)

go-iap-onestore verifies the purchase receipt via OneStore.
This repository is inspired by [go-iap](https://github.com/awa/go-iap)


# Installation
```
go get github.com/coolishbee/go-iap-onestore
```


# Quick Start

### In App Purchase via One Store

```go
import(
    "github.com/coolishbee/go-iap-onestore"
)

func main() {
	client := onestore.New("client_id", "client_secret", "purchaseToken")

	ctx := context.Background()
	resp, err := client.Verify(ctx, "package", "productID", "purchaseToken")
}
```

# ToDo
- [ ] Monthly subscription product
- [ ] subscription product


# Support

This validator uses [One store in-app payment server API (API V7)](https://dev.onestore.co.kr/devpoc/reference/view/Tools).
