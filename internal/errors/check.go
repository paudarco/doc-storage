package errors

import (
	"net/http"
)

func CheckError(err error) int {
	for errCode, list := range errorsList {
		if _, ok := list[err]; ok {
			return errCode
		}
	}

	return http.StatusInternalServerError
}
