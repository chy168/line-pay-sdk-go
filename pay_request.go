package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Amount: [form.amount != sum(packages[].amount) + sum(packages[].userFee) + shippingFee]
type PaymentsRequest struct {
	Amount       int                         `json:"amount"`
	Currency     string                      `json:"currency"`
	OrderId      string                      `json:"orderId"`
	Packages     []PaymentsPackageRequest    `json:"packages"`
	RedirectUrls PaymentsRedirectUrlsRequest `json:"redirectUrls"`
}

// Amount: packages[].amount != sum(packages[].products[].quantity * packages[].products[].price)]
type PaymentsPackageRequest struct {
	Id       string                          `json:"id"`
	Amount   int                             `json:"amount"`
	Name     string                          `json:"name"`
	Products []PaymentsPackageProductRequest `json:"products"`
}

type PaymentsPackageProductRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

const (
	PaymentsConfirmUrlTypeClient string = "CLIENT"
	PaymentsConfirmUrlTypeServer string = "SERVER"
	PaymentsConfirmUrlTypeNone   string = "NONE"
)

// NOTE: for the behavior of `ConfirmUrl`, when `ConfirmUrlType` set to `PaymentsConfirmUrlTypeClient`,
// LINE server will send user to the `ConfirmUrl` with `transactionId`. `orderID` won't send for this case.
// the exception is, when user use *QR scanner* at the `waitPreLogin` page (login by LINE account or QR Code scan page)
// server will send user to `ConfirmUrl` with `orderID`
type PaymentsRedirectUrlsRequest struct {
	ConfirmUrlType string `json:"confirmUrlType"`
	ConfirmUrl     string `json:"confirmUrl"`
	CancelUrl      string `json:"cancelUrl"`
}

// response
type PaymentsResponse struct {
	ReturnCode    string               `json:"returnCode"`
	ReturnMessage string               `json:"returnMessage"`
	Info          PaymentsInfoResponse `json:"info"`
}

type PaymentsInfoResponse struct {
	PaymentURL         PaymentsInfoPaymentURLResponse `json:"paymentUrl"`
	TransactionID      int64                          `json:"transactionId"`
	PaymentAccessToken string                         `json:"paymentAccessToken"`
}

type PaymentsInfoPaymentURLResponse struct {
	Web string `json:"web"`
	App string `json:"app"`
}

func (client *Client) PaymentsRequest(ctx context.Context, request *PaymentsRequest) (response *PaymentsResponse, err error) {

	body, err := json.Marshal(request)
	res, err := client.post(ctx, endpointV3PaymentsRequest, body)
	if err != nil {
		err = fmt.Errorf("PaymentsRequest post error = %v", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, ioerr := ioutil.ReadAll(res.Body)
		if ioerr != nil {
			err = fmt.Errorf("ReadAll read body failed: %s", ioerr.Error())
			return
		}
		response = &PaymentsResponse{}
		if err = json.Unmarshal(bodyBytes, response); err != nil {
			return
		}

	} else {
		err = fmt.Errorf("failed response, StatusCode: %d", res.StatusCode)
		return
	}

	return
}
