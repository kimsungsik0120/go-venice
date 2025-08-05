package dtos

type EVMError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EvmRPCRequest struct {
	JsonRpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Id      int      `json:"id"`
	Params  []string `json:"params"`
}

type EvmRPCResponse struct {
	JsonRpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Result  string    `json:"result,omitempty"`
	Error   *EVMError `json:"error,omitempty"`
}
