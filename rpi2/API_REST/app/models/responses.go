package models

// MessageAPIM should be used for common JSON responses
type MessageAPIM struct {
	Message string `json:"message"`
}

// ErrorAPIM should be used for resposes containing a JSON with an error message
type ErrorAPIM struct {
	Error string `json:"error"`
}

// CPDStatusAPIM should be used to serve current CPD status
type CPDStatusAPIM struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	UPSStatus   string  `json:"ups_status (LDI rack)"`
}
