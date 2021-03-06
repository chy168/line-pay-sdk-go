package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	PaymentsDetailsFieldsTransaction string = "RANSACTION"
	PaymentsDetailsFieldsOrder       string = "ORDER"
	PaymentsDetailsFieldsDefault     string = "ALL"
)

// if assign `TransactionIDs` and `OrderIDs` both at the same time, they should mean for the same record (like `AND` query).
type PaymentsDetailsRequest struct {
	TransactionIDs []int64  ``
	OrderIDs       []string ``
	Fields         string   ``
}

type PaymentsDetailsResponse struct {
	ReturnCode    string                        `json:"returnCode"`
	ReturnMessage string                        `json:"returnMessage"`
	Info          []PaymentsDetailsInfoResponse `json:"info"`
}

type PaymentsDetailsInfoResponse struct {
	TransactionID           int64                                   `json:"transactionId"`
	TransactionDate         time.Time                               `json:"transactionDate"`
	TransactionType         string                                  `json:"transactionType"`
	PayStatus               string                                  `json:"payStatus"` // AUTHORIZATION, VOIDED_AUTHORIZATION, EXPIRED_AUTHORIZATION
	ProductName             string                                  `json:"productName"`
	MerchantName            string                                  `json:"merchantName"`
	Currency                string                                  `json:"currency"`
	AuthorizationExpireDate time.Time                               `json:"authorizationExpireDate"`
	PayInfo                 []PaymentsDetailsInfoPayInfoResponse    `json:"payInfo"`
	RefundList              []PaymentsDetailsInfoRefundListResponse `json:"refundList"` // in case of `Transaction` type
	OriginalTransactionID   int64                                   `json:"originalTransactionId"`
	Packages                []PaymentsDetailsInfoPackagesResponse   `json:"packages"`
	Shipping                PaymentsDetailsInfoShippingResponse     `json:"shipping"`
}

type PaymentsDetailsInfoPayInfoResponse struct {
	Method string `json:"method"` // CREDIT_CARD, BALANCE, DISCOUNT
	Amount int    `json:"amount"` // sum(info[].payInfo[].amount) – sum(refundList[].refundAmount)
}

type PaymentsDetailsInfoRefundListResponse struct {
	RefundTransactionID   int64     `json:"refundTransactionId"`
	TransactionType       string    `json:"transactionType"` // PAYMENT_REFUND, PARTIAL_REFUND
	RefundAmount          int       `json:"refundAmount"`
	RefundTransactionDate time.Time `json:"refundTransactionDate"`
}

type PaymentsDetailsInfoPackagesResponse struct {
	ID            string                                        `json:"id"`
	Amount        int                                           `json:"amount"`
	UserFeeAmount int                                           `json:"userFeeAmount"`
	Name          string                                        `json:"name"`
	Products      []PaymentsDetailsInfoPackagesProductsResponse `json:"products"`
}

type PaymentsDetailsInfoPackagesProductsResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ImageURL      string `json:"imageUrl"`
	Quantity      int    `json:"quantity"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"originalPrice"`
}

type PaymentsDetailsInfoShippingResponse struct {
	MethodID  string                                     `json:"methodId"`
	FeeAmount int                                        `json:"feeAmount"`
	Address   PaymentsDetailsInfoShippingAddressResponse `json:"address"`
}

type PaymentsDetailsInfoShippingAddressResponse struct {
	Country    string                                              `json:"country"`
	PostalCode string                                              `json:"postalCode"`
	State      string                                              `json:"state"`
	City       string                                              `json:"city"`
	Detail     string                                              `json:"detail"`
	Optional   string                                              `json:"optional"`
	Recipient  PaymentsDetailsInfoShippingAddressRecipientResponse `json:"recipient"`
}

type PaymentsDetailsInfoShippingAddressRecipientResponse struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	FirstNameOptional string `json:"firstNameOptional"`
	LastNameOptional  string `json:"lastNameOptional"`
	Email             string `json:"email"`
	PhoneNo           string `json:"phoneNo"`
}

// PaymentsDetails
func (client *Client) PaymentsDetails(ctx context.Context, request *PaymentsDetailsRequest) (response *PaymentsDetailsResponse, err error) {

	params := url.Values{}

	for _, u := range request.TransactionIDs {
		params.Add("transactionId", strconv.FormatInt(u, 10))
	}

	for _, u := range request.OrderIDs {
		params.Add("orderId", u)
	}

	res, err := client.get(context.Background(), endpointV3PaymentsDetails, &params)
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

		response = &PaymentsDetailsResponse{}
		if err = json.Unmarshal(bodyBytes, response); err != nil {
			return
		}

	} else {
		err = fmt.Errorf("failed response, StatusCode: %d", res.StatusCode)
		return
	}

	return
}
