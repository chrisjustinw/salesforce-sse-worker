package request

type CreateConversationRequest struct {
	ConversationId    string            `json:"conversationId" validate:"required"`
	RoutingAttributes RoutingAttributes `json:"routingAttributes" validate:"required"`
	EsDeveloperName   string            `json:"esDeveloperName" validate:"required"`
	Language          string            `json:"language" validate:"required"`
}

type RoutingAttributes struct {
	CaseId        string `json:"caseId" validate:"required"`
	AccountId     string `json:"accountId" validate:"required"`
	CustomerName  string `json:"customerName" validate:"required"`
	CustomerPhone string `json:"customerPhone" validate:"required"`
	CustomerEmail string `json:"customerEmail" validate:"required"`
	Origin        string `json:"origin" validate:"required"`
	SourceType    string `json:"sourceType" validate:"required"`
}
