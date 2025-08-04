package rest

type ErrorResp struct {
	Status  string
	Message string
}

func errorResponse(message string) ErrorResp {
	return ErrorResp{
		Status:  "error",
		Message: message,
	}
}
