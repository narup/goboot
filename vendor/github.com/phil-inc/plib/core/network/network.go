package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 60,
}

// HTTPGet - makes a get request to the given URL and HTTP headers.
// it returns response data byte or error
func HTTPGet(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//add headers!
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	log.Printf("GET request to url: %s\n", url)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, fmt.Errorf("Http response NOT_OK. Status: %s, Code:%d", res.Status, res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}

// HTTPPost makes a POST data to the given url with headers
func HTTPPost(url string, values url.Values, headers map[string]string) ([]byte, error) {
	rb := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", url, rb)

	//add headers!
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("POST request to url: %s\n", url)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ioutil.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("Http response NOT_OK. Status: %s, Code:%d", resp.Status, resp.StatusCode)
}

// HTPPJsonPost - sends JSON string data as post request
func HTPPJsonPost(url, jsonBody string, headers map[string]string) ([]byte, error) {
	return httpJSONSend(url, "POST", jsonBody, headers)
}

// HTPPJsonPut - sends JSON string data as put request
func HTPPJsonPut(url, jsonBody string, headers map[string]string) ([]byte, error) {
	return httpJSONSend(url, "PUT", jsonBody, headers)
}

func httpJSONSend(url, method, jsonBody string, headers map[string]string) ([]byte, error) {
	reader := strings.NewReader(jsonBody)
	req, err := http.NewRequest(method, url, reader)
	//add headers!
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	log.Printf("%s request to url: %s\n", method, url)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {

		errResp := fmt.Sprintf("Http response NOT_OK. Status: %s, Code:%d", res.Status, res.StatusCode)
		if res.Body != nil {
			resp, _ := ioutil.ReadAll(res.Body)
			errResp = errResp + fmt.Sprintf(", Body: %s", resp)
		}

		return nil, fmt.Errorf(errResp)
	}
	return ioutil.ReadAll(res.Body)
}
