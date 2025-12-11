package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"thugcorp.io/final_bot/logger"
)

type AuthService interface {
	ExchangeCodeForToken(code string) (*oauth2.Token, error)
	FetchAlpacaUserInfo(ctx context.Context, token string) ( string, error)
}

type authService struct {
	authConfig *oauth2.Config
	logger    *logger.Logger
}

func NewAuthService(authConfig *oauth2.Config, logger *logger.Logger) AuthService {
	return &authService{
		authConfig: authConfig,
		logger: logger,
	}
}

func (s *authService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	// Exchange code for token
	ctx := context.Background()
	token, err := s.authConfig.Exchange(ctx, code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return nil, err
	}
	return token, nil
}

func (s *authService) FetchAlpacaUserInfo(ctx context.Context, token string) (string, error) {
	client := s.authConfig.Client(ctx, &oauth2.Token{AccessToken: token})
	resp, err := client.Get("https://api.alpaca.markets/v2/account")
	if err != nil {
		s.logger.Error("Failed to fetch user info from Alpaca", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Non-OK HTTP status from Alpaca user info endpoint", nil)
		return "", fmt.Errorf("failed to fetch user info from Alpaca, status code: %d %s", resp.StatusCode, resp.Status)
	}

	var result struct {
		ID    string `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		s.logger.Error("Failed to decode Alpaca user info response", err)
		return "", err
	}

	return result.ID, nil
}