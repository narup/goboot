package util

// Error represents the API level error for the client apps
type Error struct {
	ID      string `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// Errors represents json errors
type Errors struct {
	Errors []*Error `json:"errors"`
}

//DataNotFoundError returns true if the given error is data not found
func DataNotFoundError(err error) bool {
	return err.Error() == "not found"
}

// APIError construct Error
func APIError(ID string, code int, message string, err error) *Error {
	return &Error{ID, code, err.Error(), err}
}
