package request

import (
	"net/http"
)

type Header struct {
	Header *http.Header
}

func NewHeader(header map[string]string) (h *Header) {
	for k, v := range header {
		h.Header.Add(k, v)
	}
	return
}
