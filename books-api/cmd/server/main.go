package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"books/internal/cache"
	"books/internal/config"
	"books/internal/database"
	"books/internal/handlers"
	authMw "books/internal/middleware"
	"books/internal/repository"
	svc "books/internal/services"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := database.Migrate(pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	var redisCache *cache.Cache
	if cfg.Cache.Enabled {
		redisAddr := fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port)
		redisCache, err = cache.NewClient(redisAddr)
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v. Caching disabled.", err)
		} else {
			log.Printf("Connected to Redis at %s", redisAddr)
			defer redisCache.Close()
		}
	}

	userRepo := repository.NewUserRepository(pool)
	authorRepo := repository.NewAuthorRepository(pool)
	bookRepo := repository.NewBookRepository(pool, authorRepo)
	progressRepo := repository.NewReadingProgressRepository(pool)

	bookService := svc.NewBookService(bookRepo, authorRepo, redisCache)
	readingService := svc.NewReadingService(progressRepo, bookRepo)
	goodreadsService := svc.NewGoodreadsService(bookService)
	authService := svc.NewAuthService(userRepo, cfg.JWT)

	bookHandler := handlers.NewBookHandler(bookService)
	readingHandler := handlers.NewReadingHandler(readingService)
	goodreadsHandler := handlers.NewGoodreadsHandler(goodreadsService)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/books", func(r chi.Router) {
		r.Post("/", bookHandler.Create)
		r.Post("/import", goodreadsHandler.Import)
		r.Get("/", bookHandler.List)
		r.Get("/{id}", bookHandler.Get)
		r.Put("/{id}", bookHandler.Update)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMw.AuthMiddleware(authService))

		r.Route("/me/books", func(r chi.Router) {
			r.Get("/", readingHandler.GetUserBooks)
			r.Route("/{bookId}", func(r chi.Router) {
				r.Get("/status", readingHandler.GetStatus)
				r.Put("/status", readingHandler.UpdateStatus)
			})
		})
	})

	log.Printf("Server starting on :%d", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.App.Port), r))
}
