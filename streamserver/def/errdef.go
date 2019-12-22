package def

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyPaseFailed = ErrorResponse{HttpSC: 400, Error: Err{Error: "request body is not correct", ErrorCode: "001"}}
	ErrorNotAuthUser           = ErrorResponse{HttpSC: 401, Error: Err{Error: "User anthentication failed.", ErrorCode: "002"}}
	ErrorDBError               = ErrorResponse{HttpSC: 500, Error: Err{Error: "DB ops failed", ErrorCode: "003"}}
	ErrorInternalFaults        = ErrorResponse{HttpSC: 500, Error: Err{Error: "Internal service error", ErrorCode: "004"}}
	ErrorConnectLimit          = ErrorResponse{HttpSC: 429, Error: Err{Error: "Too many request!.", ErrorCode: "005"}}
)
