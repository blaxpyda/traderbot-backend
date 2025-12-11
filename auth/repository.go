package auth

import (
	"time"

	"gorm.io/gorm"
	"thugcorp.io/final_bot/logger"
)

type AuthRepository interface {
	UpsertUser(alpacaUserID, accessToken, refreshToken string, expiry int64) error
	GetUserByAlpacaID(alpacaUserID string) (*User, error)
}

type authRepository struct {
	db *gorm.DB
	logger *logger.Logger
}

func NewAuthRepository(db *gorm.DB, logger *logger.Logger) AuthRepository {
	return &authRepository{
		db: db,
		logger: logger,
	}
}

func (r *authRepository) UpsertUser(alpacaUserID,  accessToken, refreshToken string, expiry int64) error {
	user := &User{}
	err := r.db.Where("alpaca_user_id = ?", alpacaUserID).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	user.AlpacaUserID = alpacaUserID
	user.AlpacaAccessToken = accessToken
	user.AlpacaRefreshToken = refreshToken
	user.TokenExpiry = time.Unix(expiry, 0)

	return r.db.Save(user).Error
}

func (r *authRepository) GetUserByAlpacaID(alpacaUserID string) (*User, error) {
	user := &User{}
	err := r.db.Where("alpaca_user_id = ?", alpacaUserID).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}