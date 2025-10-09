package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/YuanJey/goutils2/pkg/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func Get(url, token string, req interface{}, resp interface{}) error {
	body := strings.NewReader("")
	if req != nil {
		jsonStr, err := json.Marshal(req)
		if err != nil {
			return err
		}
		body = strings.NewReader(string(jsonStr))
	}
	request, err := http.NewRequest("GET", url, body)
	if err != nil {
		return err
	}
	request.Header.Add("token", token)
	request.Header.Add("x-trace-id", utils.OperationIDGenerator())
	client := http.Client{Timeout: 5 * time.Second}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}
	result, err := io.ReadAll(httpResponse.Body)
	if httpResponse.StatusCode != 200 {
		return utils.Wrap(errors.New(httpResponse.Status), "status code failed "+url+string(result))
	}
	err = utils.JsonStringToStruct(string(result), &resp)
	if err != nil {
		return err
	}
	return nil
}
func Post(url, token string, data interface{}, resp interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Add("x-trace-id", utils.OperationIDGenerator())
	request.Header.Add("token", token)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//request.Header.Add("Cookie", "PHPSESSID="+ssid.Get()+";think_language=zh-CN")
	client := http.Client{Timeout: 10 * time.Second}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}
	result, err := io.ReadAll(httpResponse.Body)
	if httpResponse.StatusCode != 200 {
		log.Printf("api request err url is "+url, httpResponse.Status, string(result))
		return utils.Wrap(errors.New(httpResponse.Status), "status code failed "+url+string(result))
	}
	err = utils.JsonStringToStruct(string(result), &resp)
	if err != nil {
		return err
	}
	return nil
}
