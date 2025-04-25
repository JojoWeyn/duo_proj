package postgres_test

import (
	"context"
	"testing"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/repository/db/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.Identity{}))
	return db
}

func TestIdentityRepository_CRUD(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id := &entity.Identity{
		Email:    "user@example.com",
		Login:    "user",
		UserUUID: uuid.New(),
	}

	// Create
	err := repo.Create(ctx, id)
	require.NoError(t, err)
	require.NotZero(t, id.ID)

	// FindByEmail
	foundByEmail, err := repo.FindByEmail(ctx, "user@example.com")
	require.NoError(t, err)
	require.Equal(t, id.ID, foundByEmail.ID)

	// FindByUUID
	foundByUUID, err := repo.FindByUUID(ctx, id.UserUUID.String())
	require.NoError(t, err)
	require.Equal(t, id.Email, foundByUUID.Email)

	// Update
	id.Login = "newlogin"
	err = repo.Update(ctx, id)
	require.NoError(t, err)

	foundByLogin, err := repo.FindByLogin(ctx, "newlogin")
	require.NoError(t, err)
	require.Equal(t, id.Email, foundByLogin.Email)

	// Delete
	err = repo.Delete(ctx, int(id.ID))
	require.NoError(t, err)

	deleted, err := repo.FindByEmail(ctx, "user@example.com")
	require.Error(t, err)
	require.Nil(t, deleted)
}
