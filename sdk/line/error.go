package line

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrInvalidRequest          = "INVALID_REQUEST"
	ErrAccessDenied            = "ACCESS_DENIED"
	ErrUnsupportedResponseType = "UNSUPPORTED_RESPONSE_TYPE"
	ErrInvalidScope            = "INVALID_SCOPE"
	ErrServerError             = "SERVER_ERROR"
)

type Error struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("type: %s, description: %s", e.Type, e.Description)
}

type RespError struct {
	Type        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *RespError) Error() string {
	return fmt.Sprintf("type: %s, description: %s", e.Type, e.Description)
}

func ErrHandler(fullUrl string) (*Error, error) {
	parse, err := url.Parse(fullUrl)
	if err != nil {
		return nil, errors.New("url parse error")
	}

	errorCode := parse.Query().Get("error")
	errorDescription := parse.Query().Get("error_description")

	if errorCode == "" {
		return nil, errors.New("not found error")
	}

	if errorDescription != "" {
		return &Error{errorCode, errorDescription}, nil
	}

	var Description string

	switch errorCode {
	case ErrInvalidRequest:
		Description = "Problem with the request. Check the query parameters of the authorization URL."
	case ErrAccessDenied:
		Description = "The user canceled on the consent screen and declined to grant permissions to your app."
	case ErrUnsupportedResponseType:
		Description = "Problem with the value of the \"response_type\" query parameter. The LINE Login only supports \"code\"."
	case ErrInvalidScope:
		Description = "Problem with the value of the scope query parameter. Make sure you've specified an appropriate value.\n\nprofile or openid is required.\nIf you specify email, you also have to specify openid."
	case ErrServerError:
		Description = "An unexpected error occurred on the LINE Login server."
	}

	return &Error{
		Type:        errorCode,
		Description: Description,
	}, nil
}
