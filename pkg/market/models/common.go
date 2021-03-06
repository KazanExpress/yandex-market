package models

import (
	"fmt"
	"strings"
)

// CommonResponse response structure common for most of responses.
type CommonResponse struct {
	Errors CommonErrors `json:"errors"`
	Status Status       `json:"status"`
}

// CommonErrors list of CommonError.
type CommonErrors []CommonError

func (e CommonErrors) Error() string {
	var b strings.Builder
	for i, e := range e {
		fmt.Fprintf(&b, "err[%v]: %s;", i, e.Error())
	}

	return b.String()
}

// CommonError error structure common for most of responses.
type CommonError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implement error interface.
func (e CommonError) Error() string {
	return fmt.Sprintf("msg: %s, code: %s;", e.Message, e.Code)
}

// CommonPagingOptions common paging options.
type CommonPagingOptions struct {
	PageNumber int32
	PageSize   int32
	Limit      int32
	Offset     int32
}

// Currency is enum for yandex market currencies.
type Currency string

const (
	// CurrencyRUR is russian ruble.
	CurrencyRUR Currency = "RUR"
	// CurrencyBYN is belarussian ruble.
	CurrencyBYN Currency = "BYN"
	// CurrencyKZT is kazakh tenge.
	CurrencyKZT Currency = "KZT"
	// CurrencyUAH is ukrainian grivna.
	CurrencyUAH Currency = "UAH"
)
