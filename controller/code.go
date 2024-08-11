package controller

type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

var resp Response

func (*Response) InternalError(err error) Response {
	return Response{
		Code:    500,
		Message: err.Error(),
	}
}

func (*Response) InValidRequest(err error) Response {
	return Response{
		Code:    400,
		Message: err.Error(),
	}
}
func (*Response) Success() Response {
	return Response{
		Code:    200,
		Message: "OK",
	}
}
