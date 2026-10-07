package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"ctracking/backend/internal/auth"
	"ctracking/backend/internal/obs"
	"ctracking/backend/internal/store"

	"github.com/gin-gonic/gin"
)

// AuthHandler serves registration, login, refresh, and the public picker.
type AuthHandler struct {
	store            *store.Store
	tok              *auth.Manager
	keycloakVerifier *auth.KeycloakVerifier

	activationMu sync.Mutex
	activations  map[string]desktopActivation
}

// NewAuthHandler wires the auth handler.
func NewAuthHandler(
	s *store.Store,
	tok *auth.Manager,
	keycloakVerifier *auth.KeycloakVerifier,
) *AuthHandler {
	return &AuthHandler{
		store:            s,
		tok:              tok,
		keycloakVerifier: keycloakVerifier,
		activations:      make(map[string]desktopActivation),
	}
}

type registerReq struct {
	Email       string `json:"email"`        // optional if username is set
	Username    string `json:"username"`     // optional if email is set
	Password    string `json:"password"`     // required
	DisplayName string `json:"display_name"` // required
	AccountType string `json:"account_type"` // manager | parent
}

// Register creates a new account and returns tokens.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.DisplayName == "" {
		badRequest(c, "display_name is required")
		return
	}

	if req.Email == "" && req.Username == "" {
		badRequest(c, "an email or username is required")
		return
	}

	if req.Username != "" && !usernameRe.MatchString(req.Username) {
		badRequest(
			c,
			"username must be 3-32 chars: lowercase letters, digits, underscores",
		)
		return
	}

	if len(req.Password) < 8 {
		badRequest(c, "password must be at least 8 characters")
		return
	}

	if req.AccountType == "" {
		req.AccountType = "manager"
	}

	if req.AccountType != "manager" && req.AccountType != "parent" {
		badRequest(c, "account_type must be 'manager' or 'parent'")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		serverError(c, err)
		return
	}

	u, err := h.store.CreateUser(
		c.Request.Context(),
		req.Email,
		req.Username,
		hash,
		req.DisplayName,
		req.AccountType,
	)

	if errors.Is(err, store.ErrConflict) {
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": "that email or username is already taken",
			},
		)
		return
	}

	if err != nil {
		serverError(c, err)
		return
	}

	h.issue(c, http.StatusCreated, u)
}

type loginReq struct {
	Identifier string `json:"identifier"` // email or username
	Email      string `json:"email"`      // legacy field
	Password   string `json:"password"`
	BusinessID string `json:"business_id"` // optional
}

// Login verifies credentials and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid body")
		return
	}

	identifier := req.Identifier

	if identifier == "" {
		identifier = req.Email
	}

	u, hash, err := h.store.GetUserByIdentifier(
		c.Request.Context(),
		identifier,
	)

	if err != nil {
		unauthorized(c, "invalid credentials")
		return
	}

	ok, err := auth.VerifyPassword(hash, req.Password)

	if err != nil || !ok {
		unauthorized(c, "invalid credentials")
		return
	}

	if req.BusinessID != "" {
		member, err := h.store.IsMember(
			c.Request.Context(),
			u.ID,
			req.BusinessID,
		)

		if err != nil {
			serverError(c, err)
			return
		}

		if !member {
			c.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "not a member of that business",
				},
			)
			return
		}
	}

	h.issue(c, http.StatusOK, u)
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type keycloakLoginReq struct {
	Token string `json:"token"`
}

// KeycloakLogin verifies a Keycloak ID token and creates a normal Bibo session.
func (h *AuthHandler) KeycloakLogin(c *gin.Context) {
	if h.keycloakVerifier == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"error": "keycloak authentication is not configured",
			},
		)
		return
	}

	var req keycloakLoginReq

	if err := c.ShouldBindJSON(&req); err != nil ||
		strings.TrimSpace(req.Token) == "" {

		badRequest(c, "token is required")
		return
	}

	claims, err := h.keycloakVerifier.Verify(
		c.Request.Context(),
		req.Token,
	)

	if err != nil {
		unauthorized(c, "invalid keycloak token")
		return
	}

	email := strings.ToLower(
		strings.TrimSpace(claims.Email),
	)

	if email == "" {
		unauthorized(c, "keycloak token has no email")
		return
	}

	u, _, err := h.store.GetUserByIdentifier(
		c.Request.Context(),
		email,
	)

	if err != nil {
		unauthorized(
			c,
			"no BiBo account found for this employee",
		)
		return
	}

	h.issue(c, http.StatusOK, u)
}

// ============================================================
// BIBO DESKTOP ACTIVATION
// ============================================================

type desktopActivation struct {
	UserID    string
	ExpiresAt time.Time
}

type desktopActivationReq struct {
	Token string `json:"token"`
}

// newActivationCode creates a cryptographically random
// one-time activation code.
func newActivationCode() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// CreateDesktopActivation:
//
// 1. Receives Keycloak ID token from HR application.
// 2. Verifies Keycloak token.
// 3. Finds matching Bibo employee using email.
// 4. Creates a short-lived one-time activation code.
// 5. Stores the code in backend memory.
// 6. Sends the code back to HR application.
func (h *AuthHandler) CreateDesktopActivation(c *gin.Context) {
	if h.keycloakVerifier == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"error": "keycloak authentication is not configured",
			},
		)
		return
	}

	var req desktopActivationReq

	if err := c.ShouldBindJSON(&req); err != nil ||
		strings.TrimSpace(req.Token) == "" {

		badRequest(c, "token is required")
		return
	}

	// --------------------------------------------------------
	// Verify Keycloak token
	// --------------------------------------------------------

	claims, err := h.keycloakVerifier.Verify(
		c.Request.Context(),
		req.Token,
	)

	if err != nil {
		fmt.Printf(
			"[desktop-activation] Keycloak token verification FAILED: %v\n",
			err,
		)

		unauthorized(c, "invalid keycloak token")
		return
	}

	// --------------------------------------------------------
	// Get employee email
	// --------------------------------------------------------

	email := strings.ToLower(
		strings.TrimSpace(claims.Email),
	)

	if email == "" {
		fmt.Println(
			"[desktop-activation] Keycloak token has no email",
		)

		unauthorized(
			c,
			"keycloak token has no email",
		)
		return
	}

	// --------------------------------------------------------
	// Find Bibo employee
	// --------------------------------------------------------

	u, _, err := h.store.GetUserByIdentifier(
		c.Request.Context(),
		email,
	)

	if err != nil {
		fmt.Printf(
			"[desktop-activation] Bibo user NOT FOUND email=%s error=%v\n",
			email,
			err,
		)

		unauthorized(
			c,
			"no BiBo account found for this employee",
		)
		return
	}

	// --------------------------------------------------------
	// Generate activation code
	// --------------------------------------------------------

	code, err := newActivationCode()

	if err != nil {
		serverError(c, err)
		return
	}

	// --------------------------------------------------------
	// Store activation code
	// --------------------------------------------------------

	h.activationMu.Lock()

	h.activations[code] = desktopActivation{
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(60 * time.Second),
	}

	activationCount := len(h.activations)

	h.activationMu.Unlock()

	// --------------------------------------------------------
	// DEBUG LOG
	//
	// IMPORTANT:
	// We do NOT print the actual activation code.
	// --------------------------------------------------------

	fmt.Printf(
		"[desktop-activation] CREATED code_len=%d activations=%d user=%s email=%s expires_in=60s\n",
		len(code),
		activationCount,
		u.ID,
		email,
	)

	// --------------------------------------------------------
	// Return activation code
	// --------------------------------------------------------

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":       code,
			"expires_in": 60,
		},
	)
}

// ActivateDesktop:
//
// 1. Receives one-time activation code from Bibo Desktop.
// 2. Looks up the code.
// 3. Makes sure it has not expired.
// 4. Deletes the code so it cannot be reused.
// 5. Creates Bibo access + refresh tokens.
// 6. Sends tokens to Bibo Desktop.
func (h *AuthHandler) ActivateDesktop(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil ||
		strings.TrimSpace(req.Code) == "" {

		badRequest(c, "code is required")
		return
	}

	code := strings.TrimSpace(req.Code)

	// --------------------------------------------------------
	// Look up activation code
	// --------------------------------------------------------

	h.activationMu.Lock()

	activation, ok := h.activations[code]

	activationCount := len(h.activations)

	if ok {
		// One-time activation.
		delete(h.activations, code)
	}

	h.activationMu.Unlock()

	// --------------------------------------------------------
	// DEBUG LOG
	// --------------------------------------------------------

	fmt.Printf(
		"[desktop-activate] RECEIVED code_len=%d found=%v activations=%d\n",
		len(code),
		ok,
		activationCount,
	)

	// --------------------------------------------------------
	// Code not found
	// --------------------------------------------------------

	if !ok {
		fmt.Println(
			"[desktop-activate] ERROR: activation code was NOT FOUND",
		)

		unauthorized(
			c,
			"invalid or expired activation code",
		)
		return
	}

	// --------------------------------------------------------
	// Code expired
	// --------------------------------------------------------

	if time.Now().After(activation.ExpiresAt) {
		fmt.Println(
			"[desktop-activate] ERROR: activation code EXPIRED",
		)

		unauthorized(
			c,
			"invalid or expired activation code",
		)
		return
	}

	// --------------------------------------------------------
	// Generate normal Bibo session
	// --------------------------------------------------------

	pair, err := h.tok.Issue(
		activation.UserID,
	)

	if err != nil {
		fmt.Printf(
			"[desktop-activate] TOKEN ISSUE FAILED user=%s error=%v\n",
			activation.UserID,
			err,
		)

		serverError(c, err)
		return
	}

	fmt.Printf(
		"[desktop-activate] SUCCESS user=%s\n",
		activation.UserID,
	)

	c.JSON(
		http.StatusOK,
		pair,
	)
}

// ============================================================
// REFRESH
// ============================================================

// Refresh exchanges a valid refresh token for a new token pair.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid body")
		return
	}

	userID, err := h.tok.ParseRefresh(
		req.RefreshToken,
	)

	if err != nil {
		unauthorized(c, "invalid refresh token")
		return
	}

	pair, err := h.tok.Issue(userID)

	if err != nil {
		serverError(c, err)
		return
	}

	obs.Info(
		"login ok",
		"user",
		userID,
	)

	c.JSON(
		http.StatusOK,
		pair,
	)
}

// ============================================================
// PUBLIC BUSINESSES
// ============================================================

// PublicBusinesses lists businesses + owner names
// for the login picker.
func (h *AuthHandler) PublicBusinesses(c *gin.Context) {
	list, err := h.store.ListPublicBusinesses(
		c.Request.Context(),
	)

	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"businesses": list,
		},
	)
}

// ============================================================
// ME
// ============================================================

// Me returns the authenticated user.
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := auth.UserID(c)

	if !ok {
		unauthorized(c, "unauthenticated")
		return
	}

	u, err := h.store.GetUserByID(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"id":           u.ID,
			"email":        u.Email,
			"username":     u.Username,
			"display_name": u.DisplayName,
			"account_type": u.AccountType,
		},
	)
}

// ============================================================
// ISSUE TOKENS
// ============================================================

func (h *AuthHandler) issue(
	c *gin.Context,
	status int,
	u store.User,
) {
	pair, err := h.tok.Issue(u.ID)

	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(
		status,
		gin.H{
			"user": gin.H{
				"id":           u.ID,
				"email":        u.Email,
				"username":     u.Username,
				"display_name": u.DisplayName,
				"account_type": u.AccountType,
			},
			"tokens": pair,
		},
	)
}
