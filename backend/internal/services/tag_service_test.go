package services

import (
	"testing"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	testutil "github.com/homecooking/backend/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestCreateTag(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	created, err := service.CreateTag(tag)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Vegetarian", created.Name)
	assert.Equal(t, "vegetarian", created.Slug)
	assert.Equal(t, "#6366f1", created.Color)
}

func TestCreateTag_EmptyName(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name: "",
	}

	_, err = service.CreateTag(tag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestCreateTag_AutoColor(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Quick",
		Color: "",
	}

	created, err := service.CreateTag(tag)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "#6366f1", created.Color)
}

func TestGetTag(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	created, err := tagRepo.Create(tag)
	require.NoError(t, err)

	fetched, err := service.GetTag(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Vegetarian", fetched.Name)
}

func TestGetTag_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	_, err = service.GetTag("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestListTags(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag1 := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	tag2 := &models.Tag{
		Name:  "Quick",
		Slug:  "quick",
		Color: "#10b981",
	}
	tag3 := &models.Tag{
		Name:  "Spicy",
		Slug:  "spicy",
		Color: "#ef4444",
	}

	_, err = tagRepo.Create(tag1)
	require.NoError(t, err)
	_, err = tagRepo.Create(tag2)
	require.NoError(t, err)
	_, err = tagRepo.Create(tag3)
	require.NoError(t, err)

	tags, err := service.ListTags()
	require.NoError(t, err)
	assert.Len(t, tags, 3)
}

func TestUpdateTag(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	created, err := tagRepo.Create(tag)
	require.NoError(t, err)

	updatedTag := &models.Tag{
		Name:  "Vegan",
		Color: "#10b981",
	}

	updated, err := service.UpdateTag(created.ID.String(), updatedTag)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Vegan", updated.Name)
	assert.Equal(t, "#10b981", updated.Color)
}

func TestUpdateTag_EmptyName(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name: "Vegetarian",
	}
	created, err := tagRepo.Create(tag)
	require.NoError(t, err)

	updatedTag := &models.Tag{
		Name: "",
	}

	_, err = service.UpdateTag(created.ID.String(), updatedTag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestUpdateTag_AutoColor(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Color: "#6366f1",
	}
	created, err := tagRepo.Create(tag)
	require.NoError(t, err)

	updatedTag := &models.Tag{
		Name:  "Vegetarian",
		Color: "",
	}

	updated, err := service.UpdateTag(created.ID.String(), updatedTag)
	require.NoError(t, err)
	assert.Equal(t, "#6366f1", updated.Color)
}

func TestDeleteTag(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	created, err := tagRepo.Create(tag)
	require.NoError(t, err)

	err = service.DeleteTag(created.ID.String())
	require.NoError(t, err)

	_, err = service.GetTag(created.ID.String())
	assert.Error(t, err)
}

func TestDeleteTag_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	service := NewTagService(tagRepo)

	err = service.DeleteTag("00000000-0000-0000-0000-000000000000")
	assert.NoError(t, err)
}

func TestAddTagToRecipe(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewTagService(tagRepo)

	recipe := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     false,
	}
	createdRecipe, err := recipeRepo.Create(recipe)
	require.NoError(t, err)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	createdTag, err := tagRepo.Create(tag)
	require.NoError(t, err)

	err = service.AddTagToRecipe(createdRecipe.ID.String(), createdTag.ID.String())
	require.NoError(t, err)

	tags, err := service.GetRecipeTags(createdRecipe.ID.String())
	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, createdTag.ID, tags[0].ID)
}

func TestRemoveTagFromRecipe(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewTagService(tagRepo)

	recipe := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     false,
	}
	createdRecipe, err := recipeRepo.Create(recipe)
	require.NoError(t, err)

	tag := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	createdTag, err := tagRepo.Create(tag)
	require.NoError(t, err)

	err = service.AddTagToRecipe(createdRecipe.ID.String(), createdTag.ID.String())
	require.NoError(t, err)

	err = service.RemoveTagFromRecipe(createdRecipe.ID.String(), createdTag.ID.String())
	require.NoError(t, err)

	tags, err := service.GetRecipeTags(createdRecipe.ID.String())
	require.NoError(t, err)
	assert.Len(t, tags, 0)
}

func TestGetRecipeTags(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	tagRepo := repository.NewTagRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewTagService(tagRepo)

	recipe := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     false,
	}
	createdRecipe, err := recipeRepo.Create(recipe)
	require.NoError(t, err)

	tag1 := &models.Tag{
		Name:  "Vegetarian",
		Slug:  "vegetarian",
		Color: "#6366f1",
	}
	tag2 := &models.Tag{
		Name:  "Quick",
		Slug:  "quick",
		Color: "#10b981",
	}

	createdTag1, err := tagRepo.Create(tag1)
	require.NoError(t, err)
	createdTag2, err := tagRepo.Create(tag2)
	require.NoError(t, err)

	err = service.AddTagToRecipe(createdRecipe.ID.String(), createdTag1.ID.String())
	require.NoError(t, err)
	err = service.AddTagToRecipe(createdRecipe.ID.String(), createdTag2.ID.String())
	require.NoError(t, err)

	tags, err := service.GetRecipeTags(createdRecipe.ID.String())
	require.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestGenerateTagSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "Vegetarian", "vegetarian"},
		{"with spaces", "Super Quick", "super-quick"},
		{"with apostrophe", "Chef's Choice", "chefs-choice"},
		{"with quotes", "\"Special\" Dish", "special-dish"},
		{"already lowercase", "dinner", "dinner"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := generateTagSlug(tt.input)
			assert.Equal(t, tt.expected, slug)
		})
	}
}
