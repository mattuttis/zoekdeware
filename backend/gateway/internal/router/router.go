package router

import (
	"github.com/gorilla/mux"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/config"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/handlers"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/middleware"
)

func New(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	r.HandleFunc("/health", handlers.Health).Methods("GET")

	api := r.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	protected.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", handlers.UpdateProfile).Methods("PUT")

	protected.HandleFunc("/discover", handlers.Discover).Methods("GET")
	protected.HandleFunc("/swipe", handlers.Swipe).Methods("POST")
	protected.HandleFunc("/matches", handlers.GetMatches).Methods("GET")

	protected.HandleFunc("/conversations", handlers.GetConversations).Methods("GET")
	protected.HandleFunc("/conversations/{id}", handlers.GetConversation).Methods("GET")
	protected.HandleFunc("/conversations/{id}/messages", handlers.SendMessage).Methods("POST")

	protected.HandleFunc("/media", handlers.UploadMedia).Methods("POST")
	protected.HandleFunc("/media/{id}", handlers.DeleteMedia).Methods("DELETE")

	protected.HandleFunc("/location", handlers.UpdateLocation).Methods("PUT")

	ws := api.PathPrefix("/ws").Subrouter()
	ws.Use(middleware.Auth(cfg.JWTSecret))
	ws.HandleFunc("/chat", handlers.WebSocketChat)

	return r
}
