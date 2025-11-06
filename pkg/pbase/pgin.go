package pbase

import "net/http"

const (
	PXHTTPHeaderUser  = "PX-User"
	PXHTTPCCRequestID = "x-request-id"
)

const (
	ContextRequestID = "request_id"
)

func GetUser(header http.Header) string {
	return header.Get(PXHTTPHeaderUser)
}

func GetHTTPRequestID(header http.Header) string {
	return header.Get(PXHTTPCCRequestID)
}
