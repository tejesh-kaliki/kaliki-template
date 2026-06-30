package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/example/jwt-basic-app/backend/internal/testsupport"
)

func decode(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v (%s)", err, string(b))
	}
	return m
}

func TestSignup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTest(t)

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
			`{"email":"a@b.com","password":"password123","name":"A"}`)
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201 (%s)", w.Code, w.Body.String())
		}
		body := decode(t, w.Body.Bytes())
		if body["token"] == "" || body["token"] == nil {
			t.Fatal("expected a session token")
		}

		// Password must be stored hashed, never in plaintext.
		var hash string
		if err := testDB.Pool.QueryRow(context.Background(),
			`SELECT password_hash FROM users WHERE email = $1`, "a@b.com").Scan(&hash); err != nil {
			t.Fatalf("query: %v", err)
		}
		if hash == "" || hash == "password123" {
			t.Fatalf("password not hashed: %q", hash)
		}
		if body["refresh_token"] == "" || body["refresh_token"] == nil {
			t.Fatal("expected a refresh token")
		}
	})

	t.Run("conflict", func(t *testing.T) {
		setupTest(t)

		testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
			`{"email":"dup@b.com","password":"password123"}`)
		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
			`{"email":"dup@b.com","password":"password123"}`)
		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", w.Code)
		}
	})

	t.Run("validation", func(t *testing.T) {
		setupTest(t)

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
			`{"email":"a@b.com","password":"short"}`)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	signup := func() {
		testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
			`{"email":"u@b.com","password":"password123"}`)
	}

	t.Run("success", func(t *testing.T) {
		setupTest(t)
		signup()

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/login",
			`{"email":"u@b.com","password":"password123"}`)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (%s)", w.Code, w.Body.String())
		}
		if decode(t, w.Body.Bytes())["token"] == nil {
			t.Fatal("expected a token")
		}
	})

	t.Run("authorization", func(t *testing.T) {
		setupTest(t)
		signup()

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/login",
			`{"email":"u@b.com","password":"wrongpass"}`)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", w.Code)
		}
	})

	t.Run("validation", func(t *testing.T) {
		setupTest(t)

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/login", `{}`)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", w.Code)
		}
	})
}

// signupSession signs up and returns the (access, refresh) pair.
func signupSession(t *testing.T, email string) (string, string) {
	t.Helper()
	w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/signup",
		`{"email":"`+email+`","password":"password123"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("signup status = %d (%s)", w.Code, w.Body.String())
	}
	b := decode(t, w.Body.Bytes())
	return b["token"].(string), b["refresh_token"].(string)
}

func TestGetCurrentUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTest(t)
		access, _ := signupSession(t, "me@b.com")

		w := testsupport.DoJSONAuth(router, http.MethodGet, "/api/v1/auth/me", "", access)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (%s)", w.Code, w.Body.String())
		}
		if decode(t, w.Body.Bytes())["email"] != "me@b.com" {
			t.Fatalf("unexpected user: %s", w.Body.String())
		}
	})

	t.Run("authorization", func(t *testing.T) {
		setupTest(t)

		w := testsupport.DoJSON(router, http.MethodGet, "/api/v1/auth/me", "")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", w.Code)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTest(t)
		_, refresh := signupSession(t, "rt@b.com")

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/refresh",
			`{"refresh_token":"`+refresh+`"}`)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (%s)", w.Code, w.Body.String())
		}
		if decode(t, w.Body.Bytes())["token"] == nil {
			t.Fatal("expected a new access token")
		}

		// Rotation: the old refresh token is revoked and cannot be reused.
		reuse := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/refresh",
			`{"refresh_token":"`+refresh+`"}`)
		if reuse.Code != http.StatusUnauthorized {
			t.Fatalf("reused refresh status = %d, want 401", reuse.Code)
		}
	})

	t.Run("authorization", func(t *testing.T) {
		setupTest(t)

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/refresh",
			`{"refresh_token":"nope"}`)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", w.Code)
		}
	})
}

func TestLogout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTest(t)
		_, refresh := signupSession(t, "lo@b.com")

		w := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/logout",
			`{"refresh_token":"`+refresh+`"}`)
		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204 (%s)", w.Code, w.Body.String())
		}

		// The refresh token is revoked after logout.
		after := testsupport.DoJSON(router, http.MethodPost, "/api/v1/auth/refresh",
			`{"refresh_token":"`+refresh+`"}`)
		if after.Code != http.StatusUnauthorized {
			t.Fatalf("refresh after logout = %d, want 401", after.Code)
		}
	})
}
