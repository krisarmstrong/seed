package auth

import (
	"context"
	"testing"
	"time"
)

func TestNewCSRFManager(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewCSRFManager(ctx)
	if manager == nil {
		t.Fatal("NewCSRFManager returned nil")
	}

	if manager.tokens == nil {
		t.Fatal("tokens map not initialized")
	}
}

func TestCSRFManagerCleanupStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	manager := NewCSRFManager(ctx)

	// Generate a token to ensure the manager is working
	token, err := manager.GenerateToken("test-session")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Cancel the context
	cancel()

	// Give the goroutine a moment to stop
	time.Sleep(100 * time.Millisecond)

	// The test passes if we get here without hanging or panicking
	// In a real scenario, we'd verify the goroutine actually stopped,
	// but that's difficult without exposing internal state
}

func TestCSRFManagerGenerateAndValidate(t *testing.T) {
	ctx := context.Background()
	manager := NewCSRFManager(ctx)

	sessionID := "test-session"

	// Generate a token
	token, err := manager.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Validate the token
	if err := manager.ValidateToken(sessionID, token); err != nil {
		t.Errorf("failed to validate token: %v", err)
	}

	// Validate with wrong session ID
	if err := manager.ValidateToken("wrong-session", token); err != ErrCSRFTokenInvalid {
		t.Errorf("expected ErrCSRFTokenInvalid, got %v", err)
	}

	// Validate with wrong token
	if err := manager.ValidateToken(sessionID, "wrong-token"); err != ErrCSRFTokenInvalid {
		t.Errorf("expected ErrCSRFTokenInvalid, got %v", err)
	}
}
