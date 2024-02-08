package web

type StandartResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewStandartResponse(status int, message string, data interface{}) *StandartResponse {
	return &StandartResponse{Status: status, Message: message, Data: data}
}
