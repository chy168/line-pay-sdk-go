package linepay

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	APIHostSandbox    = "https://sandbox-api-pay.line.me"
	APIHostProduction = "https://api-pay.line.me"

	// POST /v3/payments/request
	endpointV3PaymentsRequest = "/v3/payments/request"

	// POST /v3/payments/{transactionId}/confirm
	endpointV3PaymentsConfirm = "/v3/payments/%d/confirm"

	// GET /v3/payments
	endpointV3PaymentsDetails = "/v3/payments"
)

type Client struct {
	channelID     string
	channelSecret string
	apiEndpoint   *url.URL
	httpClient    *http.Client
	signer        *Signer
}

type ClientOpts struct {
	ProductionEnabled bool
}

func NewClient(channelID, channelSecret string, signer *Signer, opts *ClientOpts) (*Client, error) {
	if channelSecret == "" || channelID == "" {
		return nil, errors.New("channel id or secret not correct")
	}

	apiEndpoint := APIHostSandbox
	if opts.ProductionEnabled {
		apiEndpoint = APIHostProduction
	}

	uu, err := url.ParseRequestURI(apiEndpoint)
	if err != nil {
		return nil, err
	}

	c := &Client{
		channelID:     channelID,
		channelSecret: channelSecret,
		apiEndpoint:   uu,
		httpClient:    http.DefaultClient,
		signer:        signer,
	}

	return c, nil
}

func (client *Client) post(ctx context.Context, endpoint string, body []byte) (res *http.Response, err error) {

	req, err := http.NewRequestWithContext(ctx, "POST", client.url(endpoint), bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("post request error: %s", err.Error())
		return
	}

	header, err := client.signer.SignWithBody(req, client.channelSecret, string(body))

	req.Header = header
	req.Header.Add("Content-Type", "application/json")

	res, err = client.do(ctx, req)
	return
}

func (client *Client) get(ctx context.Context, endpoint string, params *url.Values) (res *http.Response, err error) {

	myURL := client.url(endpoint)
	targetURL, err := url.ParseRequestURI(myURL)
	if err != nil {
		err = fmt.Errorf("error ParseRequestURI '%s'", myURL)
		return
	}
	targetURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL.String(), nil)
	if err != nil {
		err = fmt.Errorf("get request error: %s", err.Error())
		return
	}

	header, err := client.signer.SignWithBody(req, client.channelSecret, params.Encode())

	req.Header = header

	res, err = client.do(ctx, req)
	return
}

func (client *Client) url(endpoint string) string {
	u := *client.apiEndpoint
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (client *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {

	req.Header.Set("User-Agent", "line-pay-sdk-go")

	res, err := client.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return res, err

}
