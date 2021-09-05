package dtos

type AccessTokenResponse struct {
	Status string            `json:"status"`
	Error  bool              `json:"error"`
	Result AccessTokenResult `json:"result"`
}

type AccessTokenResult struct {
	Token string `json:"token"`
}
