package common

type UserInfo struct {
	PartnerId uint64 `json:"partner_id"`
	UserId    uint64 `json:"user_id"`
	ScopeId   uint64 `json:"scope_id"`
}
