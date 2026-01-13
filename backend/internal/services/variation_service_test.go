package services

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	testutil "github.com/homecooking/backend/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func vCreateTestUser(db *sql.DB, q *sqlc.Queries, email string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	userID := uuid.New().String()

	_, err := db.Exec(`
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, userID, email, string(hashedPassword), "user")

	if err != nil {
		return ""
	}

	return userID
}

func vCreateTestRecipe(db *sql.DB, q *sqlc.Queries, authorID string) string {
	recipeID := uuid.New().String()

	_, err := db.Exec(`
		INSERT INTO recipes (id, title, slug, markdown_content, author_id, is_published, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 0, datetime('now'), datetime('now'))
	`, recipeID, "Pancakes", "pancakes", "Ingredients: Flour, Eggs.", authorID)

	if err != nil {
		return ""
	}

	return recipeID
}

func vInt32Ptr(i int32) *int32 {
	return &i
}

func vStringPtr(s string) *string {
	return &s
}

func vBoolPtr(b bool) *bool {
	return &b
}

func TestVariationService_CreateVariation(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "My variation: Use 2 cups of flour instead of 1",
		PrepTimeMinutes: vInt32Ptr(15),
		CookTimeMinutes: vInt32Ptr(20),
		Servings:        vInt32Ptr(6),
		Difficulty:      vStringPtr("medium"),
		Notes:           vStringPtr("Makes fluffier pancakes"),
		IsPublished:     true,
	}

	variation, err := service.CreateVariation(req, recipeID, authorID)
	require.NoError(t, err)
	assert.NotNil(t, variation)
	assert.Equal(t, "My variation: Use 2 cups of flour instead of 1", variation.MarkdownContent)
	assert.Equal(t, int32(15), *variation.PrepTimeMinutes)
	assert.Equal(t, int32(20), *variation.CookTimeMinutes)
	assert.Equal(t, int32(6), *variation.Servings)
	assert.Equal(t, "medium", *variation.Difficulty)
	assert.Equal(t, "Makes fluffier pancakes", *variation.Notes)
	assert.True(t, variation.IsPublished)
}

func TestVariationService_CreateVariation_EmptyContent(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "",
		IsPublished:     false,
	}

	_, err = service.CreateVariation(req, recipeID, authorID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "markdown content is required")
}

func TestVariationService_GetVariation(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "My variation",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID)
	require.NoError(t, err)

	fetched, err := service.GetVariation(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "My variation", fetched.MarkdownContent)
}

func TestVariationService_GetVariation_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	_, err = service.GetVariation("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestVariationService_GetVariationsByRecipe(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID1 := vCreateTestUser(db, q, "user1@example.com")
	authorID2 := vCreateTestUser(db, q, "user2@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID1)

	req1 := &models.CreateVariationRequest{
		MarkdownContent: "User 1's variation",
		IsPublished:     true,
	}

	req2 := &models.CreateVariationRequest{
		MarkdownContent: "User 2's variation",
		IsPublished:     false,
	}

	_, err = service.CreateVariation(req1, recipeID, authorID1)
	require.NoError(t, err)

	_, err = service.CreateVariation(req2, recipeID, authorID2)
	require.NoError(t, err)

	variations, err := service.GetVariationsByRecipe(recipeID)
	require.NoError(t, err)
	assert.Len(t, variations, 2)
}

func TestVariationService_GetPublishedVariationsByRecipe(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID1 := vCreateTestUser(db, q, "user1@example.com")
	authorID2 := vCreateTestUser(db, q, "user2@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID1)

	req1 := &models.CreateVariationRequest{
		MarkdownContent: "Published variation",
		IsPublished:     true,
	}

	req2 := &models.CreateVariationRequest{
		MarkdownContent: "Draft variation",
		IsPublished:     false,
	}

	_, err = service.CreateVariation(req1, recipeID, authorID1)
	require.NoError(t, err)

	_, err = service.CreateVariation(req2, recipeID, authorID2)
	require.NoError(t, err)

	variations, err := service.GetPublishedVariationsByRecipe(recipeID)
	require.NoError(t, err)
	assert.Len(t, variations, 1)
	assert.Equal(t, "Published variation", variations[0].MarkdownContent)
}

func TestVariationService_GetVariationByRecipeAndAuthor(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "My variation",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID)
	require.NoError(t, err)

	fetched, err := service.GetVariationByRecipeAndAuthor(recipeID, authorID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "My variation", fetched.MarkdownContent)
}

func TestVariationService_UpdateVariation(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID)
	require.NoError(t, err)

	updateReq := &models.UpdateVariationRequest{
		MarkdownContent: vStringPtr("Updated content"),
		PrepTimeMinutes: vInt32Ptr(25),
		Notes:           vStringPtr("Updated notes"),
		IsPublished:     vBoolPtr(false),
	}

	updated, err := service.UpdateVariation(created.ID.String(), updateReq, authorID)
	require.NoError(t, err)
	assert.Equal(t, "Updated content", updated.MarkdownContent)
	assert.Equal(t, int32(25), *updated.PrepTimeMinutes)
	assert.Equal(t, "Updated notes", *updated.Notes)
	assert.False(t, updated.IsPublished)
}

func TestVariationService_UpdateVariation_NotOwner(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID1 := vCreateTestUser(db, q, "user1@example.com")
	authorID2 := vCreateTestUser(db, q, "user2@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID1)

	req := &models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID1)
	require.NoError(t, err)

	updateReq := &models.UpdateVariationRequest{
		MarkdownContent: vStringPtr("Hacked content"),
	}

	_, err = service.UpdateVariation(created.ID.String(), updateReq, authorID2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestVariationService_DeleteVariation(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID)
	require.NoError(t, err)

	err = service.DeleteVariation(created.ID.String(), authorID)
	require.NoError(t, err)

	_, err = service.GetVariation(created.ID.String())
	assert.Error(t, err)
}

func TestVariationService_DeleteVariation_NotOwner(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID1 := vCreateTestUser(db, q, "user1@example.com")
	authorID2 := vCreateTestUser(db, q, "user2@example.com")
	recipeID := vCreateTestRecipe(db, q, authorID1)

	req := &models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := service.CreateVariation(req, recipeID, authorID1)
	require.NoError(t, err)

	err = service.DeleteVariation(created.ID.String(), authorID2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestVariationService_ListVariationsByAuthor(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	recipeRepo := repository.NewRecipeRepository(db, q)
	variationRepo := repository.NewVariationRepository(db, q)
	service := NewVariationService(variationRepo, recipeRepo)

	authorID := vCreateTestUser(db, q, "test@example.com")

	recipeID1 := vCreateTestRecipe(db, q, authorID)
	recipeID2 := vCreateTestRecipe(db, q, authorID)

	req := &models.CreateVariationRequest{
		MarkdownContent: "My variation",
		IsPublished:     true,
	}

	_, err = service.CreateVariation(req, recipeID1, authorID)
	require.NoError(t, err)

	_, err = service.CreateVariation(req, recipeID2, authorID)
	require.NoError(t, err)

	variations, err := service.ListVariationsByAuthor(authorID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, variations, 2)
}
