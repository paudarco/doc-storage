package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRequestBody = errors.New("invalid request body")

	ErrInvalidAdminToken  = errors.New("invalid admin token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenRequired      = errors.New("token required")
	ErrTokenExpired       = errors.New("token expired")
	ErrWrongToken         = errors.New("wrong token")
	ErrInvalidAuthHeader  = errors.New("invalid auth header forman")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")

	ErrWrongPswrdLength   = errors.New("login and password must be at least 8 characters long")
	ErrWrongLoginLength   = errors.New("login and password must be at least 8 characters long")
	ErrLoginWithoutLatin  = errors.New("login must contain only latin letters and digits")
	ErrPswrdWithoutLatin  = errors.New("password must contain only latin letters, digits and punct symbols")
	ErrPswrdWithoutUpper  = errors.New("password must contain at least upper 2 letters")
	ErrPswrdWithoutLower  = errors.New("password must contain at least lower 2 letters")
	ErrPswrdWithoutDigit  = errors.New("password must contain at least 1 digit")
	ErrPswrdWithoutSymbol = errors.New("password must contain at least 1 special character")

	ErrDocNotFound      = errors.New("document not found")
	ErrDocListNotFound  = errors.New("document list not found")
	ErrMetaNameRequired = errors.New("meta.name is required")

	ErrAccessDenied = errors.New("access denied")
	ErrUnauthorized = errors.New("unautharized")
)

var badReqErrList map[error]interface{} = map[error]interface{}{
	ErrInvalidRequestBody: nil,
	ErrLoginWithoutLatin:  nil,
	ErrPswrdWithoutLatin:  nil,
	ErrWrongPswrdLength:   nil,
	ErrWrongLoginLength:   nil,
	ErrPswrdWithoutLower:  nil,
	ErrPswrdWithoutUpper:  nil,
	ErrPswrdWithoutDigit:  nil,
	ErrPswrdWithoutSymbol: nil,
	ErrMetaNameRequired:   nil,
}

var notFoundErrList map[error]interface{} = map[error]interface{}{
	ErrDocNotFound:     nil,
	ErrDocListNotFound: nil,
	ErrUserNotFound:    nil,
}

var unauthErrList map[error]interface{} = map[error]interface{}{
	ErrTokenRequired:      nil,
	ErrInvalidAuthHeader:  nil,
	ErrInvalidCredentials: nil,
	ErrWrongToken:         nil,
	ErrInvalidToken:       nil,
	ErrTokenExpired:       nil,
}

var forbiddenErrList map[error]interface{} = map[error]interface{}{
	ErrInvalidAdminToken: nil,
	ErrAccessDenied:      nil,
}

var conflictErrList map[error]interface{} = map[error]interface{}{
	ErrUserAlreadyExist: nil,
}

var errorsList map[int]map[error]interface{} = map[int]map[error]interface{}{
	http.StatusBadRequest:   badReqErrList,
	http.StatusNotFound:     notFoundErrList,
	http.StatusUnauthorized: unauthErrList,
	http.StatusForbidden:    forbiddenErrList,
	http.StatusConflict:     conflictErrList,
}
