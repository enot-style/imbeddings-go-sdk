package imbeddings

// RequestError describes failures in building or executing a request.
type RequestOp string

const (
	RequestOpMarshalRequest RequestOp = "marshal request"
	RequestOpCreateRequest  RequestOp = "create request"
	RequestOpCallService    RequestOp = "call service"
	RequestOpReadImageData  RequestOp = "read image data"
	RequestOpEncodeImage    RequestOp = "encode image data"
)

type RequestError struct {
	Op  RequestOp
	Err error
}

// Error formats RequestError as "op: err" (or "op" when Err is nil).
func (e RequestError) Error() string {
	if e.Err == nil {
		return string(e.Op)
	}
	return string(e.Op) + ": " + e.Err.Error()
}

// ValidationError describes invalid inputs before a request is made.
type ValidationOp string

const (
	ValidationOpInit      ValidationOp = "init"
	ValidationOpSetModel  ValidationOp = "set model"
	ValidationOpParams    ValidationOp = "params"
	ValidationOpImage ValidationOp = "image"
)

type ValidationError struct {
	Op  ValidationOp
	Err error
}

// Error formats ValidationError as "op: err" (or "op" when Err is nil).
func (e ValidationError) Error() string {
	if e.Err == nil {
		return string(e.Op)
	}
	return string(e.Op) + ": " + e.Err.Error()
}

// ResponseError describes failures while decoding or interpreting a response.
type ResponseOp string

const (
	ResponseOpDecodeResponse ResponseOp = "decode response"
	ResponseOpEmptyResponse  ResponseOp = "empty response"
	ResponseOpResult         ResponseOp = "result"
)

type ResponseError struct {
	Op  ResponseOp
	Err error
}

// Error formats ResponseError as "op: err" (or "op" when Err is nil).
func (e ResponseError) Error() string {
	if e.Err == nil {
		return string(e.Op)
	}
	return string(e.Op) + ": " + e.Err.Error()
}
