package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/yeyee2901/lord-bidoof-bot/pkg/datasource"

	"github.com/rs/zerolog/log"
)

type TelegramService struct {
	*datasource.DataSource
}

func NewTelegramService(ds *datasource.DataSource) *TelegramService {
	return &TelegramService{ds}
}

func (t *TelegramService) GetBotStatus(ctx context.Context) (*RespGetMe, error) {
	var (
		msgType = "telegram.GetBotStatus"
		path    = "/getMe"
		method  = "GET"
	)

	resp := new(RespGetMe)
	httpStatus, err := t.sendToTelegram(ctx, method, path, nil, nil, resp)
	if err != nil {
		log.Error().Err(err).Msg(msgType + "-sendToTelegram()")
		return nil, err
	}

	// normal logging
	loggingDone := make(chan struct{})
	go func(done chan<- struct{}) {
		defer func() {
			done <- struct{}{}
		}()

		log.Info().
			AnErr("error", err).
			Interface("header", nil).
			Interface("req", nil).
			Interface("resp", resp).
			Str("endpoint", path).
			Str("method", method).
			Msg(msgType)
	}(loggingDone)

	defer func() {
		<-loggingDone
	}()

	// handle http status
	if httpStatus != http.StatusOK {
		if len(resp.Description) > 0 && !resp.Ok {
			err = fmt.Errorf("API telegram returned HTTP %d: %s", httpStatus, resp.Description)
			log.Error().Err(err).Msg(msgType + "-api.telegram.org")
			return nil, err
		} else {
			err = fmt.Errorf("Unknown error with HTTP %d", httpStatus)
			log.Error().Err(err).Msg(msgType + "unknown")
			return nil, err
		}
	}

	return resp, nil
}

func (t *TelegramService) sendToTelegram(
	ctx context.Context,
	method string,
	endpoint string,
	header map[string]string,
	payload any,
	resp any,
) (httpStatus int, err error) {
	client := new(http.Client)

	// create context for this http request (with timeout)
	tgContext, cancel := context.WithTimeout(ctx, time.Duration(t.Config.Telegram.RequestTimeout)*time.Second)
	defer cancel()

	// get bot token
	token, err := t.getBotToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// create the request
	body := new(bytes.Buffer)
	if err = json.NewEncoder(body).Encode(payload); err != nil {
		return http.StatusInternalServerError, err
	}

	urlEndpoint := t.Config.Telegram.Url + "/bot" + token + endpoint
	httpReq, err := http.NewRequestWithContext(tgContext, method, urlEndpoint, body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// read response body
	if resp != nil {
		if err = json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return httpResp.StatusCode, nil
}

func (t *TelegramService) getBotToken() (string, error) {
	token := os.Getenv(t.Config.Telegram.TokenEnv)
	if len(token) == 0 {
		return "", fmt.Errorf("Empty token in environment")
	}

	return token, nil
}
