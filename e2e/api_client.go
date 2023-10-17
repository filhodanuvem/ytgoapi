package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var logger = log.Default()

type ApiClient struct {
	baseUrl string
}

func NewApiClient() ApiClient {
	return ApiClient{
		baseUrl: "http://localhost:8080",
	}
}

func (api *ApiClient) Post(path string, data map[string]string) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	logger.Println("POST", url, payload)

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return nil, err
	}

	logger.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Get(path string) (*http.Response, error) {
	url := api.baseUrl + path

	logger.Println("GET", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	logger.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Put(path string, data map[string]string) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	logger.Println("PUT", url, payload)

	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Delete(path string) (*http.Response, error) {
	url := api.baseUrl + path

	logger.Println("DELETE", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) ParseBody(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Printf("BODY %s", body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
