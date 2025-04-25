package postgres_test

import (
	"context"
	"testing"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/repository/db/postgres"
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

// --- Create ---

func TestCreateIdentity_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id, err := entity.NewIdentity("user@example.com", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, id)
	require.NoError(t, err)
	require.NotZero(t, id.ID)
}

// --- FindByEmail ---

func TestFindByEmail_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id, _ := entity.NewIdentity("user@example.com", "password")
	repo.Create(ctx, id)

	result, err := repo.FindByEmail(ctx, "user@example.com")
	require.NoError(t, err)
	require.Equal(t, id.ID, result.ID)
}

func TestFindByEmail_NotFound(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	result, err := repo.FindByEmail(ctx, "notfound@example.com")
	require.Error(t, err)
	require.Nil(t, result)
}

// --- FindByUUID ---

func TestFindByUUID_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id, _ := entity.NewIdentity("user@example.com", "password")
	repo.Create(ctx, id)

	result, err := repo.FindByUUID(ctx, id.UserUUID.String())
	require.NoError(t, err)
	require.Equal(t, id.Email, result.Email)
}

func TestFindByUUID_NotFound(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	result, err := repo.FindByUUID(ctx, "00000000-0000-0000-0000-000000000000")
	require.Error(t, err)
	require.Nil(t, result)
}

// --- Update ---

func TestUpdateIdentity_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id, _ := entity.NewIdentity("user@example.com", "password")
	repo.Create(ctx, id)

	id.Email = "new@example.com"
	err := repo.Update(ctx, id)
	require.NoError(t, err)

	updated, _ := repo.FindByEmail(ctx, "new@example.com")
	require.Equal(t, id.ID, updated.ID)
}

func TestUpdateIdentity_Fail_NoID(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id := &entity.Identity{Email: "broken@example.com"}

	err := repo.Update(ctx, id)
	require.Error(t, err)
}

// --- Delete ---

func TestDeleteIdentity_Success(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	id, _ := entity.NewIdentity("user@example.com", "password")
	repo.Create(ctx, id)

	err := repo.Delete(ctx, int(id.ID))
	require.NoError(t, err)

	result, err := repo.FindByEmail(ctx, "user@example.com")
	require.Error(t, err)
	require.Nil(t, result)
}

func TestDeleteIdentity_Fail_InvalidID(t *testing.T) {
	ctx := context.TODO()
	db := setupTestDB(t)
	repo := postgres.NewIdentityRepository(db)

	err := repo.Delete(ctx, -1)
	require.NoError(t, err) // GORM не падает при удалении несуществующего ID
}
