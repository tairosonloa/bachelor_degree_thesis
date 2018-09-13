package models

// MessageAPI should be used for common JSON responses
type MessageAPI struct {
	Message string `json:"message"`
}

// ErrorAPI should be used for resposes containing a JSON with an error message
type ErrorAPI struct {
	Error string `json:"error"`
}

// CPDStatusAPI should be used to serve current CPD status
type CPDStatusAPI struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	UPSStatus   string `json:"ups_status (LDI rack)"`
}
