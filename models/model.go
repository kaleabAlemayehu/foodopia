package models

type Payload struct {
	Id       int64  `json:"id"`
	Username string `json:"username"  binding:"required"`
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	Error    string `json:"error"`
}

type Params struct {
	Payload Payload `json:"params"`
}
type UserActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Params                 `json:"input"`
}

type ImageUploadArgs struct {
	Name      string `json:"name"`
	Base64Str string `json:"base64Str"`
}
type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            ImageUploadArgs        `json:"input"`
}

type SaveImageOutput struct {
	ImageUrl string `json:"image_url"`
}
