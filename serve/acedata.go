package serve

//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"fksunoapi/cfg"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/google/uuid"
//	"io"
//	"log"
//	"net/http"
//	"time"
//)
//
//type aceData struct {
//}
//
//var AceData aceData
//
////歌曲创建
//
//type GenerateData struct {
//	Id                string      `json:"id"`
//	Clips             interface{} `json:"clips"`
//	MajorModelVersion string      `json:"major_model_version"`
//	Status            string      `json:"status"`
//	CreatedAt         time.Time   `json:"created_at"`
//	BatchSize         int         `json:"batch_size"`
//}
//
////func (p *aceData) Create(title, prompt, lyric string) {
////
////	return c.Status(fiber.StatusOK).Send(body)
////
////}
//
//func CreateTask() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		//func (p *aceData) Create(title, prompt, lyric string)  {
//		url := "https://api.acedata.cloud/suno/audios"
//		custom := false
//		instrumentalboolean := true
//
//		if lyric != "" {
//			custom = true
//			instrumentalboolean = false
//		}
//
//		param := struct {
//			Action              string `json:"action"`
//			Prompt              string `json:"prompt"` // 提示词
//			Title               string `json:"title"`  // 标题
//			Model               string `json:"model"`
//			Lyric               string `json:"lyric"`               // 具体歌词
//			Custom              bool   `json:"custom"`              // true 通过具体歌词生成，false 提示词生成
//			InstrumentalBoolean bool   `json:"instrumentalboolean"` // true 忽略Lyric这个参数
//			CallbackUrl         string `json:"callback_url"`        // 如果这个参数不为空，就马上结束请求返回一个taskId,异步通知这个地址歌曲是否成功，或者通过taskId查询歌曲状态；如果这个参数为空 就一直等着歌曲生成结束，这个过程很长很长
//
//		}{
//			Action:              "generate",
//			Prompt:              prompt,
//			Model:               "chirp-v3-5",
//			Lyric:               lyric,
//			Custom:              custom,
//			InstrumentalBoolean: instrumentalboolean,
//			Title:               title,
//			CallbackUrl:         "http://13.214.179.63:8808/x/callback",
//		}
//
//		jsonData, _ := json.Marshal(param)
//		client := &http.Client{}
//		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
//
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
//		}
//		req.Header.Add("authorization", "Bearer e8e01042216c4e49a96c1b5371b2a495")
//		req.Header.Add("content-type", "application/json")
//
//		resp, err := client.Do(req)
//		if err != nil {
//			fmt.Println(err)
//			return "", err
//		}
//		defer func(Body io.ReadCloser) {
//			err := Body.Close()
//			if err != nil {
//
//			}
//		}(resp.Body)
//
//		body, err := io.ReadAll(resp.Body)
//		if err != nil {
//			fmt.Println(err)
//			return "", err
//		}
//		log.Print(string(body))
//
//		if resp.StatusCode != 200 {
//			return "", errors.New(string(body))
//		}
//
//		var data map[string]interface{}
//
//		err = json.Unmarshal(body, &data)
//
//		if err != nil {
//			return "", errors.New(string(body))
//		}
//
//		var generateData GenerateData
//		u := uuid.NewString()
//		generateData.Id = u
//		generateData.Status = "0"
//		generateData.CreatedAt = time.Now()
//		generateData.MajorModelVersion = "v3"
//		generateData.Clips = data
//
//		var data map[string]interface{}
//
//		if err := c.BodyParser(&data); err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
//		}
//
//
//
//
//
//		return c.Status(fiber.StatusOK).Send(body)
//}
