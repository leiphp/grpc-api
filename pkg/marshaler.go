package pkg

import (
	"encoding/json"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type HttpResponse struct {
	StatusCode int             `json:"status_code"`
	ErrorCode  string          `json:"error_code,omitempty"`
	Message    string          `json:"message,omitempty"`
	Source     string          `json:"source,omitempty"`
	Detail     interface{}     `json:"detail,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

type HexMarshaler struct {
	marshaler runtime.Marshaler
}

func NewHexMarshaler() runtime.Marshaler {
	mm := &HexMarshaler{
		marshaler: &runtime.JSONPb{OrigName: true},
	}
	return mm
}

func (m *HexMarshaler) Marshal(v interface{}) ([]byte, error) {
	data, _ := m.marshaler.Marshal(v)
	r := &HttpResponse{
		StatusCode: 0,
		Payload:    data,
	}
	return json.MarshalIndent(r, "", "")
}

func (m *HexMarshaler) Unmarshal(data []byte, v interface{}) error {
	return m.marshaler.Unmarshal(data, v)
}

func (m *HexMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return m.marshaler.NewDecoder(r)
}

func (m *HexMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return m.marshaler.NewEncoder(w)
}

func (m *HexMarshaler) ContentType() string {
	return m.marshaler.ContentType()
}
