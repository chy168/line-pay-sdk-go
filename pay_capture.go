package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PaymentsCaptureRequest request body of capture api
type PaymentsCaptureRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"` // USD, JPY, TWD, THB
}

// PaymentsCaptureResponse response body of capture api
// `info[].payInfo[].method`: CREDIT_CARD, BALANCE, DISCOUNT
type PaymentsCaptureResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		TransactionID int64  `json:"transactionId"`
		OrderID       string `json:"orderId"`
		PayInfo       []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
	} `json:"info"`
}

// PaymentsCapture Transactions that have set options.payment.capture as false when requesting the Request API payment will be put on hold when the payment is completed with the Confirm API. In order to finalize the payment, an additional purchase with Capture API is required.
func (client *Client) PaymentsCapture(ctx context.Context, transactionId int64, request *PaymentsCaptureRequest) (response *PaymentsCaptureResponse, err error) {

	body, err := json.Marshal(request)
	res, err := client.post(ctx, fmt.Sprintf(endpointV3PaymentsCapture, transactionId), body)
	if err != nil {
		err = fmt.Errorf("PaymentsCapture post error = %v", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, ioerr := ioutil.ReadAll(res.Body)
		if ioerr != nil {
			err = fmt.Errorf("ReadAll read body failed: %s", ioerr.Error())
			return
		}
		response = &PaymentsCaptureResponse{}
		if err = json.Unmarshal(bodyBytes, response); err != nil {
			return
		}

	} else {
		err = fmt.Errorf("failed response, StatusCode: %d", res.StatusCode)
		return
	}

	return
}
