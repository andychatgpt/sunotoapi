package serve

import (
	"encoding/json"
	"fksunoapi/cfg"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Rsp struct {
	Token string
}

func CaptchaHandlers() (string, *ErrorResponse) {
	url := "https://api.acedata.cloud/captcha/token/hcaptcha"

	payload := strings.NewReader(`{
	"website_key": "d65453de-3f1a-4aac-9366-a0f06e52b2ce",
	"website_url": "https://suno.com/create"
	}`)

	client := &http.Client{}
	req, err1 := http.NewRequest("POST", url, payload)

	if err1 != nil {
		return "", NewErrorResponse(ErrCodeTimeout, err1.Error())
	}

	log.Println("cfg.Config.Auth.Captcha", cfg.Config.Auth.Captcha)

	req.Header.Add("authorization", "Bearer "+cfg.Config.Auth.Captcha)
	req.Header.Add("content-type", "application/json")

	var res, err = client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", NewErrorResponse(ErrCodeTimeout, err.Error())
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)

	log.Printf(string(body))

	if err != nil {
		return "", nil
	}

	if res.StatusCode != 200 {
		return "", NewErrorResponse(ErrCodeTimeout, "CaptchaHandlers code error")
	}

	var rps Rsp

	err = json.Unmarshal(body, &rps)
	if rps.Token == "" {
		return "", NewErrorResponse(ErrCodeTimeout, "CaptchaHandlers code 101 error")

	}
	return rps.Token, nil
}
