package apigw

import "fmt"

type ServiceRes struct {
	ServiceName   string `json:"serviceCode,omitempty"`
	EffectiveDate string `json:"effDate,omitempty"`
	ExpiryDate    string `json:"expDate,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
}

func (ws *ClientService) ServiceBuy(mdn, serviceCode string) (response ServiceRes, err error) {

	data := map[string]interface{}{
		"mdn":         NormalizeMDN(mdn),
		"serviceCode": serviceCode,
	}

	var res interface{}
	if res, err = ws.Post("/crm/service/buy", data); err != nil {
		return
	}

	var ok bool
	if response, ok = res.(ServiceRes); !ok {
		err = fmt.Errorf("Invalid response format")
		return

	}

	return
}
