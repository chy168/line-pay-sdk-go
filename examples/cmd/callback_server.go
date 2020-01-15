package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	linepay "github.com/chy168/line-pay-sdk-go"
	"github.com/sirupsen/logrus"
)

var (
	ChannelID     *string = flag.String("channel-id", "", "Channel ID")
	ChannelSecret *string = flag.String("channel-secret", "", "Channel ID")
)

func main() {

	flag.Parse()

	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello: %q", html.EscapeString(r.URL.Path))

		transactionStringID := r.URL.Query().Get("transactionId")
		transactionID, err := strconv.ParseInt(transactionStringID, 10, 64)
		if err != nil {
			logrus.Errorf("transactionId is nil, err: %s", err.Error())
			return
		}

		orderID := r.URL.Query().Get("orderId")

		logrus.Infof("txid: '%d', orderid: '%s'", transactionID, orderID)

		client, err := linepay.NewClient(*ChannelID, *ChannelSecret, &linepay.Signer{ChannelId: *ChannelID}, &linepay.ClientOpts{})
		if err != nil {
			logrus.Errorf("init linepay client error: %s", err.Error())
			return
		}

		// Get Detail

		dataDetail1 := linepay.PaymentsDetailsRequest{
			TransactionIDs: []int64{transactionID},
			// OrderIDs:       []string{"order_0be9807d-88cf-42fe-bf69-75a51f1ad83f", "order_9583d466-6c47-488b-813f-894c0a26d7e8", "order_d776f2dd-eb7a-4611-b8cc-53242b9d7e71"},
			Fields: linepay.PaymentsDetailsFieldsTransaction,
		}

		resp1, err := client.PaymentsDetails(context.Background(), &dataDetail1)
		if err != nil {
			logrus.Errorf("Test PaymentsDetails failed: %s", err.Error())
			return
		}
		logrus.Infof("dump Detail#1 body: %+v", resp1)

		// Confirm
		data := linepay.PaymentsConfirmRequest{
			Amount:   100,
			Currency: "TWD",
		}

		resp, err := client.PaymentsConfirm(context.Background(), transactionID, &data)
		if err != nil {
			logrus.Errorf("client PaymentsConfirm error: %s", err.Error())
			return
		}
		logrus.Infof("dump PaymentsConfirm body: %+v", resp)

		// Detail #2
		resp2, err := client.PaymentsDetails(context.Background(), &dataDetail1)
		if err != nil {
			logrus.Errorf("Test PaymentsDetails #2 failed: %s", err.Error())
			return
		}

		logrus.Infof("PayStatus: %s", resp2.Info[0].PayStatus)
		logrus.Infof("dump Detail#2 body: %+v", resp2)

		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":9876", nil))

}
