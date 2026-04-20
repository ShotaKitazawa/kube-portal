package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShotaKitazawa/kube-portal/internal/api"
	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type SecurityHandler struct {
	oidcVerifier *oidc.IDTokenVerifier
	allowedSubs  []string
}

var _ api.SecurityHandler = (*SecurityHandler)(nil)

func NewSecurityHandler(v *oidc.IDTokenVerifier, allowedSubs []string) *SecurityHandler {
	return &SecurityHandler{oidcVerifier: v, allowedSubs: allowedSubs}
}

func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, _ api.OperationName, t api.BearerAuth) (context.Context, error) {
	if s.oidcVerifier == nil {
		// --disable-oidc mode: skip verification
		return ctx, nil
	}
	idToken, err := s.oidcVerifier.Verify(ctx, t.Token)
	if err != nil {
		return ctx, fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}

	if len(s.allowedSubs) > 0 {
		allowed := false
		for _, sub := range s.allowedSubs {
			if sub == idToken.Subject {
				allowed = true
				break
			}
		}
		if !allowed {
			return ctx, fmt.Errorf("%w: subject %q is not in allowed list", ErrForbidden, idToken.Subject)
		}
	}

	ctx = setToken(ctx, t.Token)
	return setIsLogin(ctx, true), nil
}
