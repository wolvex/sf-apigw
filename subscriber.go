package apigw

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SubscriberData struct {
	Mdn               string `json:"mdn,omitempty"`
	State             string `json:"state,omitempty"`
	Imsi              string `json:"imsi,omitempty"`
	PUK1              string `json:"puk1,omitempty"`
	PUK2              string `json:"puk2,omitempty"`
	MarketingCategory string `json:"marketingCategory,omitempty"`
	FraudLocked       string `json:"fraudLocked,omitempty"`
	AccountNumber     string `json:"acctNbr,omitempty"`
	ActiveDate        string `json:"activeDate,omitempty"`
	SettlementMethod  string `json:"settlementMethod,omitempty"`
}

func (ws *ClientService) SubscriberQuery(mdn string) (response SubscriberData, err error) {

	data := map[string]interface{}{
		"mdn": NormalizeMDN(mdn),
	}

	var res []byte
	if res, err = ws.Post("/crm/subscriber/query", data); err != nil {
		return
	}

	if err = json.Unmarshal(res, &response); err != nil {
		return
	}

	if response.State == "" || !strings.Contains("ABDEG", response.State) {
		err = fmt.Errorf("Failure query subscriber: %#v", response)
	}

	if activeDate, e := time.Parse("02/01/2006 15:04:05", response.ActiveDate); e == nil {
		response.ActiveDate = activeDate.Format("2006-01-02 15:04:05")
	} else {
		fmt.Printf("%#v", e)
	}

	return
}
