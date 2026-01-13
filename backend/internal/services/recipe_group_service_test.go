package services

import (
	"testing"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
	testutil "github.com/homecooking/backend/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeGroup_Create(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name:        "Comfort Food",
		Slug:        "comfort-food",
		Icon:        stringPtr("🍲"),
		Description: stringPtr("Classic comfort dishes"),
	}
	created, err := service.Create(group)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Comfort Food", created.Name)
	assert.Equal(t, "🍲", *created.Icon)
}

func TestRecipeGroup_Create_AutoSlug(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name:        "Weekend Meals",
		Description: stringPtr("Meals for the weekend"),
	}

	created, err := service.Create(group)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "weekend-meals", created.Slug)
}

func TestRecipeGroup_GetByID(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name: "Comfort Food",
		Slug: "comfort-food",
	}
	created, err := groupRepo.Create(group)
	require.NoError(t, err)

	fetched, err := service.GetByID(created.ID.String())
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Comfort Food", fetched.Name)
}

func TestRecipeGroup_GetByID_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	_, err = service.GetByID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestRecipeGroup_GetBySlug(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name: "Comfort Food",
		Slug: "comfort-food",
	}
	created, err := groupRepo.Create(group)
	require.NoError(t, err)

	fetched, err := service.GetBySlug("comfort-food")
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "comfort-food", fetched.Slug)
}

func TestRecipeGroup_GetBySlug_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	_, err = service.GetBySlug("non-existent-slug")
	assert.Error(t, err)
}

func TestRecipeGroup_List(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group1 := &models.RecipeGroup{
		Name: "Comfort Food",
		Slug: "comfort-food",
	}
	group2 := &models.RecipeGroup{
		Name: "Quick Meals",
		Slug: "quick-meals",
	}

	_, err = groupRepo.Create(group1)
	require.NoError(t, err)
	_, err = groupRepo.Create(group2)
	require.NoError(t, err)

	groups, err := service.List()
	require.NoError(t, err)
	assert.Len(t, groups, 2)
}

func TestRecipeGroup_List_Empty(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	groups, err := service.List()
	require.NoError(t, err)
	assert.Len(t, groups, 0)
}

func TestRecipeGroup_Update(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name: "Comfort Food",
		Slug: "comfort-food",
	}
	created, err := groupRepo.Create(group)
	require.NoError(t, err)

	updatedGroup := &models.RecipeGroup{
		Name:        "Cozy Food",
		Description: stringPtr("Updated description"),
	}

	updated, err := service.Update(created.ID.String(), updatedGroup)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Cozy Food", updated.Name)
	assert.Equal(t, "cozy-food", updated.Slug)
}

func TestRecipeGroup_Update_NotFound(t *testing.T) {
	t.Skip("Skip: SQLite COALESCE query issue with sqlc.narg - needs PostgreSQL-specific query handling")

	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	updatedGroup := &models.RecipeGroup{
		Name: "Updated Name",
	}

	_, err = service.Update("00000000-0000-0000-0000-000000000000", updatedGroup)
	assert.Error(t, err)
}

func TestRecipeGroup_Delete(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	group := &models.RecipeGroup{
		Name: "Comfort Food",
		Slug: "comfort-food",
	}
	created, err := groupRepo.Create(group)
	require.NoError(t, err)

	err = service.Delete(created.ID.String())
	require.NoError(t, err)

	_, err = service.GetByID(created.ID.String())
	assert.Error(t, err)
}

func TestRecipeGroup_Delete_NotFound(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	err = service.Delete("00000000-0000-0000-0000-000000000000")
	assert.NoError(t, err)
}

func TestRecipeGroup_AddRecipeToGroup(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	recipe := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     true,
	}
	createdRecipe, err := recipeRepo.Create(recipe)
	require.NoError(t, err)

	group := &models.RecipeGroup{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	createdGroup, err := groupRepo.Create(group)
	require.NoError(t, err)

	err = service.AddRecipeToGroup(createdGroup.ID.String(), createdRecipe.ID.String())
	require.NoError(t, err)

	recipes, err := service.GetRecipesInGroup(createdGroup.ID.String())
	require.NoError(t, err)
	assert.Len(t, recipes, 1)
	assert.Equal(t, createdRecipe.ID, recipes[0].ID)
}

func TestRecipeGroup_RemoveRecipeFromGroup(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	recipe := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     true,
	}
	createdRecipe, err := recipeRepo.Create(recipe)
	require.NoError(t, err)

	group := &models.RecipeGroup{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	createdGroup, err := groupRepo.Create(group)
	require.NoError(t, err)

	err = service.AddRecipeToGroup(createdGroup.ID.String(), createdRecipe.ID.String())
	require.NoError(t, err)

	err = service.RemoveRecipeFromGroup(createdGroup.ID.String(), createdRecipe.ID.String())
	require.NoError(t, err)

	recipes, err := service.GetRecipesInGroup(createdGroup.ID.String())
	require.NoError(t, err)
	assert.Len(t, recipes, 0)
}

func TestRecipeGroup_GetRecipesInGroup(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	recipeRepo := repository.NewRecipeRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	recipe1 := &models.Recipe{
		Title:           "Pancakes",
		Slug:            "pancakes",
		MarkdownContent: "Ingredients: Flour, Eggs. Instructions: Mix and cook.",
		Description:     stringPtr("Simple pancakes"),
		IsPublished:     true,
	}
	recipe2 := &models.Recipe{
		Title:           "Omelette",
		Slug:            "omelette",
		MarkdownContent: "Ingredients: Eggs. Instructions: Cook eggs.",
		Description:     stringPtr("Simple omelette"),
		IsPublished:     true,
	}
	createdRecipe1, err := recipeRepo.Create(recipe1)
	require.NoError(t, err)
	createdRecipe2, err := recipeRepo.Create(recipe2)
	require.NoError(t, err)

	group := &models.RecipeGroup{
		Name: "Breakfast",
		Slug: "breakfast",
	}
	createdGroup, err := groupRepo.Create(group)
	require.NoError(t, err)

	err = service.AddRecipeToGroup(createdGroup.ID.String(), createdRecipe1.ID.String())
	require.NoError(t, err)
	err = service.AddRecipeToGroup(createdGroup.ID.String(), createdRecipe2.ID.String())
	require.NoError(t, err)

	recipes, err := service.GetRecipesInGroup(createdGroup.ID.String())
	require.NoError(t, err)
	assert.Len(t, recipes, 2)
}

func TestRecipeGroup_GenerateSlug(t *testing.T) {
	db, q, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.TeardownTestDB(db)

	groupRepo := repository.NewRecipeGroupRepository(db, q)
	service := NewRecipeGroupService(groupRepo)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "Breakfast", "breakfast"},
		{"with spaces", "Comfort Food", "comfort-food"},
		{"with special chars", "Sunday's Best!", "sunday-s-best"},
		{"already lowercase", "dinner", "dinner"},
		{"with multiple spaces", "  Quick  Meals  ", "quick--meals"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := service.GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, slug)
		})
	}
}
