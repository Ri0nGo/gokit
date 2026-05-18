package oauth2Client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientAuthCodeURL(t *testing.T) {
	client, err := New(Config{BaseURL: " https://iam.example.com/ ", ClientID: " app ", ClientSecret: " secret ", RedirectURI: " https://app.example.com/callback ", Scopes: []string{"profile"}})
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := client.AuthCodeURL("state-1")
	if err != nil {
		t.Fatal(err)
	}

	want := "https://iam.example.com/api/v1/oauth/authorize?client_id=app&redirect_uri=https%3A%2F%2Fapp.example.com%2Fcallback&response_type=code&scope=profile&state=state-1"
	if authURL != want {
		t.Fatalf("auth url = %q, want %q", authURL, want)
	}
}

func TestClientGetUserTokenAndInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case tokenPath:
			query := r.URL.Query()
			if r.Method != http.MethodGet || query.Get("grant_type") != "authorization_code" || query.Get("client_id") != "app" || query.Get("secret") != "secret" || query.Get("code") != "code-1" {
				t.Fatalf("unexpected token request: method=%s query=%s", r.Method, r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(apiResponse{Code: 0, Msg: "ok", Data: mustRawMessage(t, TokenResponse{AccessToken: "token-1", ExpiresIn: 3600, RefreshToken: "refresh-1", OpenID: "openid-1", Scope: "profile"})})
		case refreshPath:
			query := r.URL.Query()
			if r.Method != http.MethodGet || query.Get("grant_type") != "refresh_token" || query.Get("client_id") != "app" || query.Get("refresh_token") != "refresh-1" {
				t.Fatalf("unexpected refresh request: method=%s query=%s", r.Method, r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(apiResponse{Code: 0, Msg: "ok", Data: mustRawMessage(t, TokenResponse{AccessToken: "token-2", ExpiresIn: 3600, RefreshToken: "refresh-1", OpenID: "openid-1", Scope: "profile"})})
		case authPath:
			query := r.URL.Query()
			if query.Get("access_token") != "token-1" || query.Get("openid") != "openid-1" {
				t.Fatalf("unexpected auth query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(apiResponse{Code: 0, Msg: "ok", Data: json.RawMessage(`{"success":true}`)})
		case userinfoPath:
			query := r.URL.Query()
			if query.Get("access_token") != "token-1" || query.Get("openid") != "openid-1" {
				t.Fatalf("unexpected userinfo query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(apiResponse{Code: 0, Msg: "ok", Data: json.RawMessage(`{"openid":"openid-1","username":"alice","display_name":"Alice","status":1,"roles":["admin"]}`)})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := New(Config{BaseURL: server.URL, ClientID: "app", ClientSecret: "secret", RedirectURI: "https://app.example.com/callback"})
	if err != nil {
		t.Fatal(err)
	}

	token, err := client.GetUserToken(context.Background(), "code-1")
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "token-1" {
		t.Fatalf("unexpected token: %#v", token)
	}
	if token.RefreshToken != "refresh-1" || token.OpenID != "openid-1" {
		t.Fatalf("unexpected token identity: %#v", token)
	}

	if err := client.CheckUserToken(context.Background(), token.AccessToken, token.OpenID); err != nil {
		t.Fatal(err)
	}

	refreshed, err := client.RefreshUserToken(context.Background(), token.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.AccessToken != "token-2" || refreshed.OpenID != "openid-1" {
		t.Fatalf("unexpected refreshed token: %#v", refreshed)
	}

	user, err := client.GetUserInfo(context.Background(), token.AccessToken, token.OpenID)
	if err != nil {
		t.Fatal(err)
	}
	if user.OpenID != "openid-1" || user.Username != "alice" || len(user.Roles) != 1 || user.Roles[0] != "admin" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestClientErrorsOnMissingData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apiResponse{Code: 0, Msg: "ok"})
	}))
	defer server.Close()

	client, err := New(Config{BaseURL: server.URL, ClientID: "app", ClientSecret: "secret", RedirectURI: "https://app.example.com/callback"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetUserToken(context.Background(), "code-1"); err == nil {
		t.Fatal("expected missing data error")
	}
}

func mustRawMessage(t *testing.T, value any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
