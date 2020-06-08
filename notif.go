package apigw

func (ws *ClientService) SendSms(mdn, text string) (err error) {

	data := map[string]interface{}{
		"from": "62888",
		"to":   NormalizeMDN(mdn)[2:],
		"text": text,
	}

	if _, err = ws.Post("/crm/notification/sms", data); err != nil {
		return
	}

	return
}
