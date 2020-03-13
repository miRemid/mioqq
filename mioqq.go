package mioqq

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// API is the cqhttp api sdk
type API struct {
	API   string
	Token string

	client *http.Client
	wsPool *sync.Pool
}

var (
	// Timeout set the mio message handler time
	Timeout int = 10
)

// New return a mio client
func New(token, api string) (*API, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		return NewHTTP(token, api)
	case "ws":
		fallthrough
	case "wss":
		return NewWebsocket(token, api)
	default:
		return nil, errors.New("invalid api's scheme")
	}
}

// NewHTTP will create a Mio HTTP Client
func NewHTTP(token, api string) (*API, error) {
	return &API{
		Token: token,
		API:   api,
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: time.Duration(Timeout) * time.Second,
				}).Dial,
				TLSHandshakeTimeout: time.Duration(Timeout) * time.Second,
			},
			Timeout: time.Second * time.Duration(Timeout),
		},
	}, nil
}

// NewWebsocket will create a Mio Websocket Client
func NewWebsocket(token, wsurl string) (*API, error) {
	var mio API
	mio.Token = token
	mio.wsPool = &sync.Pool{
		New: func() interface{} {
			var api WebsocketAPI
			api.API = wsurl
			api.Token = token
			return &api
		},
	}
	return &mio, nil
}

// Request will save the message into the mio's message channel
func (m *API) Request(endpoint string, reader io.Reader) (CQAPIResponse, error) {
	if m.client == nil {
		return m.sendWsRequest(endpoint, reader)
	}
	return m.sendHTTPRequest(endpoint, reader)
}

func (m *API) sendHTTPRequest(endpoint string, reader io.Reader) (CQAPIResponse, error) {
	var uri = m.API + "/" + endpoint
	var response CQAPIResponse
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return response, err
	}
	if m.Token != "" {
		req.Header.Add("Authorization", "Bearer "+m.Token)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(data, &response)
	return response, err
}

func (m *API) sendWsRequest(endpoint string, reader io.Reader) (CQAPIResponse, error) {
	var params CQParams
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return CQAPIResponse{}, err
	}
	if err := json.Unmarshal(data, &params); err != nil {
		return CQAPIResponse{}, err
	}
	api := m.wsPool.Get().(*WebsocketAPI)
	defer m.wsPool.Put(api)
	if err := api.Connect(); err != nil {
		return CQAPIResponse{}, err
	}
	return api.Call(endpoint, params)
}
