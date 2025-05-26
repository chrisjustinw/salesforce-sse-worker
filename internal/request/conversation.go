package request

type CreateConversationRequest struct {
	ConversationId    string                              `json:"conversationId" validate:"required"`
	EsDeveloperName   string                              `json:"esDeveloperName" validate:"required"`
	Language          string                              `json:"language" validate:"required"`
	RoutingAttributes CreateConversationRoutingAttributes `json:"routingAttributes" validate:"required"`
}

type CreateConversationRoutingAttributes struct {
	CaseId        string `json:"caseId" validate:"required"`
	AccountId     string `json:"accountId" validate:"required"`
	CustomerName  string `json:"customerName" validate:"required"`
	CustomerPhone string `json:"customerPhone" validate:"required"`
	CustomerEmail string `json:"customerEmail" validate:"required"`
	Origin        string `json:"origin" validate:"required"`
	SourceType    string `json:"sourceType" validate:"required"`
}

type GenerateTokenRequest struct {
	OrgId               string               `json:"orgId" validate:"required"`
	EsDeveloperName     string               `json:"esDeveloperName" validate:"required"`
	CapabilitiesVersion string               `json:"capabilitiesVersion" validate:"required"`
	Platform            string               `json:"platform" validate:"required"`
	Context             GenerateTokenContext `json:"context" validate:"required"`
}

type GenerateTokenContext struct {
	AppName       string `json:"appName" validate:"required"`
	ClientVersion string `json:"clientVersion" validate:"required"`
}
