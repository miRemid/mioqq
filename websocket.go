package mioqq

// WebsocketAPI websocket
type WebsocketAPI struct {
	API   string
	Token string
}

// Connect 链接cqhttp websocket
func (websocket *WebsocketAPI) Connect() error {
	return nil
}

// Call 调用api
func (websocket *WebsocketAPI) Call(endpoint string, params CQParams) (CQAPIResponse, error) {
	return CQAPIResponse{}, nil
}
