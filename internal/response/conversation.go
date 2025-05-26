package response

type GenerateTokenResponse struct {
	AccessToken string `json:"accessToken"`
	LastEventId string `json:"lastEventId"`
}
