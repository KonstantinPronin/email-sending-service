package constants

//Internal
const (
	MaxPageSize        = 1000
	DefaultPageSize    = 20
	MaxConnectAttempts = 5
	MaxSendingWorkers  = 10
)

//http ContentType
const (
	Text = "text/plain"
	Json = "application/json"
)

//http x-headers
const (
	Total      = "X-Total"
	TotalPages = "X-Total-Pages"
	PerPage    = "X-Per-Page"
	Page       = "X-Page"
	NextPage   = "X-Next-Page"
	PrevPage   = "X-Prev-Page"
)
