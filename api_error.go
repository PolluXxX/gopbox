package gopbox

import (
	"fmt"
)

type ApiError struct {
	HTTPCode int
	Message  string `json:"error"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("HTTP %d : %s", e.HTTPCode, e.Message)
}
