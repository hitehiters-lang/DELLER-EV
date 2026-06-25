package dto

type UploadMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Link    string `json:"link"`
}
