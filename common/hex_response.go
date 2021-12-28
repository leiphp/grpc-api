package common


type HexResponse struct {
	StatusCode  uint64      `json:"status_code"`
	Code        string      `json:"code,omitempty"`
	Source      string      `json:"source,omitempty"`
	Detail      interface{} `json:"detail"`
	Payload     interface{} `json:"payload"`
	Description string      `json:"description"`
}
