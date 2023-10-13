package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ApiClient struct {
	baseUrl string
}

func NewApiClient(baseUrl string) ApiClient {
	return ApiClient{baseUrl: baseUrl}
}

func (api *ApiClient) Post(path string, data map[string]string) *http.Response {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	logger.Println("POST", url, payload)

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		panic(err)
	}

	logger.Println("RESPONSE", resp.Status)

	return resp
}

func (api *ApiClient) Get(path string) *http.Response {
	url := api.baseUrl + path

	logger.Println("GET", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	logger.Println("RESPONSE", resp.Status)

	return resp
}

func (api *ApiClient) Put(path string, data map[string]string) *http.Response {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	logger.Println("PUT", url, payload)

	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Println("RESPONSE", resp.Status)

	return resp
}

func (api *ApiClient) Delete(path string) *http.Response {
	url := api.baseUrl + path

	logger.Println("DELETE", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Println("RESPONSE", resp.Status)

	return resp
}

func (api *ApiClient) ParseBody(resp *http.Response) map[string]interface{} {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Printf("BODY %s", body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return data
}
