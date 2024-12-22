package response

import "github.com/ryvasa/go-super-farmer/service_api/model/dto"

type ResponseMessage struct {
	Status  int     `json:"status"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    Message `json:"data"`
	Errors  Error   `json:"errors"`
}

type Message struct {
	Message string `json:"message"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponseDownload struct {
	Status  int                     `json:"status"`
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    dto.DownloadResponseDTO `json:"data"`
	Errors  Error                   `json:"errors"`
}
