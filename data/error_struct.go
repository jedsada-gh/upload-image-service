package data

// Error is exported because it starts with a capital letter
type Error struct {
	ErrorDetail ErrorDetail `json:"error"`
}

// ErrorDetail is exported because it starts with a capital letter
type ErrorDetail struct {
	Message string `json:"message"`
}
