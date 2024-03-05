package short

import "time"

// Short represents a record in the short table
type Short struct {
	ID  string `gorm:"primaryKey;default:gen_uid()" json:"id"`
	URL string `json:"url"`
}

// AccessLog represents a record in the access log table
type AccessLog struct {
	ID         int       `gorm:"primarykey" json:"id"`
	ShortID    string    `json:"short_id"`
	AccessTime time.Time `json:"access_time"`
}

// Audit represents a record in the audits table
type Audit struct {
	ID      int    `gorm:"primarykey" json:"id"`
	Action  string `json:"action"`
	Path    string `json:"path"`
	Request string `json:"request"`
	Invoker string `json:"invoker"`
	Latency string `json:"latency"`
}

// ShortRequest is a request sent
type ShortRequest struct {
	URL string `json:"url" binding:"url"`
}

// ShortResponse is the response from the details endpoint. It contains the url as well as a few metrics.
type ShortResponse struct {
	Short
	Last24Hours int64 `json:"last_24_hours"`
	PastWeek    int64 `json:"past_week"`
	AllTime     int64 `json:"all_time"`
}
