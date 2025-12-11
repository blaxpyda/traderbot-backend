package auth

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"thugcorp.io/final_bot/logger"
	"thugcorp.io/final_bot/utils"
)

type authHandler struct {
	logger *logger.Logger
	repo AuthRepository
	service AuthService
}

func NewAuthHandler(  logger *logger.Logger, repo AuthRepository, service AuthService) *authHandler {
	return &authHandler{
		logger: logger,
		repo: repo,
		service: service,
	}
}

func (h *authHandler) HandleCallBack(w http.ResponseWriter, r *http.Request) {
	// Get code from the redirect URL
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("No code in callback URL", nil)
		http.Error(w, "No code in callback URL", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Exchange code for token
	token, err := h.service.ExchangeCodeForToken(code)
	if err != nil {
		h.logger.Error("Failed to exchange code for token", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// Get real user Info from Alpaca
	alpacaID, err := h.service.FetchAlpacaUserInfo(ctx, token.AccessToken)
	if err != nil {
		h.logger.Error("Failed to fetch user info from Alpaca", err)
		http.Error(w, "Failed to fetch user info from Alpaca", http.StatusInternalServerError)
		return
	}

	// Upsert user in DB
	err = h.repo.UpsertUser(alpacaID, token.AccessToken, token.RefreshToken, token.Expiry.Unix())
	if err != nil {
		h.logger.Error("Failed to upsert user in DB", err)
		http.Error(w, "Failed to upsert user in DB", http.StatusInternalServerError)
		return
	}

	// Generate JWT token for mobile app
	jwtToken, err := utils.GenerateJWT(alpacaID, time.Hour*24)
	if err != nil {
		h.logger.Error("Failed to generate JWT token", err)
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Redirect back to mobile app with JWT token
	redirectURL := "myapp://auth?token=" + jwtToken
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *authHandler) RegisterRoutes(r mux.Router) {
	r.HandleFunc("/auth/callback", h.HandleCallBack).Methods("GET")
}