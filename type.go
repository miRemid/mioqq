package mioqq

import "encoding/json"

// CQAPIResponse is the response from cqhttp's api
type CQAPIResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	RetCode int             `json:"retcode"`
}

// CQParams 参数
type CQParams map[string]interface{}

// CQEvent is the interface of cqhttp's event
type CQEvent interface{}
