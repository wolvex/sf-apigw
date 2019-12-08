package apigw

import (
	"encoding/json"
)

type DISMessage struct {
	Mdn           string `json:"mdn,omitempty"`
	TransactionID int    `yaml:"transactionId,omitempty"`
	BankCode      string `yaml:"bankCode,omitempty"`
	ReferenceNo   string `yaml:"referenceNo,omitempty"`
	QRCode        string `yaml:"qrcode,omitempty"`
	Error         string `yaml:"error,omitempty"`
}

func (ws *ClientService) DISOrderPayment(mdn, bankCode, referenceNo string, transactionId int) (response DISMessage, err error) {

	data := map[string]interface{}{
		"mdn":           NormalizeMDN(mdn),
		"bankCode":      bankCode,
		"referenceNo":   referenceNo,
		"transactionId": transactionId,
	}

	var res []byte
	if res, err = ws.Post("/dis/order/payment", data); err != nil {
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		return
	}

	return
}
