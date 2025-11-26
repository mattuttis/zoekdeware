package router

import (
	"github.com/gorilla/mux"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/config"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/handlers"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/middleware"
)

func New(cfg *config.Config, h *handlers.Handlers) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	r.HandleFunc("/health", h.Health).Methods("GET")

	api := r.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", h.Register).Methods("POST")
	auth.HandleFunc("/login", h.Login).Methods("POST")
	auth.HandleFunc("/refresh", h.RefreshToken).Methods("POST")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	protected.HandleFunc("/profile", h.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", h.UpdateProfile).Methods("PUT")

	protected.HandleFunc("/discover", h.Discover).Methods("GET")
	protected.HandleFunc("/swipe", h.Swipe).Methods("POST")
	protected.HandleFunc("/matches", h.GetMatches).Methods("GET")

	protected.HandleFunc("/conversations", h.GetConversations).Methods("GET")
	protected.HandleFunc("/conversations/{id}", h.GetConversation).Methods("GET")
	protected.HandleFunc("/conversations/{id}/messages", h.SendMessage).Methods("POST")

	protected.HandleFunc("/media", h.UploadMedia).Methods("POST")
	protected.HandleFunc("/media/{id}", h.DeleteMedia).Methods("DELETE")

	protected.HandleFunc("/location", h.UpdateLocation).Methods("PUT")

	ws := api.PathPrefix("/ws").Subrouter()
	ws.Use(middleware.Auth(cfg.JWTSecret))
	ws.HandleFunc("/chat", h.WebSocketChat)

	return r
}
