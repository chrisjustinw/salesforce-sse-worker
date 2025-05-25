package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"salesforce-sse-worker/internal/request"
	"salesforce-sse-worker/internal/service"
)

type ConversationHandler interface {
	CreateConversation(e echo.Context) error
	GenerateContinuationToken(e echo.Context) error
}

type ConversationHandlerImpl struct {
	conversationService service.ConversationService
}

func NewConversationHandler(conversationService service.ConversationService) ConversationHandler {
	return &ConversationHandlerImpl{conversationService: conversationService}
}

func (m *ConversationHandlerImpl) CreateConversation(e echo.Context) error {
	var req request.CreateConversationRequest

	if err := e.Bind(&req); err != nil {
		return e.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return e.JSON(400, map[string]string{"error": "Validation failed", "details": err.Error()})
	}

	resp, err := m.conversationService.PublishConversation(e.Request().Context(), req)
	if err != nil {
		return e.JSON(500, map[string]string{"error": err.Error()})
	}

	return e.JSON(200, resp)
}

func (m *ConversationHandlerImpl) GenerateContinuationToken(e echo.Context) error {
	resp, err := m.conversationService.GenerateContinuationToken(e.Request().Context())
	if err != nil {
		return e.JSON(500, map[string]string{"error": err.Error()})
	}

	return e.JSON(200, resp)
}
