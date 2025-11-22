package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/middleware"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to member service
	w.WriteHeader(http.StatusNotImplemented)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to member service, return JWT
	w.WriteHeader(http.StatusNotImplemented)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Refresh JWT
	w.WriteHeader(http.StatusNotImplemented)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	// TODO: Forward to member service
	_ = json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to member service
	w.WriteHeader(http.StatusNotImplemented)
}

func Discover(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service with location context
	w.WriteHeader(http.StatusNotImplemented)
}

func Swipe(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service
	w.WriteHeader(http.StatusNotImplemented)
}

func GetMatches(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service
	w.WriteHeader(http.StatusNotImplemented)
}

func GetConversations(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.WriteHeader(http.StatusNotImplemented)
}

func GetConversation(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.WriteHeader(http.StatusNotImplemented)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.WriteHeader(http.StatusNotImplemented)
}

func UploadMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to media service
	w.WriteHeader(http.StatusNotImplemented)
}

func DeleteMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to media service
	w.WriteHeader(http.StatusNotImplemented)
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to location service
	w.WriteHeader(http.StatusNotImplemented)
}

func WebSocketChat(w http.ResponseWriter, r *http.Request) {
	// TODO: Upgrade to WebSocket, handle real-time messaging
	w.WriteHeader(http.StatusNotImplemented)
}
