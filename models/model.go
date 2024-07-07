package models

type Payload struct {
	Id       int64  `json:"id"`
	Username string `json:"username"  binding:"required"`
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type Params struct {
	Payload Payload `json:"params"`
}
type UserActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Params                 `json:"input"`
}
