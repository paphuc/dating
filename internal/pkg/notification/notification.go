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
	Content string `json:"content"`
	Avatar  string `json:"avatar"`
}
type Notification struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}
type NotificationPayLoad struct {
	RegistrationIds []string     `json:"registration_ids"`
	Data            Data         `json:"data"`
	Foreground      bool         `json:"foreground"`
	Notification    Notification `json:"notification"`
}

func PushNotification(conf *config.Configs, payLoad []byte, result chan<- error) {
	logger := glog.New().WithField("package", "notificationpkg")

	req, err := http.NewRequest("POST", conf.Notification.Firebase.Url, bytes.NewBuffer(payLoad))
	if err != nil {
		result <- err
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", conf.Notification.Firebase.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result <- err
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed when push notification to firebase: %v", err)
	}

	logger.Infof("Push notification to completed: %v", string(bytes))
	defer resp.Body.Close()
	result <- nil
	return
}
