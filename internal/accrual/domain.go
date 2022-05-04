package accrual

const (
	NEW        = "NEW"
	REGISTERED = "REGISTERED"
	INVALID    = "INVALID"
	PROCESSING = "PROCESSING"
	PROCESSED  = "PROCESSED"
)

type AccrualInfo struct {
	OrderID string
	Amount  float64 `json:"accrual"`
	Status  string  `json:"status"`
}

type Response struct {
	OrderID string
	Amount  float64
	Status  string
	IsFinal bool
	Timeout int64
}
