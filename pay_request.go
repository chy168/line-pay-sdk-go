package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// `Amount` required, valid amount `form.amount != sum(packages[].amount) + sum(packages[].userFee) + shippingFee`
// `Currency` required, is ISO 4217, supported: USD, JPY, TWD, THB
// `OrderId` required
// if `Capture` true: only need to call `Confirm API` to process payments. false: call `Confirm API` and then `Capture API`
type PaymentsRequest struct {
	Amount       int                         `json:"amount"`
	Currency     string                      `json:"currency"`
	OrderID      string                      `json:"orderId"`
	Packages     []PaymentsPackageRequest    `json:"packages"`
	RedirectUrls PaymentsRedirectUrlsRequest `json:"redirectUrls"`
	Options      PaymentsOptionsRequest      `json:"options"`
}

// `Id` required
// `Amount` required, valid amount `packages[].amount != sum(packages[].products[].quantity * packages[].products[].price)`
// `Name` required
type PaymentsPackageRequest struct {
	ID       string                          `json:"id"`
	Amount   int                             `json:"amount"`
	UserFee  int                             `json:"userFee,omitempty"`
	Name     string                          `json:"name"`
	Products []PaymentsPackageProductRequest `json:"products"`
}

// `Name` required
// `Quantity` required
// `Price` required
type PaymentsPackageProductRequest struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	ImageUrl      string `json:"imageUrl,omitempty"`
	Quantity      int    `json:"quantity"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"originalPrice,omitempty"`
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
// `ConfirmUrl` required
// `CancelUrl` required
type PaymentsRedirectUrlsRequest struct {
	AppPackageName string `json:"appPackageName,omitempty"`
	ConfirmUrl     string `json:"confirmUrl"`
	ConfirmUrlType string `json:"confirmUrlType,omitempty"`
	CancelUrl      string `json:"cancelUrl"`
}

type PaymentsOptionsRequest struct {
	Payment       PaymentsOptionsPaymentRequest       `json:"payment"`
	Display       PaymentsOptionsDisplayRequest       `json:"display"`
	Shipping      PaymentsOptionsShippingRequest      `json:"shipping"`
	FamilyService PaymentsOptionsFamilyServiceRequest `json:"familyService"`
	Extra         PaymentsOptionsExtraRequest         `json:"extra"`
}

type PaymentsOptionsPaymentRequest struct {
	Capture bool   `json:"capture,omitempty"`
	PayType string `json:"payType,omitempty"` // NORMAL, PREAPPROVED
}

type PaymentsOptionsDisplayRequest struct {
	Locale                 string `json:"locale,omitempty"` // en, ja, ko, th, zh_TW, zh_CN
	CheckConfirmUrlBrowser bool   `json:"checkConfirmUrlBrowser,omitempty"`
}

type PaymentsOptionsShippingRequest struct {
	ShippintType   string                                `json:"type,omitempty"`      // NO_SHIPPING, FIXED_ADDRESS, SHIPPING
	FeeAmount      string                                `json:"feeAmount,omitempty"` //why string?
	FeeInquiryUrl  string                                `json:"feeInquiryUrl,omitempty"`
	FeeInquiryType string                                `json:"feeInquiryType,omitempty"` // CONDITION, FIXED
	Address        PaymentsOptionsShippingAddressRequest `json:"address"`
}

type PaymentsOptionsShippingAddressRequest struct {
	Country    string                                         `json:"country,omitempty"`
	PostalCode string                                         `json:"postalCode,omitempty"`
	State      string                                         `json:"state,omitempty"`
	City       string                                         `json:"city,omitempty"`
	Detail     string                                         `json:"detail,omitempty"`
	Optional   string                                         `json:"optional,omitempty"`
	Recipient  PaymentsOptionsShippingAddressRecipientRequest `json:"recipient,omitempty"`
}

type PaymentsOptionsShippingAddressRecipientRequest struct {
	FirstName         string `json:"firstName,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	FirstNameOptional string `json:"firstNameOptional,omitempty"`
	LastNameOptional  string `json:"lastNameOptional,omitempty"`
	Email             string `json:"email,omitempty"`
	PhoneNo           string `json:"phoneNo,omitempty"`
}

type PaymentsOptionsFamilyServiceRequest struct {
	AddFriends []PaymentsOptionsFamilyServiceAddFriendsRequest `json:"addFriends"`
}

type PaymentsOptionsFamilyServiceAddFriendsRequest struct {
	AddType string   `json:"type,omitempty"` // line@
	IDs     []string `json:"ids,omitempty"`
}

type PaymentsOptionsExtraRequest struct {
	BranchName string `json:"branchName,omitempty"`
	BranchId   string `json:"branchId,omitempty"`
}

// response
type PaymentsResponse struct {
	ReturnCode    string               `json:"returnCode"`
	ReturnMessage string               `json:"returnMessage"`
	Info          PaymentsInfoResponse `json:"info"`
}

type PaymentsInfoResponse struct {
	TransactionID      int64                          `json:"transactionId"`
	PaymentAccessToken string                         `json:"paymentAccessToken"`
	PaymentURL         PaymentsInfoPaymentURLResponse `json:"paymentUrl"`
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
