package apigw

import (
	"encoding/json"
	"fmt"
)

type ServiceData struct {
	Mdn           string `json:"mdn,omitempty"`
	ServiceCode   string `json:"serviceCode,omitempty"`
	EffectiveDate string `json:"effDate,omitempty"`
	ExpiryDate    string `json:"expDate,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
	ReturnCode    string `json:"returnCode,omitempty"`
	ResultMsg     string `json:"resultMsg,omitempty"`
}

func (ws *ClientService) AddService(mdn, serviceCode string) (response ServiceData, err error) {

	data := map[string]interface{}{
		"mdn":         NormalizeMDN(mdn),
		"serviceCode": serviceCode,
	}

	var res []byte
	if res, err = ws.Post("/crm/service/buy", data); err != nil {
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		return
	}

	if response.TransactionID == "" || response.TransactionID == "map[-nil:true]" {
		err = fmt.Errorf("Failure adding service: %#v", response)
	}

	return
}
