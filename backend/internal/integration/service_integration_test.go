package integration

import (
	"testing"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	"github.com/homecooking/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func TestRecipeServiceIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := GetTestUser(t, server, "recipe@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))

	createReq := models.CreateRecipeRequest{
		Title:           "Integration Test Recipe",
		MarkdownContent: "Test content for integration",
		Description:     stringPtr("Test description"),
		IsPublished:     false,
	}

	recipe, err := recipeService.CreateRecipe(&createReq, user.ID.String())
	require.NoError(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "Integration Test Recipe", recipe.Title)
	assert.Equal(t, user.ID, *recipe.AuthorID)

	retrieved, err := recipeService.GetRecipe(recipe.ID.String())
	require.NoError(t, err)
	assert.Equal(t, recipe.ID, retrieved.ID)
	assert.Equal(t, recipe.Title, retrieved.Title)
}

func TestRecipeListIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := GetTestUser(t, server, "list@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))

	recipes := []models.CreateRecipeRequest{
		{Title: "Recipe 1", MarkdownContent: "Content 1", IsPublished: true},
		{Title: "Recipe 2", MarkdownContent: "Content 2", IsPublished: true},
		{Title: "Recipe 3", MarkdownContent: "Content 3", IsPublished: true},
	}

	for _, req := range recipes {
		_, err := recipeService.CreateRecipe(&req, user.ID.String())
		require.NoError(t, err)
	}

	list, err := recipeService.ListRecipes(10, 0)
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestRecipeDeleteIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := GetTestUser(t, server, "delete@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))

	createReq := models.CreateRecipeRequest{
		Title:           "Delete Me",
		MarkdownContent: "Content to delete",
		IsPublished:     false,
	}

	created, err := recipeService.CreateRecipe(&createReq, user.ID.String())
	require.NoError(t, err)

	err = recipeService.DeleteRecipe(created.ID.String(), user.ID.String())
	require.NoError(t, err)

	_, err = recipeService.GetRecipe(created.ID.String())
	assert.Error(t, err)
}

func TestRecipeOwnershipIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user1 := GetTestUser(t, server, "owner1@example.com", "password123")
	user2 := GetTestUser(t, server, "owner2@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))

	createReq := models.CreateRecipeRequest{
		Title:           "My Recipe",
		MarkdownContent: "Content",
		IsPublished:     false,
	}

	created, err := recipeService.CreateRecipe(&createReq, user1.ID.String())
	require.NoError(t, err)

	err = recipeService.DeleteRecipe(created.ID.String(), user2.ID.String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestCategoryServiceIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	categoryService := services.NewCategoryService(repository.NewCategoryRepository(server.DB, server.Queries))

	category := &models.Category{
		Name:        "Integration Category",
		Slug:        "integration-category",
		Description: stringPtr("Test category"),
	}

	created, err := categoryService.CreateCategory(category)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Integration Category", created.Name)

	retrieved, err := categoryService.GetCategory(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)

	err = categoryService.DeleteCategory(created.ID.String())
	require.NoError(t, err)

	_, err = categoryService.GetCategory(created.ID.String())
	assert.Error(t, err)
}

func TestTagServiceIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	tagService := services.NewTagService(repository.NewTagRepository(server.DB, server.Queries))

	tag := &models.Tag{
		Name:  "Integration Tag",
		Slug:  "integration-tag",
		Color: "#ff0000",
	}

	created, err := tagService.CreateTag(tag)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Integration Tag", created.Name)

	retrieved, err := tagService.GetTag(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)

	list, err := tagService.ListTags()
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestRecipeGroupServiceIntegration(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	groupService := services.NewRecipeGroupService(repository.NewRecipeGroupRepository(server.DB, server.Queries))

	group := &models.RecipeGroup{
		Name:        "Integration Group",
		Slug:        "integration-group",
		Description: stringPtr("Test group"),
	}

	created, err := groupService.Create(group)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Integration Group", created.Name)

	retrieved, err := groupService.GetByID(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)

	list, err := groupService.List()
	require.NoError(t, err)
	assert.Len(t, list, 1)
}
