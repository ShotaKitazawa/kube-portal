package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"github.com/ShotaKitazawa/kube-portal/backend/util"
)

type AuthController struct {
	logger *slog.Logger
	config *oauth2.Config

	PkceCodeVerifiers      map[uuid.UUID]string
	PkceCodeVerifiersMutex sync.Mutex
}

func NewAuthController(l *slog.Logger, config *oauth2.Config) *AuthController {
	return &AuthController{l, config, map[uuid.UUID]string{}, sync.Mutex{}}
}

func (c *AuthController) Login(ctx echo.Context) error {
	state := uuid.New()
	pkceVerifier := oauth2.GenerateVerifier()
	c.setPkceCodeVerifier(state, pkceVerifier)
	url := c.config.AuthCodeURL(state.String(),
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(pkceVerifier),
	)
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *AuthController) Logout(ctx echo.Context) error {
	// delete cookie
	ctx.SetCookie(&http.Cookie{
		Name:   util.CookieKeyIDToken,
		Path:   "/",
		Value:  "",
		MaxAge: -1,
	})

	// TODO: redirect to provider logout URL
	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (c *AuthController) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code == "" {
		return fmt.Errorf("missing query parameter: code")
	}
	state := ctx.QueryParam("state")
	if state == "" {
		return fmt.Errorf("missing query parameter: state")
	}
	id, err := uuid.Parse(state)
	if err != nil {
		return fmt.Errorf("failed to parse state as UUID")
	}
	pkceVerifier, err := c.getPkceCodeVerifier(id)
	if err != nil {
		return fmt.Errorf("failed to get PKCE verifier from state")
	}

	token, err := c.config.Exchange(ctx.Request().Context(), code, oauth2.VerifierOption(pkceVerifier))
	if err != nil {
		return err
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("id_token could not be acuired")
	}

	ctx.SetCookie(&http.Cookie{
		Name:     util.CookieKeyIDToken,
		Path:     "/",
		HttpOnly: false, // for frontend to refer to JWT's profile
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Value:    idToken,
		MaxAge:   int(token.ExpiresIn),
	})
	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (c *AuthController) setPkceCodeVerifier(key uuid.UUID, val string) {
	c.PkceCodeVerifiersMutex.Lock()
	defer c.PkceCodeVerifiersMutex.Unlock()
	c.PkceCodeVerifiers[key] = val
}

func (c *AuthController) getPkceCodeVerifier(key uuid.UUID) (string, error) {
	c.PkceCodeVerifiersMutex.Lock()
	defer c.PkceCodeVerifiersMutex.Unlock()
	res, ok := c.PkceCodeVerifiers[key]
	if !ok {
		return "", fmt.Errorf("failed to get PKCE verifier")
	}
	return res, nil
}
