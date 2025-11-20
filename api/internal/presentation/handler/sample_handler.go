package handler

import (
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"net/http"
)

type SampleHandler struct {
}

func NewSampleHandler() *SampleHandler {
	return &SampleHandler{}
}

type GreetingHandlerResponse struct {
	Message string `json:"message"`
}

func (h *SampleHandler) GreetingHandler(w http.ResponseWriter, r *http.Request) {
	helper.RespondWithJSON(w, http.StatusOK, GreetingHandlerResponse{Message: "Hello, World!"})
}
