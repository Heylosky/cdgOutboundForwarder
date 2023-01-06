package main

import (
	"encoding/json"
	"github.com/cdgProcessor/outboundForwarder/logger"
	"github.com/cdgProcessor/outboundForwarder/messageQ"
	"github.com/cdgProcessor/outboundForwarder/models"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	logger.InitLogger("./outboundForwarder.log")
	zap.L().Info("Outbound forwarder processor starting...")

	waitSendChan := make(chan []byte)
	sentChan := make(chan models.MbRc)

	go messageQ.MQRead(waitSendChan, "outboundSMS", "outboundSmsToCu", "outboundToCu")
	go messageQ.Publish(sentChan, "obSentRecordToDb")

	Processor(waitSendChan, sentChan)
}

func Processor(c <-chan []byte, sentChan chan<- models.MbRc) {
	var msg models.MbRc
	client := &http.Client{}

	for payload := range c {
		err := json.Unmarshal(payload, &msg)
		if err != nil {
			zap.L().Error("json unmarshal sms error, message field is uncorrect.", zap.Error(err))
		}
		time.Sleep(2 * time.Millisecond) //simulate message field replace time
		zap.L().Info("Replace some value in message", zap.Any("msg", msg))

		urlMap := url.Values{}
		urlMap.Add("payload", msg.Payload)

		request, _ := http.NewRequest("POST", "http://127.0.0.1:8081/test", strings.NewReader(urlMap.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(request)
		if err != nil {
			zap.L().Error("Client request err", zap.Error(err))
		} else if res.StatusCode != 200 {
			zap.L().Info("Failed to post message to customer api")
			res.Body.Close()
		} else {
			res.Body.Close()
		}

		sentChan <- msg
	}
}
