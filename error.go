package freshdesk

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	Err         error
	ReqBody     string
	ResBody     string
	Description string       `json:"description"`
	Errors      []FieldError `json:"errors"`
	StatusCode  int
}

type FieldError struct {
	Field          string                 `json:"field"`
	AdditionalInfo map[string]interface{} `json:"additional_info"`
	Message        string                 `json:"message"`
	Code           string                 `json:"code"`
}

func NewApiError(statusCode, expectedStatus int, req, res string) *APIError {
	i := &APIError{
		Err:        fmt.Errorf("received status code %d (%d expected)", statusCode, expectedStatus),
		ReqBody:    req,
		ResBody:    res,
		StatusCode: statusCode,
	}
	_ = json.Unmarshal([]byte(res), i)
	return i
}

func (e *APIError) Error() string {
	return e.Err.Error() + " - " + e.ResBody
}

func (e *APIError) IsDuplicate() (bool, uint64) {
	for _, err := range e.Errors {
		if err.Field == "email" && err.Code == "duplicate_value" {
			userID, ok := err.AdditionalInfo["user_id"]
			if !ok || userID == nil {
				return false, 0
			}
			return true, uint64(userID.(float64))
		}
	}
	return false, 0
}
