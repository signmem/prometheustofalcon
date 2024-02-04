package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/signmem/prometheustofalcon/g"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
)

var (
	FalconAuthName = "xxx"
	FalconAuthSig = "xxxx"
)

func FalconToken() (string, error) {

	// crate falcon api header token access

	token, err := json.Marshal(map[string]string{"name": FalconAuthName,
		"sig": FalconAuthSig})

	if err != nil {
		return "", err
	}

	return  string(token), nil
}

func HttpApiPut(fullApiUrl string, jsonData []byte, tokenType string) (status bool, err error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, fullApiUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if tokenType == "falcon" {
		token, err := FalconToken()
		if err == nil {
			req.Header.Add("Apitoken", token)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	if  ( resp.StatusCode  == 200 ) {
		return true, nil
	} else {
		return false, errors.New("ttpApiPut() response not 200")
	}
}

func HttpApiGet(fullApiUrl string, params string, tokenType string) (io.ReadCloser, error) {

	var httpUrl string
	client := &http.Client{}
	if params != ""  {
		httpUrl = fullApiUrl + params
	} else {
		httpUrl = fullApiUrl
	}


	req, err := http.NewRequest("GET", httpUrl, nil)

	if err != nil {
		return nil, errors.New("HttpApiGet() http get error with NewRequest")
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if tokenType == "falcon" {
		token, err := FalconToken()
		if err == nil {
			req.Header.Add("Apitoken", token)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, err
	}

	if ( resp.StatusCode  == 200 ) {
		return resp.Body, nil
	} else {
		return nil, errors.New("HttpApiGet() resp status code not 200.")
	}

}

func HttpApiPost(fullApiUrl string, params []byte, tokenType string) (io.ReadCloser, error) {
	// use to access http post
	// params = post params  [must be []byte format]
	// return http response


	tr := &http.Transport{
		MaxIdleConns: 10,
		IdleConnTimeout: 10 * time.Second,
		DisableCompression: true,
	}

	req, err := http.NewRequest("POST", fullApiUrl, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	if tokenType == "falcon" {
		token, err := FalconToken()
		if err == nil {
			req.Header.Add("Apitoken", token)
		}
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	request := req.WithContext(ctx)
	if err != nil {
		return nil, errors.New("HttpApiPost() http post error with NewRequest")
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)

	if err != nil {
		return nil, errors.New("HttpApiPost()  client access error.")
	}

	defer cancelFunc()

	if resp.Body == nil {
		return nil, err
	}

	defer resp.Body.Close()

	if ( resp.StatusCode  == 200 ) {
		return resp.Body, nil
	} else {
		return nil, errors.New("HttpApiPost() resp status code not 200.")
	}
}

func HttpApiDelete(fullApiUrl string, params string, tokenType string) (io.ReadCloser, error) {
	// use to do http Delete request
	// METHOD: DELETE

	client := &http.Client{}
	httpUrl := fullApiUrl + params
	req, err := http.NewRequest("DELETE", httpUrl, nil)

	if err != nil {
		return nil, errors.New("HttpApiDelete() http delete error with NewRequest")
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if tokenType == "falcon" {
		token, err := FalconToken()
		if err == nil {
			req.Header.Add("Apitoken", token)
		}
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.New("HttpApiDelete() http delete error")
	}

	if ( resp.StatusCode  == 200 ) {
		return resp.Body, nil
	} else {
		return nil, errors.New("HttpApiDelete() resp status code not 200.")
	}
}



func HttpsApiGet(fullApiUrl string, params string) (io.ReadCloser, error) {

	caCert, _ := ioutil.ReadFile(g.Config().TLS.CaFile)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(g.Config().TLS.CertFile, g.Config().TLS.KeyFile)

	if err != nil {
		return nil, err
	}

	client := &http.Client {
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	var httpUrl string

	if params != ""  {
		httpUrl = fullApiUrl + params
	} else {
		httpUrl = fullApiUrl
	}

	resp, err := client.Get(httpUrl)

	if err != nil {
		return nil, errors.New("HttpsApiGet() https " + httpUrl + " get error with NewRequest")
	}

	defer resp.Body.Close()

	if ( resp.StatusCode  == 200 ) {

		data, _ := ioutil.ReadAll(resp.Body)
		respons := ioutil.NopCloser(bytes.NewReader(data))
		return respons, nil
	} else {
		return nil, errors.New("HttpsApiGet() resp status code not 200.")
	}

}

