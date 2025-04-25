package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/repository/db/postgres"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTokenTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.BlacklistedToken{}))
	return db
}

func closeDB(t *testing.T, db *gorm.DB) *gorm.DB {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	require.NoError(t, err)
	return db
}

// ---- BlacklistToken ----

func TestBlacklistToken_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	token := &entity.BlacklistedToken{
		Token:     "token123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	err := repo.BlacklistToken(ctx, token)
	require.NoError(t, err)
}

func TestBlacklistToken_Fail_DBClosed(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	// Закрываем соединение
	closeDB(t, db)

	token := &entity.BlacklistedToken{
		Token:     "token123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	err := repo.BlacklistToken(ctx, token)
	require.Error(t, err)
}

// ---- IsBlacklisted ----

func TestIsBlacklisted_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	token := &entity.BlacklistedToken{
		Token:     "blocked-token",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	_ = repo.BlacklistToken(ctx, token)

	isBlacklisted, err := repo.IsBlacklisted(ctx, "blocked-token")
	require.NoError(t, err)
	require.True(t, isBlacklisted)
}

func TestIsBlacklisted_Fail_DBClosed(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	closeDB(t, db)

	isBlacklisted, err := repo.IsBlacklisted(ctx, "blocked-token")
	require.Error(t, err)
	require.False(t, isBlacklisted)
}

// ---- CleanupExpired ----

func TestCleanupExpired_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	expired := &entity.BlacklistedToken{
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	valid := &entity.BlacklistedToken{
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	_ = repo.BlacklistToken(ctx, expired)
	_ = repo.BlacklistToken(ctx, valid)

	err := repo.CleanupExpired(ctx)
	require.NoError(t, err)

	isExpired, _ := repo.IsBlacklisted(ctx, "expired-token")
	isValid, _ := repo.IsBlacklisted(ctx, "valid-token")

	require.False(t, isExpired)
	require.True(t, isValid)
}

func TestCleanupExpired_Fail_DBClosed(t *testing.T) {
	ctx := context.TODO()
	db := setupTokenTestDB(t)
	repo := postgres.NewTokenRepository(db)

	closeDB(t, db)

	err := repo.CleanupExpired(ctx)
	require.Error(t, err)
}
