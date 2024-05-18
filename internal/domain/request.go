package domain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server interface {
	SendRequest(ctx context.Context, req *http.Request) (*http.Response, error)
	IsHealthy(ctx context.Context) bool
}

type HttpServerSender struct {
	host string
	port string
	// You can include any dependencies here.
}

func NewHttpServerSender(host string, port string) *HttpServerSender {
	// Initialize any dependencies here.
	return &HttpServerSender{
		host: host,
		port: port,
	}
}

func (s *HttpServerSender) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	request, err := http.NewRequest(req.Method, fmt.Sprintf("http://%s:%s/%s", s.host, s.port, req.URL.Path), req.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("HttpServerSender: SendRequest", s.host, s.port)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *HttpServerSender) IsHealthy(ctx context.Context) bool {
	method := "GET"

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%s/", s.host, s.port), nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 (OK)
	if resp.StatusCode == http.StatusOK {
		return true
	}

	// You can add more status code checks if needed
	return false
}

type httpServerSenderLogger struct {
	base Server
}

func NewHttpServerSenderLogger(base Server) Server {
	return httpServerSenderLogger{
		base: base,
	}
}

func (h httpServerSenderLogger) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := h.base.SendRequest(ctx, req)
	if err != nil {
		log.Printf("HttpServerSender: SendRequest error: %s\n", err)
	} else {
		log.Printf("HttpServerSender: SendRequest: %s\n", resp.Status)
	}
	return resp, err
}

func (h httpServerSenderLogger) IsHealthy(ctx context.Context) bool {
	isHealthy := h.base.IsHealthy(ctx)
	log.Printf("HttpServerSender: IsHealthy Method Called: %t\n", isHealthy)
	return isHealthy
}

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type Header map[string][]string

type Request struct {
	// You can include any request parameters here.
	Headers Header
	Body    string
	Method  Method
	URL     string
	Host    string
}

type Status int

const (
	Success             Status = 200
	NotFound            Status = 404
	InternalServerError Status = 500
)

type Response struct {
	// You can include any response parameters here.
	Headers map[string]string
	Body    string
	Status  Status
}

type Service struct {
	Host string
	Port string
}

func SwToRequest(request *http.Request) Request {
	return Request{
		URL:     request.URL.String(),
		Method:  Method(request.Method),
		Body:    SwToRequestBody(request.Body),
		Headers: SwRequestHeaderTo(request.Header),
		Host:    request.Host,
	}
}

func RequestToHttpRequest(request Request) *http.Request {
	req, _ := http.NewRequest(string(request.Method), request.URL, bytes.NewBuffer([]byte(request.Body)))
	req.Host = request.Host
	req.Header = http.Header(request.Headers)
	return req
}

func SwToRequestBody(body io.ReadCloser) string {
	if body == http.NoBody {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func SwRequestHeaderTo(header http.Header) Header {
	headers := make(Header)
	for key, value := range header {
		headers[key] = value
	}
	return headers
}
