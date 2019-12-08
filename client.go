package apigw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ClientServices interface {
	New(url string, username string, password string, timeout int64) ClientService
	SetProxy(uri string)
}

type ClientService struct {
	BaseURL   string
	Version   string
	SecretKey string
	KeyID     string
	Timeout   int
	Transport http.Transport
}

func New(baseUrl, version, keyId, secretKey string, timeout int) *ClientService {
	api := &ClientService{
		BaseURL:   baseUrl,
		Version:   version,
		SecretKey: secretKey,
		KeyID:     keyId,
		Timeout:   timeout,
	}

	if strings.HasPrefix(baseUrl, "https") {
		api.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return api
}

func (ws *ClientService) SetProxy(uri string) error {
	if uri != "" {
		if proxy, err := url.Parse(uri); err != nil {
			return err
		} else {
			ws.Transport.Proxy = http.ProxyURL(proxy)
		}
	}
	return nil
}

func (ws *ClientService) Post(uri string, data map[string]interface{}) (response []byte, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		return
	}

	path := ws.BaseURL + uri

	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", path, bytes.NewBufferString(string(payload[:])))
	if err != nil {
		return
	}

	loc, _ := time.LoadLocation("MST")
	date := time.Now().In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Version", ws.Version)
	req.Header.Add("Date", date)

	//Adds authorization headers
	if ws.SecretKey != "" && ws.KeyID != "" {
		// Prepare the signature to include those headers:
		//data := "(request-target): post " + uri + "\n"
		data := "date: " + date
		//fmt.Printf("SECRET: [%s] , DATA: [%s]\n", ws.SecretKey, data)

		h := hmac.New(sha256.New, []byte(ws.SecretKey))
		h.Write([]byte(data))

		// Base64 and URL Encode the string
		sigString := base64.StdEncoding.EncodeToString(h.Sum(nil))
		encodedString := url.QueryEscape(sigString)

		//fmt.Printf("BASE64: [%s] , ESCAPE: [%s]\n", sigString, encodedString)

		//req.Header.Add("Authorization", "Signature keyid="+ws.KeyID+",algorithm=hmac-sha256,headers=date,signature="+encodedString)
		req.Header.Add("Authorization", "Signature keyid="+ws.KeyID+",algorithm=hmac-sha256,signature="+encodedString)
	}

	client := &http.Client{
		Timeout:   time.Duration(ws.Timeout) * time.Millisecond,
		Transport: &ws.Transport,
	}

	requestDump, _ := httputil.DumpRequest(req, true)
	fmt.Printf("HTTP Request:\n%q\n", requestDump)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	responseDump, _ := httputil.DumpResponse(res, true)
	fmt.Printf("HTTP Response:\n%q\n", responseDump)

	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	//fmt.Println(string(bRes))
	//return bRes, nil

	//if err = json.Unmarshal(bRes, &response); err != nil {
	//	return
	//}

	return
}
