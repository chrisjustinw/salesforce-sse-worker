package response

type GenerateContinuationTokenResponse struct {
	AccessToken string `json:"accessToken"`
	LastEventId string `json:"lastEventId"`
}
