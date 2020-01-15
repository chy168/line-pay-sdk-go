[![GitHub license](https://img.shields.io/badge/license-Apache--2.0-blue)](https://raw.githubusercontent.com/chy168/line-pay-sdk-go/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/chy168/line-pay-sdk-go?status.svg)](https://godoc.org/github.com/chy168/line-pay-sdk-go)

# line-pay-sdk-go
LINE Pay API SDK for Go
[LINE Pay Developers Website](https://pay.line.me/developers/apis/onlineApis)

# Supported LINE Pay v3 APIs
---------------
- [x] Request API 
- [x] Confirm API 
- [x] Capture API (not test yet)
- [ ] Void API
- [ ] Refund API
- [x] Payment Details API 
- [ ] Check Payment Status API
- [ ] Check RegKey API
- [ ] Pay Preapproved API
- [ ] Expire RegKey API

# Usage
```
go get -v github.com/chy168/line-pay-sdk-go
```

# How to test
## develop
replace necessary information in `data_test.go`, then you can `go test` what you want to try.

## test
there is a built in web server to perform confirm api by transaction (can be used as confirmURL)
```
go run examples/cmd/callback_server.go --channel-id=<YOUR_CHANNEL_ID> --channel-secret=<YOUR_CHANNEL_SECRET>
```

# LICENSE
Apache 2.0


