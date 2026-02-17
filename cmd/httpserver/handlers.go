package main

import (
	"bytes"
	"os"

	"github.com/pavychsasha/httpfromtcp/internal/request"
	"github.com/pavychsasha/httpfromtcp/internal/response"
)

func handleRoutes(w *response.Writer, req *request.Request) {
	var headerMessage []byte
	var pMessage []byte
	var statusCode response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.StatusBadRequest
		headerMessage = []byte(response.GetStatusLine(statusCode))
		pMessage = []byte("Your request honestly kinda sucked.")

	case "/myproblem":
		statusCode = response.StatusInternalServerError
		headerMessage = []byte(response.GetStatusLine(statusCode))
		pMessage = []byte("Okay, you know what? This one is on me.")

	default:
		statusCode = response.StatusResponseOK
		headerMessage = []byte("Success!")
		pMessage = []byte("Your request was an absolute banger.")
	}

	body := renderHTMLTemplateBody(statusCode, headerMessage, pMessage)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	for respKey, respVal := range req.Headers {
		headers.Override(respKey, respVal)
	}

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func renderHTMLTemplateBody(statusCode response.StatusCode, headerMessage, pMessage []byte) []byte {

	data, err := os.ReadFile("template.html")
	if err != nil {
		return []byte{}
	}
	title := response.GetWriterStatusLine(statusCode)
	data = bytes.ReplaceAll(data, []byte("{{title}}"), []byte(title))
	data = bytes.ReplaceAll(data, []byte("{{header_message}}"), headerMessage)
	data = bytes.ReplaceAll(data, []byte("{{p_message}}"), pMessage)
	return data
}
