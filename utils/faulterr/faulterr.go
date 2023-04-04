package faulterr

// FaultErr structure
type FaultErr struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
