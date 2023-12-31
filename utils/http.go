package utils

import (
	"EmptyClassroom/logs"
	"bytes"
	"context"
	"io"
	"net/http"
)

func HttpPostJson(ctx context.Context, url string, jsonStr []byte) (int, http.Header, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.CtxError(ctx, "http post json error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := io.ReadAll(resp.Body)
	return statusCode, header, body, nil
}

func HttpPostForm(ctx context.Context, url string, data map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logs.CtxError(ctx, "http post form error: %v", err)
		return 0, nil, nil, err
	}
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.CtxError(ctx, "http post form error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := io.ReadAll(resp.Body)
	return statusCode, header, body, nil
}

func HttpGetWithHeader(ctx context.Context, url string, header map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.CtxError(ctx, "http get with header error: %v", err)
		return 0, nil, nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.CtxError(ctx, "http get with header error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)
	return statusCode, resp.Header, body, nil
}
