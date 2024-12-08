package response

type ResponseMessage struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Message     `json:"data"`
	Errors  interface{} `json:"errors"`
}

type Message struct {
	Message string `json:"message"`
}
