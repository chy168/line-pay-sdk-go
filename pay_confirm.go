package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Amount: [form.amount != sum(packages[].amount) + sum(packages[].userFee) + shippingFee]
type PaymentsConfirmRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type PaymentsConfirmResponse struct {
	ReturnCode    string                      `json:"returnCode"`
	ReturnMessage string                      `json:"returnMessage"`
	Info          PaymentsConfirmInfoResponse `json:"info"`
	Amount        int                         `json:"amount"`
	Currency      string                      `json:"currency"`
}

type PaymentsConfirmInfoResponse struct {
	OrderID                 string                                `json:"orderId"`
	TransactionID           int64                                 `json:"transactionId"`
	AuthorizationExpireDate time.Time                             `json:"authorizationExpireDate"`
	RegKey                  string                                `json:"regKey"`
	PayInfo                 []PaymentsConfirmInfoPayInfoResponse  `json:"payInfo"`
	Packages                []PaymentsConfirmInfoPackagesResponse `json:"packages"`
	Shipping                PaymentsConfirmInfoShippingResponse   `json:"shipping"`
}

type PaymentsConfirmInfoPayInfoResponse struct {
	Method                 string `json:"method"`
	Amount                 int    `json:"amount"`
	CreditCardNickname     string `json:"creditCardNickname"`
	CreditCardBrand        string `json:"creditCardBrand"`        // VISA, MASTER, AMEX, DINERS, JCB
	MaskedCreditCardNumber string `json:"maskedCreditCardNumber"` // Format: **** **** **** 1234
}

type PaymentsConfirmInfoPackagesResponse struct {
	ID            string `json:"id"`
	Amount        int    `json:"amount"`
	UserFeeAmount int    `json:"userFeeAmount"`
}

type PaymentsConfirmInfoShippingResponse struct {
	MethodID  string                                     `json:"methodId"`
	FeeAmount int                                        `json:"feeAmount"`
	Address   PaymentsConfirmInfoShippingAddressResponse `json:"address"`
}

type PaymentsConfirmInfoShippingAddressResponse struct {
	Country    string                                              `json:"country,omitempty"`
	PostalCode string                                              `json:"postalCode,omitempty"`
	State      string                                              `json:"state,omitempty"`
	City       string                                              `json:"city,omitempty"`
	Detail     string                                              `json:"detail,omitempty"`
	Optional   string                                              `json:"optional,omitempty"`
	Recipient  PaymentsConfirmInfoShippingAddressRecipientResponse `json:"recipient,omitempty"`
}

type PaymentsConfirmInfoShippingAddressRecipientResponse struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	FirstNameOptional string `json:"firstNameOptional"`
	LastNameOptional  string `json:"lastNameOptional"`
	Email             string `json:"email"`
	PhoneNo           string `json:"phoneNo"`
}

func (client *Client) PaymentsConfirm(ctx context.Context, transactionId int64, request *PaymentsConfirmRequest) (response *PaymentsConfirmResponse, err error) {

	body, err := json.Marshal(request)
	res, err := client.post(ctx, fmt.Sprintf(endpointV3PaymentsConfirm, transactionId), body)
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
