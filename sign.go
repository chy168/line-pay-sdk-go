package linepay

import (
	"net/http"

	"crypto/hmac"
	"crypto/sha256"

	"encoding/base64"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Signer struct {
	ChannelId string
}

// SignWithBody implements API Authentication in `https://pay.line.me/developers/apis/onlineApis`
// HTTP Method : POST
// Signature = Base64(HMAC-SHA256(Your ChannelSecret, (Your ChannelSecret + URL Path + RequestBody + nonce)))
//
// HTTP Method : GET
// Signature = Base64(HMAC-SHA256(Your ChannelSecret, (Your ChannelSecret + URL Path + Query String + nonce))) Query String : A query string except ? (Example: Name1=Value1&Name2=Value2...)
func (v3 Signer) SignWithBody(r *http.Request, channelSecret string, requestBody string) (header http.Header, err error) {

	logrus.Debugf("sign with body dump body: %s", requestBody)

	myid, err := uuid.NewRandom()
	if err != nil {
		return
	}

	nonce := myid.String()

	sign := channelSecret + r.URL.Path + requestBody + nonce
	logrus.Debugf("To sign '%s'", sign)

	encResult := calculate(channelSecret, sign)

	// TODO
	header = http.Header{}
	header.Add("X-LINE-ChannelId", v3.ChannelId)
	header.Add("X-LINE-Authorization-Nonce", nonce)
	header.Add("X-LINE-Authorization", encResult)

	return
}

func calculate(secret, body string) (enc string) {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))
	// logrus.Infof("sha256: %s", hex.EncodeToString(h.Sum(nil)))
	enc = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return
}
