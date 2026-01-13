package services

import (
	"testing"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	testutil "github.com/homecooking/backend/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func TestCreateCategory(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := testutil.TestCategory1

	created, err := service.CreateCategory(category)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Breakfast", created.Name)
	assert.Equal(t, "breakfast", created.Slug)
	assert.Equal(t, "🍳", *created.Icon)
}

func TestCreateCategory_EmptyName(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "",
	}

	_, err = service.CreateCategory(category)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestCreateCategory_AutoSlug(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Lunch and Dinner",
		Slug: "",
	}

	created, err := service.CreateCategory(category)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "lunch-and-dinner", created.Slug)
}

func TestGetCategory(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Breakfast",
		Slug: "breakfast",
		Icon: stringPtr("🍳"),
	}
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	fetched, err := service.GetCategory(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Breakfast", fetched.Name)
	assert.Equal(t, "breakfast", fetched.Slug)
}

func TestGetCategory_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	_, err = service.GetCategory("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestGetCategoryBySlug(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := testutil.TestCategory1
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	fetched, err := service.GetBySlug("breakfast")
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "breakfast", fetched.Slug)
}

func TestListCategories(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category1 := testutil.TestCategory1
	category2 := testutil.TestCategory2
	category3 := testutil.TestCategory3

	_, err = categoryRepo.Create(category1)
	require.NoError(t, err)
	_, err = categoryRepo.Create(category2)
	require.NoError(t, err)
	_, err = categoryRepo.Create(category3)
	require.NoError(t, err)

	categories, err := service.ListCategories()
	require.NoError(t, err)
	assert.Len(t, categories, 3)
}

func TestListCategories_Empty(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	categories, err := service.ListCategories()
	require.NoError(t, err)
	assert.Len(t, categories, 0)
}

func TestUpdateCategory(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Breakfast",
		Slug: "breakfast",
		Icon: stringPtr("🍳"),
	}
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	updatedCategory := &models.Category{
		Name: "Brunch",
		Slug: "brunch",
		Icon: stringPtr("🥐"),
	}

	updated, err := service.UpdateCategory(created.ID.String(), updatedCategory)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Brunch", updated.Name)
	assert.Equal(t, "brunch", updated.Slug)
	assert.Equal(t, "🥐", *updated.Icon)
}

func TestUpdateCategory_EmptyName(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	updatedCategory := &models.Category{
		Name: "",
	}

	_, err = service.UpdateCategory(created.ID.String(), updatedCategory)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestUpdateCategory_AutoSlug(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	updatedCategory := &models.Category{
		Name: "Morning Food",
	}

	updated, err := service.UpdateCategory(created.ID.String(), updatedCategory)
	require.NoError(t, err)
	assert.Equal(t, "morning-food", updated.Slug)
}

func TestUpdateCategory_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	updatedCategory := &models.Category{
		Name: "Updated Name",
	}

	_, err = service.UpdateCategory("00000000-0000-0000-0000-000000000000", updatedCategory)
	assert.Error(t, err)
}

func TestDeleteCategory(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	created, err := categoryRepo.Create(category)
	require.NoError(t, err)

	err = service.DeleteCategory(created.ID.String())
	require.NoError(t, err)

	_, err = service.GetCategory(created.ID.String())
	assert.Error(t, err)
}

func TestDeleteCategory_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	categoryRepo := repository.NewCategoryRepository(db, q)
	service := NewCategoryService(categoryRepo)

	err = service.DeleteCategory("00000000-0000-0000-0000-000000000000")
	assert.NoError(t, err)
}

func TestGenerateCategorySlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "Breakfast", "breakfast"},
		{"with spaces", "Lunch and Dinner", "lunch-and-dinner"},
		{"with apostrophe", "Chef's Choice", "chefs-choice"},
		{"with quotes", "\"Special\" Dish", "special-dish"},
		{"already lowercase", "dinner", "dinner"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := generateCategorySlug(tt.input)
			assert.Equal(t, tt.expected, slug)
		})
	}
}
