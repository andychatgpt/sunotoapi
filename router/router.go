package router

import (
	"fksunoapi/cfg"
	"fksunoapi/serve"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func CreateTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		/**
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
		*/

		//sessionId, err := serve.GetSessionS()
		//if err != nil {
		//	return c.Status(fiber.StatusInternalServerError).JSON(serve.NewErrorResponse(200, "error 901"))
		//}
		//log.Println("sessionId:", sessionId)

		var data map[string]interface{}

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		ck := c.Get("Authorization")
		log.Printf("ck", ck)

		//ck = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsaWVudF8ycDc2cDZZQ09QRHNtQVNab0w4bTlWczFnMnoiLCJyb3RhdGluZ190b2tlbiI6IjVnbTh3bncyY3VxbTdxZ3h4amppYjZvcWxmbzRscW00ZHpteDE3YWQifQ.W26zyWRNiRCLLNlbFalv4JybeQXjxqITQ8m_RIZdaaiDnIY9taWona7RWDQn36qaH7-5cvPeAnswwZwY-awvG-SQbbuqBHKOo2o1qb7tx5-PQgAYXFkH16k5W-srkfY-bLdlTSM3evGGir1uE16IwToTMgFwfcLO8CoyOv8xF3l_V1VYoGz9NrTsNcuvShEdg5qs7jyknk74zincb3mGtIOJ5oeLK9EfkzT_HGhvFd0SZ4IGGPxp6cawhPqrIXvY3WXjHjEcrjQWah-QPywL-Vdqr5QYjKOJZwf0DNtXyCD6T2zvreHBSeNcsgIsMslrROgXJLxE6VP2k7eWnHkDDA"

		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}

		serve.Session = serve.GetSession(ck)

		var body []byte
		var errResp *serve.ErrorResponse

		if c.Path() == "/v2/generate" {
			check, errResp := serve.CheckSong(ck)

			if errResp != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errResp)
			}

			TokenCaptcha := ""
			if check {
				TokenCaptcha, errResp = serve.CaptchaHandlers()
				if errResp != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errResp)
				}
			}

			//lastUserContent := serve.GetLastUserContent(data)
			//log.Println("show1", lastUserContent)
			//uid := uuid.NewString()
			data["token"] = TokenCaptcha
			data["generation_type"] = "TEXT"
			//sessionId := uuid.NewString()
			//if sessionId != "" {
			data["metadata"] = map[string]interface{}{
				"lyrics_model": "default",
			}

			//}

			//if _, ok := data["artist_clip_id"]; !ok {
			//	data["artist_clip_id"] = nil
			//}
			//
			//if _, ok := data["artist_end_s"]; !ok {
			//	data["artist_end_s"] = nil
			//}

			//if _, ok := data["artist_start_s"]; !ok {
			//	data["artist_start_s"] = nil
			//}
			//
			//if _, ok := data["continue_at"]; !ok {
			//	data["continue_at"] = nil
			//}

			//if _, ok := data["continue_clip_id"]; !ok {
			//	data["continue_clip_id"] = nil
			//}
			//
			//if _, ok := data["continued_aligned_prompt"]; !ok {
			//	data["continued_aligned_prompt"] = nil
			//}
			//
			//if _, ok := data["infill_end_s"]; !ok {
			//	data["infill_end_s"] = nil
			//}
			//
			//if _, ok := data["infill_start_s"]; !ok {
			//	data["infill_start_s"] = nil
			//}
			//
			//if _, ok := data["negative_tags"]; !ok {
			//	data["negative_tags"] = ""
			//}
			//
			//if _, ok := data["tags"]; !ok {
			//	data["tags"] = ""
			//}
			//
			//if _, ok := data["task"]; !ok {
			//	data["task"] = nil
			//}
			//
			//if _, ok := data["title"]; !ok {
			//	data["title"] = ""
			//}

			data["user_uploaded images b64"] = []string{}

			body, errResp = serve.V2Generate(data, ck)
		} else if c.Path() == "/v2/lyrics/create" {
			body, errResp = serve.GenerateLyrics(data, ck)
		}

		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		return c.Status(fiber.StatusOK).Send(body)
	}
}

func GetTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		if len(data["ids"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeRequestFailed, "Cannot find ids"))
		}
		ck := c.Get("Authorization")

		//ck = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsaWVudF8ycDc2cDZZQ09QRHNtQVNab0w4bTlWczFnMnoiLCJyb3RhdGluZ190b2tlbiI6IjVnbTh3bncyY3VxbTdxZ3h4amppYjZvcWxmbzRscW00ZHpteDE3YWQifQ.W26zyWRNiRCLLNlbFalv4JybeQXjxqITQ8m_RIZdaaiDnIY9taWona7RWDQn36qaH7-5cvPeAnswwZwY-awvG-SQbbuqBHKOo2o1qb7tx5-PQgAYXFkH16k5W-srkfY-bLdlTSM3evGGir1uE16IwToTMgFwfcLO8CoyOv8xF3l_V1VYoGz9NrTsNcuvShEdg5qs7jyknk74zincb3mGtIOJ5oeLK9EfkzT_HGhvFd0SZ4IGGPxp6cawhPqrIXvY3WXjHjEcrjQWah-QPywL-Vdqr5QYjKOJZwf0DNtXyCD6T2zvreHBSeNcsgIsMslrROgXJLxE6VP2k7eWnHkDDA"

		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}

		serve.Session = serve.GetSession(ck)
		var body []byte
		var errResp *serve.ErrorResponse
		if c.Path() == "/v2/feed" {
			log.Println("data1111", data)
			body, errResp = serve.V2GetFeedTask(data["ids"], ck)
		} else if c.Path() == "/v2/lyrics/task" {
			body, errResp = serve.GetLyricsTask(data["ids"], ck)
		}
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func SunoChat() fiber.Handler {

	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		ck := c.Get("Authorization")
		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}
		serve.Session = serve.GetSession(ck)
		fmt.Println("serve.Session101", serve.Session)
		res, errResp := serve.SunoChat(data, ck)
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.ConfigDefault))
	app.Use(cors.New(cors.ConfigDefault))
	app.Post("/v2/generate", CreateTask())
	app.Post("/v2/feed", GetTask())
	app.Post("/v2/lyrics/create", CreateTask())
	app.Post("/v2/lyrics/task", GetTask())
	app.Post("/v1/chat/completions", SunoChat())
}
