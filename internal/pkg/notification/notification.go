package notification

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"dating/internal/app/config"
	"dating/internal/pkg/glog"
)

type Data struct {
	Feature string `json:"feature"`
	Body    string `json:"body"`
}
type NotificationPayLoad struct {
	RegistrationIds []string `json:"registration_ids"`
	Data            Data     `json:"data"`
}

func PushNotification(conf *config.Configs, payLoad []byte) {
	logger := glog.New().WithField("package", "notification")

	req, _ := http.NewRequest("POST", conf.Notification.Firebase.Url, bytes.NewBuffer(payLoad))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", conf.Notification.Firebase.Key))

	client := &http.Client{}
	resp, _ := client.Do(req)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed when push notification to firebase: %v", err)
	}
	logger.Errorf("Push notification to completed: %v", string(bytes))
	defer resp.Body.Close()

}
