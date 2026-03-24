package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"alioth-hrc/internal/auth"
	"alioth-hrc/internal/config"
	"alioth-hrc/internal/model"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const refreshKeyPrefix = "hrc:refresh:"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidRefresh     = errors.New("invalid or expired refresh token")
)

type AuthService struct {
	db      *gorm.DB
	rdb     *redis.Client
	cfg     *config.Config
	access  time.Duration
	refresh time.Duration
}

func NewAuthService(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{
		db:      db,
		rdb:     rdb,
		cfg:     cfg,
		access:  time.Duration(cfg.JWTAccessTTLMin) * time.Minute,
		refresh: time.Duration(cfg.JWTRefreshTTLHours) * time.Hour,
	}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"max=100"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*model.User, error) {
	var n int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", in.Email).Count(&n).Error; err != nil {
		return nil, err
	}
	if n > 0 {
		return nil, ErrEmailTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		Name:         in.Name,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		ledger := &model.Ledger{
			UserID:    u.ID,
			Name:      "默认账本",
			IsDefault: true,
		}
		if err := tx.Create(ledger).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	TokenType        string `json:"token_type"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*TokenPair, error) {
	var u model.User
	if err := s.db.WithContext(ctx).Where("email = ?", in.Email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokens(ctx, u.ID)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint) (*TokenPair, error) {
	access, err := auth.SignAccessToken(userID, s.cfg.JWTSecret, s.access)
	if err != nil {
		return nil, err
	}
	refreshRaw := make([]byte, 32)
	if _, err := rand.Read(refreshRaw); err != nil {
		return nil, err
	}
	refresh := hex.EncodeToString(refreshRaw)
	key := refreshKeyPrefix + refresh
	if err := s.rdb.Set(ctx, key, fmt.Sprintf("%d", userID), s.refresh).Err(); err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		ExpiresIn:        int64(s.access.Seconds()),
		TokenType:        "Bearer",
		RefreshExpiresIn: int64(s.refresh.Seconds()),
	}, nil
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *AuthService) Refresh(ctx context.Context, in RefreshInput) (*TokenPair, error) {
	key := refreshKeyPrefix + in.RefreshToken
	userIDStr, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrInvalidRefresh
	}
	if err != nil {
		return nil, err
	}
	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return nil, ErrInvalidRefresh
	}
	uid := uint(uid64)
	_ = s.rdb.Del(ctx, key).Err()
	return s.issueTokens(ctx, uid)
}

type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *AuthService) Logout(ctx context.Context, in LogoutInput) error {
	key := refreshKeyPrefix + in.RefreshToken
	return s.rdb.Del(ctx, key).Err()
}
