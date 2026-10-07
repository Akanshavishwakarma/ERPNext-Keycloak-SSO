package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeycloakVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewKeycloakVerifier(ctx context.Context, issuer, clientID string) (*KeycloakVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("keycloak provider: %w", err)
	}

	config := &oidc.Config{
		ClientID: clientID,
	}

	return &KeycloakVerifier{
		verifier: provider.Verifier(config),
	}, nil
}

type KeycloakClaims struct {
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
}

func (v *KeycloakVerifier) Verify(ctx context.Context, rawToken string) (*KeycloakClaims, error) {
	token, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("invalid keycloak token: %w", err)
	}

	var claims KeycloakClaims
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("invalid keycloak claims: %w", err)
	}

	return &claims, nil
}
