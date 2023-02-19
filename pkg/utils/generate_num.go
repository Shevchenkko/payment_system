package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateNumber(itn int) (int64, string) {
	//convert created time to int
	timeT := time.Now()
	tUnixMilli := int64(time.Nanosecond) * timeT.UnixNano() / int64(time.Millisecond)

	cardStr := fmt.Sprintf("%016d", rand.Int63n(1e16))
	cardNumber, err := strconv.Atoi(cardStr)
	if err != nil {
		return 0, ""
	}
	card := int64(cardNumber+itn) + tUnixMilli

	ibanStr := fmt.Sprintf("%026d", rand.Int63n(1e16)+tUnixMilli+int64(itn))
	iban := fmt.Sprintf("UA %s", ibanStr)

	return card, iban
}
