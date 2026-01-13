package integration

import (
	"testing"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	"github.com/homecooking/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariationService_CreateVariation(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := GetTestUser(t, server, "variation@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))
	variationService := services.NewVariationService(
		repository.NewVariationRepository(server.DB, server.Queries),
		repository.NewRecipeRepository(server.DB, server.Queries),
	)

	createRecipeReq := models.CreateRecipeRequest{
		Title:           "Pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Fluffy breakfast pancakes"),
		IsPublished:     false,
	}

	recipe, err := recipeService.CreateRecipe(&createRecipeReq, user.ID.String())
	require.NoError(t, err)
	require.NotNil(t, recipe)

	createVariationReq := models.CreateVariationRequest{
		MarkdownContent: "My variation: Use 2 cups of flour instead of 1",
		PrepTimeMinutes: int32Ptr(15),
		CookTimeMinutes: int32Ptr(20),
		Servings:        int32Ptr(6),
		Difficulty:      stringPtr("medium"),
		Notes:           stringPtr("Makes fluffier pancakes"),
		IsPublished:     true,
	}

	variation, err := variationService.CreateVariation(&createVariationReq, recipe.ID.String(), user.ID.String())
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

func TestVariationService_GetVariationsByRecipe(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user1 := CreateTestUser(t, server, "user1@example.com", "password123")
	user2 := CreateTestUser(t, server, "user2@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))
	variationService := services.NewVariationService(
		repository.NewVariationRepository(server.DB, server.Queries),
		repository.NewRecipeRepository(server.DB, server.Queries),
	)

	createRecipeReq := models.CreateRecipeRequest{
		Title:           "Pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs.",
		IsPublished:     true,
	}

	recipe, err := recipeService.CreateRecipe(&createRecipeReq, user1.ID.String())
	require.NoError(t, err)

	req1 := &models.CreateVariationRequest{
		MarkdownContent: "User 1's variation",
		IsPublished:     true,
	}

	req2 := &models.CreateVariationRequest{
		MarkdownContent: "User 2's variation",
		IsPublished:     false,
	}

	_, err = variationService.CreateVariation(req1, recipe.ID.String(), user1.ID.String())
	require.NoError(t, err)

	_, err = variationService.CreateVariation(req2, recipe.ID.String(), user2.ID.String())
	require.NoError(t, err)

	variations, err := variationService.GetVariationsByRecipe(recipe.ID.String())
	require.NoError(t, err)
	assert.Len(t, variations, 2)
}

func TestVariationService_UpdateVariation(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := CreateTestUser(t, server, "variation@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))
	variationService := services.NewVariationService(
		repository.NewVariationRepository(server.DB, server.Queries),
		repository.NewRecipeRepository(server.DB, server.Queries),
	)

	createRecipeReq := models.CreateRecipeRequest{
		Title:           "Pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs.",
		IsPublished:     true,
	}

	recipe, err := recipeService.CreateRecipe(&createRecipeReq, user.ID.String())
	require.NoError(t, err)

	createVariationReq := models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := variationService.CreateVariation(&createVariationReq, recipe.ID.String(), user.ID.String())
	require.NoError(t, err)

	updateReq := models.UpdateVariationRequest{
		MarkdownContent: stringPtr("Updated content"),
		PrepTimeMinutes: int32Ptr(25),
		Notes:           stringPtr("Updated notes"),
		IsPublished:     boolPtr(false),
	}

	updated, err := variationService.UpdateVariation(created.ID.String(), &updateReq, user.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "Updated content", updated.MarkdownContent)
	assert.Equal(t, int32(25), *updated.PrepTimeMinutes)
	assert.Equal(t, "Updated notes", *updated.Notes)
	assert.False(t, updated.IsPublished)
}

func TestVariationService_DeleteVariation(t *testing.T) {
	server := SetupTestServer(t)
	defer TeardownTestServer(server)

	user := CreateTestUser(t, server, "variation@example.com", "password123")
	recipeService := services.NewRecipeService(repository.NewRecipeRepository(server.DB, server.Queries))
	variationService := services.NewVariationService(
		repository.NewVariationRepository(server.DB, server.Queries),
		repository.NewRecipeRepository(server.DB, server.Queries),
	)

	createRecipeReq := models.CreateRecipeRequest{
		Title:           "Pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs.",
		IsPublished:     true,
	}

	recipe, err := recipeService.CreateRecipe(&createRecipeReq, user.ID.String())
	require.NoError(t, err)

	createVariationReq := models.CreateVariationRequest{
		MarkdownContent: "Original content",
		IsPublished:     true,
	}

	created, err := variationService.CreateVariation(&createVariationReq, recipe.ID.String(), user.ID.String())
	require.NoError(t, err)

	err = variationService.DeleteVariation(created.ID.String(), user.ID.String())
	require.NoError(t, err)

	_, err = variationService.GetVariation(created.ID.String())
	assert.Error(t, err)
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
