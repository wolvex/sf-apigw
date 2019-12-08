package apigw

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (ws *ClientService) CreateTicket(mdn string) (response SubscriberData, err error) {

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
