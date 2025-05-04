package rest

const (
	StatusError         = "error"
	StatusSuccess       = "success"
	ErrUndefinedMessage = "undefined"
)

type Message struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error"`
}

func ErrorMessage(err error) Message {
	if err == nil {
		return Message{
			Status: StatusError,
			Error:  ErrUndefinedMessage,
		}
	}

	return Message{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func Success(data any) Message {
	return Message{
		Status: StatusSuccess,
		Data:   data,
	}
}

type Request struct {
	Method string `json:"method"`
	Data   any    `json:"data"`
}
