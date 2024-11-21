package serve

import (
	"bytes"
	"encoding/json"
	"fksunoapi/cfg"
	"fksunoapi/common"
	"fksunoapi/models"
	"fmt"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ErrCodeRequestFailed   = 1001
	ErrCodeResponseInvalid = 1002
	ErrCodeJsonFailed      = 1003
	ErrCodeTimeout         = 1004
)

var (
	SessionExp int64
	Session    string
)

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func NewErrorResponse(errorCode int, errorMsg string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}
}

func NewErrorResponseWithError(errorCode int, err error) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  err.Error(),
	}
}

func GetSession(c string) string {
	fmt.Println("cookie1", c)

	//https://clerk.suno.com/v1/client?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3
	//_url := "https://" + cfg.Domain + "/v1/client?_clerk_js_version=4.73.3"
	_url := "https://" + cfg.Domain + "/v1/client?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3"
	//https://clerk.suno.com/v1/client?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "ajs_anonymous_id=e4825918-be5f-43f5-9f84-9d2e871fedda; __client="+c)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("GetSession failed, error: %v", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != 200 {
		log.Printf("GetSession failed, invalid status code: %d", res.StatusCode)
		return ""
	}

	body, _ := io.ReadAll(res.Body)

	//log.Printf("session", string(body))

	var data models.GetSessionData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("GetSession failed, json unmarshal error: %v", err)
		return ""
	}

	if len(data.Response.Sessions) > 0 {
		SessionExp = data.Response.Sessions[0].ExpireAt
	}

	if len(data.Response.Sessions) > 0 {
		return data.Response.Sessions[0].Id
	}

	return ""
}

func GetJwtToken(c string) (string, *ErrorResponse) {
	if time.Now().After(time.Unix(SessionExp/1000, 0)) {
		Session = GetSession(c)
	}

	Session = "sess_2lhAxcrswO5EcZgj5ghGDcRcA7F"

	log.Println("Session", Session)

	//https://clerk.suno.com/v1/client/sessions/sess_2lhAxcrswO5EcZgj5ghGDcRcA7F/tokens?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3
	_url := fmt.Sprintf("https://"+cfg.Domain+"/v1/client/sessions/%s/tokens?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3", Session)
	log.Println("_url", _url)
	method := "POST"

	req, err := fhttp.NewRequest(method, _url, nil)

	if err != nil {
		log.Printf("GetJwtToken failed, error: %v", err)
		return "", NewErrorResponse(ErrCodeRequestFailed, "create request failed")
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+c)
	req.Header.Add("origin", "https://suno.com")
	req.Header.Add("referer", "https://suno.com/")
	req.Header.Add("referer", "https://suno.com/")

	res, err := common.Client.Do(req)
	if err != nil {
		log.Printf("GetJwtToken failed, error: %v", err)
		return "", NewErrorResponse(ErrCodeRequestFailed, "send request failed")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, _ := io.ReadAll(res.Body)

	//log.Println(string(body))

	if res.StatusCode != 200 {
		log.Printf("GetJwtToken failed, invalid status code: %d, response: %s", res.StatusCode, string(body))
		return "", NewErrorResponse(ErrCodeResponseInvalid, "invalid response")
	}

	var data models.GetTokenData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("GetJwtToken failed, json unmarshal error: %v", err)
		return "", NewErrorResponse(ErrCodeJsonFailed, "parse response failed")
	}

	if len(data.Jwt) == 0 {
		log.Print("GetJwtToken failed, empty jwt token")
		return "", NewErrorResponse(ErrCodeResponseInvalid, "get empty jwt token")
	}

	return data.Jwt, nil
}

//jwt config

func GetJwtConfig(c string) (string, string) {
	fmt.Println("cookie1", c)
	_url := "https://" + cfg.Domain + "/v1/client?__clerk_api_version=2021-02-05&_clerk_js_version=5.34.3"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "ajs_anonymous_id=e4825918-be5f-43f5-9f84-9d2e871fedda; __client="+c)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("GetSession failed, error: %v", err)
		return "", ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != 200 {
		log.Printf("GetSession failed, invalid status code: %d", res.StatusCode)
		return "", ""
	}

	body, _ := io.ReadAll(res.Body)

	//log.Printf("session", string(body))

	var data models.GetSessionData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("GetSession failed, json unmarshal error: %v", err)
		return "", ""
	}

	if len(data.Response.Sessions) > 0 {
		SessionExp = data.Response.Sessions[0].ExpireAt
	}

	if len(data.Response.Sessions) > 0 {
		return data.Response.Sessions[0].Id, data.Response.Sessions[0].LastActiveToken.Jwt
	}

	return "", ""
}
func sendRequest(url, method, c string, data []byte) ([]byte, *ErrorResponse) {
	//jwt, errResp := GetJwtToken(c)

	session, jwt := GetJwtConfig(c)
	log.Println("jwt", jwt, "342342", session)

	//if errResp != nil {
	//	errMsg := fmt.Sprintf("error getting JWT: %s", errResp.ErrorMsg)
	//	log.Printf("sendRequest failed, %s", errMsg)
	//	return nil, NewErrorResponse(errResp.ErrorCode, errMsg)
	//}

	//client := &http.Client{}
	var req *fhttp.Request
	var err error
	if data != nil {
		req, err = fhttp.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, err = fhttp.NewRequest(method, url, nil)
	}

	if err != nil {
		log.Printf("sendRequest failed123111, error creating request: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeRequestFailed, err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0")
	req.Header.Add("Authorization", "Bearer "+jwt)
	//req.Header.Add("Origin", "https://suno.com")
	req.Header.Add("Referer", "https://suno.com")
	req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
	//req.Header.Add("Priority", "u=1, i")

	res, err := common.Client.Do(req)

	if err != nil {
		log.Printf("sendRequest failed2222222, error sending request: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeRequestFailed, err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, _ := io.ReadAll(res.Body)
	log.Println("bodybodybodybodybodybody", string(body))

	if res.StatusCode != 200 {
		log.Printf("sendRequest failed55555555, unexpected status code: %d, response body: %s", res.StatusCode, string(body))
		return body, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("unexpected status code: %d, response body: %s", res.StatusCode, string(body)))
	}

	return body, nil
}

func V2Generate(d map[string]interface{}, c string) ([]byte, *ErrorResponse) {
	//https://studio-api.suno.ai/api/generate/v2/
	//https://studio-api.prod.suno.com/api/generate/v2/
	//_url := "https://studio-api.suno.ai/api/generate/v2/"
	_url := "https://studio-api.prod.suno.com/api/generate/v2/"
	jsonData, err := json.Marshal(d)

	//生成语音
	log.Println("jsonData _url", _url)
	log.Println("jsonData", string(jsonData))

	if err != nil {
		log.Printf("V2Generate failed, error marshalling request data: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeJsonFailed, err)
	}
	body, errResp := sendRequest(_url, "POST", c, jsonData)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func V2GetFeedTask(ids, c string) ([]byte, *ErrorResponse) {
	ids = url.QueryEscape(ids)
	//https://studio-api.prod.suno.com/api/feed/v2?ids=ea0e897e-22fa-4b55-8876-56e8529c24a1%2C37c11f85-09be-4d84-a8be-63925d697376
	//https://studio-api.prod.suno.com/api/feed/v2?ids=e77fe186-4c9f-4192-a446-b3c80383ff80%2Ccbd2bf2a-34be-41c8-9692-22c4982eaf03&page=5000

	_url := "https://studio-api.prod.suno.com/api/feed/?ids=" + ids + "&page=5000"

	//_url := "https://studio-api.suno.ai/api/feed/?ids=" + ids

	body, errResp := sendRequest(_url, "GET", c, nil)

	//log.Println("body", string(body))
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func GenerateLyrics(d map[string]interface{}, c string) ([]byte, *ErrorResponse) {
	_url := "https://studio-api.prod.suno.com/api/generate/lyrics/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Printf("GenerateLyrics failed, error marshalling request data: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeJsonFailed, err)
	}
	body, errResp := sendRequest(_url, "POST", c, jsonData)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func GetLyricsTask(ids, c string) ([]byte, *ErrorResponse) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/" + ids
	body, errResp := sendRequest(_url, "GET", c, nil)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func CheckSong(c string) (bool, *ErrorResponse) {
	var data struct {
		Required bool `json:"required"`
	}

	Bytes, err := sendRequest("https://studio-api.prod.suno.com/api/c/check", "POST", c, []byte(`{"ctype":"generation"}`))

	log.Println("Bytes", string(Bytes))

	if err != nil {
		return false, err
	}

	_ = json.Unmarshal(Bytes, &data)
	if data.Required == true {
		return true, nil
	}

	return false, nil
}

func SunoChat(c map[string]interface{}, ck string) (interface{}, *ErrorResponse) {

	check, err := CheckSong(ck)

	if err != nil {
		return nil, err
	}

	TokenCaptcha := ""
	if check {
		TokenCaptcha, err = CaptchaHandlers()
		if err != nil {
			return nil, err
		}
	}
	lastUserContent := GetLastUserContent(c)
	uid := uuid.NewString()
	d := map[string]interface{}{
		"mv":                     "chirp-v3-5",
		"gpt_description_prompt": lastUserContent,
		"prompt":                 "",
		"make_instrumental":      false,
		"token":                  TokenCaptcha,
		"generation_type":        "TEXT",
		"metadata": map[string]interface{}{
			"create_session_token": uid,
		},
	}

	body, errResp := V2Generate(d, ck)
	log.Println("8888888", string(body), errResp)

	if errResp != nil {
		return nil, errResp
	}

	var v2GenerateData models.GenerateData
	if err := json.Unmarshal(body, &v2GenerateData); err != nil {
		log.Printf("SunoChat failed, error unmarshalling generate data: %v, response body: %s", err, string(body))
		return nil, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("parse generate data failed, response body: %s", string(body)))
	}

	log.Println("打点", 55555555, v2GenerateData.Clips)

	clipIds := make([]string, len(v2GenerateData.Clips))
	for i, clip := range v2GenerateData.Clips {
		clipIds[i] = clip.Id
	}
	ids := strings.Join(clipIds, ",")

	timeout := time.After(3 * time.Minute)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			return nil, NewErrorResponse(ErrCodeTimeout, "get feed task timeout")
		case <-tick:
			body, errResp = V2GetFeedTask(ids, ck)
			if errResp != nil {
				return nil, errResp
			}

			var v2GetFeedData []map[string]interface{}
			if err := json.Unmarshal(body, &v2GetFeedData); err != nil {
				log.Printf("SunoChat failed, error unmarshalling feed data: %v, response body: %s", err, string(body))
				return nil, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("parse feed data failed, response body: %s", string(body)))
			}

			allComplete := true
			for _, data := range v2GetFeedData {
				if data["status"] != "complete" {
					allComplete = false
					break
				}
			}

			if allComplete {
				var markdown strings.Builder
				markdown.WriteString(fmt.Sprintf("# %s\n\n", v2GetFeedData[0]["title"]))
				markdown.WriteString(fmt.Sprintf("%s\n\n", v2GetFeedData[0]["metadata"].(map[string]interface{})["prompt"]))
				markdown.WriteString(fmt.Sprintf("## 版本一\n\n"))
				markdown.WriteString(fmt.Sprintf("视频链接：%s\n\n", v2GetFeedData[0]["video_url"]))
				markdown.WriteString(fmt.Sprintf("音频链接：%s\n\n", v2GetFeedData[0]["audio_url"]))
				markdown.WriteString(fmt.Sprintf("图片链接：%s\n\n", v2GetFeedData[0]["image_large_url"]))
				markdown.WriteString(fmt.Sprintf("## 版本二\n\n"))
				markdown.WriteString(fmt.Sprintf("视频链接：%s\n\n", v2GetFeedData[1]["video_url"]))
				markdown.WriteString(fmt.Sprintf("音频链接：%s\n\n", v2GetFeedData[1]["audio_url"]))
				markdown.WriteString(fmt.Sprintf("图片链接：%s\n\n", v2GetFeedData[1]["image_large_url"]))

				response := map[string]interface{}{
					"choices": []map[string]interface{}{
						{
							"finish_reason": "stop",
							"index":         0,
							"message": map[string]string{
								"content": markdown.String(),
								"role":    "assistant",
							},
							"logprobs": nil,
						},
					},
					"created": time.Now().Unix(),
					"id":      "chatcmpl-7QyqpwdfhqwajicIEznoc6Q47XAyW",
					"model":   c["model"].(string),
					"object":  "chat.completion",
					"usage": map[string]int{
						"completion_tokens": 17,
						"prompt_tokens":     57,
						"total_tokens":      74,
					},
				}

				return response, nil
			}
		}
	}
}
