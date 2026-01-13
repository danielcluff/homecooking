package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/homecooking/backend/internal/config"
	"github.com/homecooking/backend/internal/db"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/handlers"
	"github.com/homecooking/backend/internal/middleware"
	"github.com/homecooking/backend/internal/repository"
	"github.com/homecooking/backend/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	q := sqlc.New(database)

	userRepo := repository.NewUserRepository(database.DB, q)
	recipeRepo := repository.NewRecipeRepository(database.DB, q)
	categoryRepo := repository.NewCategoryRepository(database.DB, q)
	tagRepo := repository.NewTagRepository(database.DB, q)
	recipeGroupRepo := repository.NewRecipeGroupRepository(database.DB, q)
	shareCodeRepo := repository.NewShareCodeRepository(database.DB, q)
	userInviteRepo := repository.NewUserInviteRepository(database.DB, q)
	variationRepo := repository.NewVariationRepository(database.DB, q)

	authService := services.NewAuthService(cfg, userRepo)
	recipeService := services.NewRecipeService(recipeRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)
	recipeGroupService := services.NewRecipeGroupService(recipeGroupRepo)
	shareCodeService := services.NewShareCodeService(shareCodeRepo, recipeRepo)
	userInviteService := services.NewUserInviteService(userInviteRepo, userRepo)
	variationService := services.NewVariationService(variationRepo, recipeRepo)
	storageService := services.NewStorageService(cfg.Storage.LocalPath, cfg.Storage.MaxFileSize)
	aiService := services.NewAIService(cfg)

	storageService.EnsureDirectory()

	authHandler := handlers.NewAuthHandler(authService)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	tagHandler := handlers.NewTagHandler(tagService)
	recipeGroupHandler := handlers.NewRecipeGroupHandler(recipeGroupService)
	shareCodeHandler := handlers.NewShareCodeHandler(shareCodeService)
	userInviteHandler := handlers.NewUserInviteHandler(userInviteService)
	variationHandler := handlers.NewVariationHandler(variationService)
	uploadHandler := handlers.NewUploadHandler(storageService)
	aiHandler := handlers.NewAIHandler(aiService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)

	authenticated := authMiddleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHandler.Me(w, r)
	}))
	mux.Handle("GET /api/v1/auth/me", authenticated)

	mux.HandleFunc("GET /api/v1/recipes", recipeHandler.ListRecipes)
	mux.HandleFunc("GET /api/v1/recipes/search", recipeHandler.SearchRecipes)
	mux.HandleFunc("GET /api/v1/recipes/{id}", recipeHandler.GetRecipe)
	// Slug endpoint moved after to avoid conflicts
	// mux.HandleFunc("GET /api/v1/recipes/slug/{slug}", recipeHandler.GetRecipeBySlug)

	mux.Handle("POST /api/v1/recipes", authMiddleware.Auth(http.HandlerFunc(recipeHandler.CreateRecipe)))
	mux.Handle("PUT /api/v1/recipes/{id}", authMiddleware.Auth(http.HandlerFunc(recipeHandler.UpdateRecipe)))
	mux.Handle("POST /api/v1/recipes/{id}/publish", authMiddleware.Auth(http.HandlerFunc(recipeHandler.PublishRecipe)))
	mux.Handle("DELETE /api/v1/recipes/{id}", authMiddleware.Auth(http.HandlerFunc(recipeHandler.DeleteRecipe)))

	// Variation routes
	mux.HandleFunc("GET /api/v1/recipes/{id}/variations", variationHandler.ListVariations)
	mux.HandleFunc("GET /api/v1/recipes/{id}/variations/{variationId}", variationHandler.GetVariation)
	mux.Handle("POST /api/v1/recipes/{id}/variations", authMiddleware.Auth(http.HandlerFunc(variationHandler.CreateVariation)))
	mux.Handle("PUT /api/v1/recipes/{id}/variations/{variationId}", authMiddleware.Auth(http.HandlerFunc(variationHandler.UpdateVariation)))
	mux.Handle("DELETE /api/v1/recipes/{id}/variations/{variationId}", authMiddleware.Auth(http.HandlerFunc(variationHandler.DeleteVariation)))

	// Category routes
	mux.HandleFunc("GET /api/v1/categories", categoryHandler.ListCategories)
	mux.HandleFunc("GET /api/v1/categories/{id}", categoryHandler.GetCategory)
	mux.Handle("POST /api/v1/categories", authMiddleware.Auth(http.HandlerFunc(categoryHandler.CreateCategory)))
	mux.Handle("PUT /api/v1/categories/{id}", authMiddleware.Auth(http.HandlerFunc(categoryHandler.UpdateCategory)))
	mux.Handle("DELETE /api/v1/categories/{id}", authMiddleware.Auth(http.HandlerFunc(categoryHandler.DeleteCategory)))

	// Recipe Group routes
	mux.HandleFunc("GET /api/v1/groups", recipeGroupHandler.ListGroups)
	mux.HandleFunc("GET /api/v1/groups/{id}", recipeGroupHandler.GetGroup)
	mux.Handle("POST /api/v1/groups", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.CreateGroup)))
	mux.Handle("PUT /api/v1/groups/{id}", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.UpdateGroup)))
	mux.Handle("DELETE /api/v1/groups/{id}", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.DeleteGroup)))
	mux.Handle("GET /api/v1/groups/{id}/recipes", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.GetGroupRecipes)))
	mux.Handle("POST /api/v1/groups/{id}/recipes", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.AddRecipeToGroup)))
	mux.Handle("DELETE /api/v1/groups/{id}/recipes/{recipeId}", authMiddleware.Auth(http.HandlerFunc(recipeGroupHandler.RemoveRecipeFromGroup)))

	// Tag routes
	mux.HandleFunc("GET /api/v1/tags", tagHandler.ListTags)
	mux.HandleFunc("GET /api/v1/tags/{id}", tagHandler.GetTag)
	// Recipe tags endpoint commented out due to routing conflicts
	// mux.HandleFunc("GET /api/v1/recipes/{recipeId}/tags", tagHandler.GetRecipeTags)
	mux.Handle("POST /api/v1/tags", authMiddleware.Auth(http.HandlerFunc(tagHandler.CreateTag)))
	mux.Handle("PUT /api/v1/tags/{id}", authMiddleware.Auth(http.HandlerFunc(tagHandler.UpdateTag)))
	mux.Handle("DELETE /api/v1/tags/{id}", authMiddleware.Auth(http.HandlerFunc(tagHandler.DeleteTag)))
	// mux.Handle("POST /api/v1/recipes/{recipeId}/tags/{tagId}", authMiddleware.Auth(http.HandlerFunc(tagHandler.AddTagToRecipe)))
	// mux.Handle("DELETE /api/v1/recipes/{recipeId}/tags/{tagId}", authMiddleware.Auth(http.HandlerFunc(tagHandler.RemoveTagFromRecipe)))

	// Upload routes
	mux.Handle("POST /api/v1/upload/image", authMiddleware.Auth(http.HandlerFunc(uploadHandler.UploadImage)))

	// Share Code routes
	mux.Handle("POST /api/v1/share-codes", authMiddleware.Auth(http.HandlerFunc(shareCodeHandler.CreateShareCode)))
	mux.HandleFunc("GET /api/v1/share-codes/{code}", shareCodeHandler.GetShareCode)
	mux.HandleFunc("GET /api/v1/share-codes/{code}/recipe", shareCodeHandler.AccessRecipeByShareCode)
	mux.HandleFunc("GET /api/v1/recipes/{recipeId}/share-codes", shareCodeHandler.GetShareCodesForRecipe)
	mux.Handle("DELETE /api/v1/share-codes/{id}", authMiddleware.Auth(http.HandlerFunc(shareCodeHandler.DeleteShareCode)))

	// User Invite routes
	mux.Handle("POST /api/v1/invites", authMiddleware.Auth(http.HandlerFunc(userInviteHandler.CreateInvite)))
	mux.HandleFunc("GET /api/v1/invites/{code}", userInviteHandler.GetInvite)
	mux.Handle("GET /api/v1/invites", authMiddleware.Auth(http.HandlerFunc(userInviteHandler.ListInvites)))
	mux.Handle("DELETE /api/v1/invites/{id}", authMiddleware.Auth(http.HandlerFunc(userInviteHandler.DeleteInvite)))
	mux.Handle("POST /api/v1/invites/use", authMiddleware.Auth(http.HandlerFunc(userInviteHandler.UseInvite)))

	// AI routes
	mux.HandleFunc("GET /api/v1/ai/status", aiHandler.CheckEnabled)
	mux.HandleFunc("GET /api/v1/ai/config", aiHandler.GetConfig)
	mux.Handle("POST /api/v1/ai/extract", authMiddleware.Auth(http.HandlerFunc(aiHandler.ExtractFromImage)))
	mux.Handle("POST /api/v1/ai/enhance", authMiddleware.Auth(http.HandlerFunc(aiHandler.EnhanceRecipe)))

	// Static file server for uploads
	fs := http.FileServer(http.Dir(cfg.Storage.LocalPath))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fs))

	handler := middleware.Logging(middleware.CORS(mux))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler,
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
