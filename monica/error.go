package monica

type errorResponse struct {
	Error struct {
		Message   string `json:"message"`
		ErrorCode int    `json:"error_code"`
	} `json:"error"`
}
