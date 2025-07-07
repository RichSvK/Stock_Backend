package output

type WebResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type FailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
