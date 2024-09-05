package utils

import (
	"io"
	"net/http"

	"example.com/pkg/logger"
)

func Get(url string, headers map[string]string) (Response, error) {

	method := "GET"

	output := Response{Url: url, Method: method}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return output, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)

	if err != nil {
		output = Response{Url: url, Method: method, Status: 404, Data: "", Error: string(err.Error())}
		logger.E.Error().Str("url", url).Str("method", method).Int("status", 404).Err(err).Msg("GET")
		return output, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.E.Error().Str("url", url).Str("method", method).Int("status", resp.StatusCode).Err(err).Msg("HTTP")
		return output, err
	}

	output = Response{Url: url, Method: method, Status: resp.StatusCode, Data: string(body)}

	logger.I.Info().Str("url", url).Str("method", method).Int("status", resp.StatusCode).Msg("HTTP")

	return output, nil
}

