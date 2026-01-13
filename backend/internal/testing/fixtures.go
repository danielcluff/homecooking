package testing

import (
	"time"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
)

var (
	// Test Users
	TestUser1 = &models.User{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:     "test@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	TestAdminUser = &models.User{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Email:     "admin@example.com",
		Role:      "admin",
		CreatedAt: time.Now(),
	}

	TestUser2 = &models.User{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Email:     "editor@example.com",
		Role:      "editor",
		CreatedAt: time.Now(),
	}
)

var (
	// Test Categories
	TestCategory1 = &models.Category{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
		Name:       "Breakfast",
		Slug:       "breakfast",
		Icon:       stringPtr("🍳"),
		OrderIndex: 1,
		CreatedAt:  time.Now(),
	}

	TestCategory2 = &models.Category{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000012"),
		Name:       "Dinner",
		Slug:       "dinner",
		OrderIndex: 2,
		CreatedAt:  time.Now(),
	}

	TestCategory3 = &models.Category{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000013"),
		Name:       "Desserts",
		Slug:       "desserts",
		Icon:       stringPtr("🍰"),
		OrderIndex: 3,
		CreatedAt:  time.Now(),
	}
)

var (
	// Test Tags
	TestTag1 = &models.Tag{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000021"),
		Name:      "Vegetarian",
		Slug:      "vegetarian",
		Color:     "#6366f1",
		CreatedAt: time.Now(),
	}

	TestTag2 = &models.Tag{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000022"),
		Name:      "Quick",
		Slug:      "quick",
		Color:     "#6366f1",
		CreatedAt: time.Now(),
	}

	TestTag3 = &models.Tag{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000023"),
		Name:      "Spicy",
		Slug:      "spicy",
		Color:     "#ef4444",
		CreatedAt: time.Now(),
	}
)

var (
	// Test Recipe Groups
	TestGroup1 = &models.RecipeGroup{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000031"),
		Name:        "Comfort Food",
		Slug:        "comfort-food",
		Icon:        stringPtr("🍲"),
		Description: stringPtr("Classic comfort dishes"),
		CreatedAt:   time.Now(),
	}

	TestGroup2 = &models.RecipeGroup{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000032"),
		Name:        "Biscuits & Gravy",
		Slug:        "biscuits-gravy",
		Icon:        stringPtr("🥞"),
		Description: stringPtr("Southern breakfast combo"),
		CreatedAt:   time.Now(),
	}
)

var (
	// Test Recipes
	TestRecipe1 = &models.Recipe{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000041"),
		Title:           "Fluffy Pancakes",
		Slug:            "fluffy-pancakes",
		MarkdownContent: "## Ingredients\n\n- 2 cups flour\n- 2 eggs\n- 1 cup milk\n\n## Instructions\n\nMix all ingredients and cook on hot griddle.",
		Description:     stringPtr("Fluffy breakfast pancakes"),
		PrepTimeMinutes: int32Ptr(10),
		CookTimeMinutes: int32Ptr(15),
		Servings:        int32Ptr(4),
		Difficulty:      stringPtr("easy"),
		IsPublished:     true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	TestRecipe2 = &models.Recipe{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000042"),
		Title:           "Cheese Omelette",
		Slug:            "cheese-omelette",
		MarkdownContent: "## Ingredients\n\n- 3 eggs\n- 1/2 cup shredded cheese\n- Salt and pepper\n\n## Instructions\n\nBeat eggs, cook in pan, add cheese, fold.",
		Description:     stringPtr("Simple cheese omelette"),
		PrepTimeMinutes: int32Ptr(5),
		CookTimeMinutes: int32Ptr(10),
		Servings:        int32Ptr(2),
		Difficulty:      stringPtr("easy"),
		IsPublished:     false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	TestRecipe3 = &models.Recipe{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000043"),
		Title:           "Spicy Pasta",
		Slug:            "spicy-pasta",
		MarkdownContent: "## Ingredients\n\n- 8 oz pasta\n- 1 can tomato sauce\n- Red pepper flakes\n\n## Instructions\n\nCook pasta, add sauce and spices, serve hot.",
		Description:     stringPtr("Quick spicy pasta dinner"),
		PrepTimeMinutes: int32Ptr(5),
		CookTimeMinutes: int32Ptr(20),
		Servings:        int32Ptr(4),
		Difficulty:      stringPtr("medium"),
		IsPublished:     true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
)

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
