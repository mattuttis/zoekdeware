package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	memberv1 "github.com/mattuttis/inetcontrol/zoekdeware/api/proto/member/v1"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/gateway/internal/middleware"
)

// Handlers holds the gRPC clients for all services.
type Handlers struct {
	memberClient memberv1.MemberServiceClient
	jwtSecret    string
}

// NewHandlers creates a new Handlers instance with the given gRPC clients.
func NewHandlers(memberClient memberv1.MemberServiceClient, jwtSecret string) *Handlers {
	return &Handlers{
		memberClient: memberClient,
		jwtSecret:    jwtSecret,
	}
}

// RegisterRequest represents the JSON request body for registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the JSON response for registration and login.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.memberClient.RegisterMember(ctx, &memberv1.RegisterMemberRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Generate JWT tokens
	authResp, err := h.generateTokens(resp.Member.Id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(authResp)
}

// generateTokens creates access and refresh tokens for a user.
func (h *Handlers) generateTokens(userID string) (*AuthResponse, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Access token expires in 24 hours

	// Create access token
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token (longer expiry)
	refreshExpiresAt := now.Add(7 * 24 * time.Hour) // Refresh token expires in 7 days
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  refreshExpiresAt.Unix(),
		"type": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}

// LoginRequest represents the JSON request body for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.memberClient.AuthenticateMember(ctx, &memberv1.AuthenticateMemberRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Generate JWT tokens
	authResp, err := h.generateTokens(resp.Member.Id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authResp)
}

func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Refresh JWT
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.memberClient.GetMember(ctx, &memberv1.GetMemberRequest{
		MemberId: userID,
	})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.Member)
}

func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to member service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) Discover(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service with location context
	// Return empty response for now
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"profiles": []any{}})
}

func (h *Handlers) Swipe(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) GetMatches(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to matching service
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"matches": []any{}})
}

func (h *Handlers) GetConversations(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"conversations": []any{}})
}

func (h *Handlers) GetConversation(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to messaging service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to media service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to media service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	// TODO: Forward to location service
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) WebSocketChat(w http.ResponseWriter, r *http.Request) {
	// TODO: Upgrade to WebSocket, handle real-time messaging
	w.WriteHeader(http.StatusNotImplemented)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// handleGRPCError converts gRPC errors to HTTP responses.
func handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		writeError(w, http.StatusBadRequest, st.Message())
	case codes.NotFound:
		writeError(w, http.StatusNotFound, st.Message())
	case codes.AlreadyExists:
		writeError(w, http.StatusConflict, st.Message())
	case codes.Unauthenticated:
		writeError(w, http.StatusUnauthorized, st.Message())
	case codes.PermissionDenied:
		writeError(w, http.StatusForbidden, st.Message())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
