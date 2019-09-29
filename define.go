package apigw

import "strings"

func NormalizeMDN(mdn string) string {
	if strings.HasPrefix(mdn, "6288") {
		return mdn
	}

	if strings.HasPrefix(mdn, "+62") {
		//remove leading +62
		mdn = strings.Replace(mdn, "+62", "62", 1)
	} else if strings.HasPrefix(mdn, "088") {
		//remove leading 0
		mdn = strings.Replace(mdn, "088", "6288", 1)
	} else if strings.HasPrefix(mdn, "88") {
		//remove leading 0
		mdn = strings.Replace(mdn, "88", "6288", 1)
	}
	return mdn
}
