package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Amount: [form.amount != sum(packages[].amount) + sum(packages[].userFee) + shippingFee]
type PaymentsConfirm struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type PaymentsConfirmResponse struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

func (client *Client) PaymentsConfirm(transactionId int64, request *PaymentsConfirm) (response *PaymentsConfirmResponse, err error) {

	body, err := json.Marshal(request)
	res, err := client.post(context.Background(), fmt.Sprintf(endpointV3PaymentsConfirm, transactionId), body)
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
		response = &PaymentsConfirmResponse{}
		if err = json.Unmarshal(bodyBytes, response); err != nil {
			return
		}

	} else {
		err = fmt.Errorf("failed response, StatusCode: %d", res.StatusCode)
		return
	}

	return
}
