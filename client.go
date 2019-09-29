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
	SetEnv(env string)
	SetProxy(uri string)
	NumberShift(mdn string, newMdn string) (map[string]interface{}, error)
	GetSubInfo(mdn string) (map[string]interface{}, error)
}

type ClientService struct {
	BaseURL   string
	SecretKey string
	KeyID     string
	Timeout   int
	Transport http.Transport
}

func New(baseUrl, secretKey, keyId string, timeout int) *ClientService {
	api := &ClientService{
		BaseURL:   baseUrl,
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

func (ws *ClientService) Post(uri string, data map[string]interface{}) (response map[string]interface{}, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		return
	}

	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(string(payload[:])))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Date", time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))

	//Adds authorization headers
	if ws.SecretKey != "" && ws.KeyID != "" {
		// Prepare the signature to include those headers:
		h := hmac.New(sha256.New, []byte(ws.SecretKey))
		h.Write([]byte("date: " + req.Header.Get("Date")))

		// Base64 and URL Encode the string
		sigString := base64.StdEncoding.EncodeToString(h.Sum(nil))
		encodedString := url.QueryEscape(sigString)
		req.Header.Add("Authorization", "Signature keyId=\""+ws.KeyID+"\",algorithm=\"hmac-sha256\",signature=\""+encodedString+"\"")
	}

	client := &http.Client{
		Timeout:   time.Duration(ws.Timeout) * time.Millisecond,
		Transport: &ws.Transport,
	}

	requestDump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(requestDump))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	bRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bRes, &response); err != nil {
		return
	}

	return
}

func (ws *ClientService) ServiceBuy(mdn, serviceCode string) (response map[string]interface{}, err error) {

	data := map[string]interface{}{
		"mdn":         NormalizeMDN(mdn),
		"serviceCode": serviceCode,
	}

	response, err = ws.Post("/crm/service/buy", data)

	return
}
