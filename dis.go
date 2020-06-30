package apigw

import (
	"encoding/json"
	"fmt"
)

type AddOns struct {
	ServiceCode string  `json:"serviceCode,omitempty"`
	Price       float64 `yaml:"price,omitempty"`
}

type DISMessage struct {
	Mdn            string   `json:"mdn,omitempty"`
	TransactionID  int      `yaml:"transactionId,omitempty"`
	BankCode       string   `yaml:"bankCode,omitempty"`
	ReferenceNo    string   `yaml:"referenceNo,omitempty"`
	QRCode         string   `yaml:"qrcode,omitempty"`
	ActivationCode string   `yaml:"activationCode,omitempty"`
	Error          string   `yaml:"error,omitempty"`
	Status         string   `yaml:"status,omitempty"`
	Code           int      `yaml:"Code,omitempty"`
	Price          float64  `yaml:"price,omitempty"`
	AddOn          []AddOns `yaml:"addons,omitempty"`
}

func (ws *ClientService) DISOrderPayment(mdn, bankCode, referenceNo string, transactionId int) (response DISMessage, err error) {

	data := map[string]interface{}{
		"mdn":           NormalizeMDN(mdn),
		"bankCode":      bankCode,
		"referenceNo":   referenceNo,
		"transactionId": transactionId,
	}

	var res []byte
	//if res, err = ws.Post("/onboarding/v1/order/payment", data); err != nil {
	if res, err = ws.Post("/dis/order/payment", data); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	return
}

func (ws *ClientService) DISOrderPaymentV1(mdn, bankCode, referenceNo string, transactionId int) (response DISMessage, err error) {

	data := map[string]interface{}{
		"mdn":           NormalizeMDN(mdn),
		"bankCode":      bankCode,
		"referenceNo":   referenceNo,
		"transactionId": transactionId,
	}

	var res []byte
	if res, err = ws.Post("/onboarding/v1/order/payment", data); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	return
}

func (ws *ClientService) DISOrderInquiryV1(transactionId int) (response DISMessage, err error) {

	data := map[string]interface{}{
		"transactionId": transactionId,
	}

	var res []byte
	if res, err = ws.Post("/onboarding/v1/order/query", data); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		fmt.Printf("Exception caught %#v", err)
		return
	}

	return
}
