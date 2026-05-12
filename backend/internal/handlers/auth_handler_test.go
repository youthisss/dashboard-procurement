package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rygell-dashboard/internal/config"
	"rygell-dashboard/internal/models"

	"github.com/gin-gonic/gin"
)

func TestRequireAdminUsesConfiguredAdminUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		cfg        *config.Config
		username   any
		setUser    bool
		wantOK     bool
		wantStatus int
	}{
		{
			name:       "allows configured admin username",
			cfg:        &config.Config{AdminUsername: "owner"},
			username:   "owner",
			setUser:    true,
			wantOK:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "rejects hard coded admin when config uses another username",
			cfg:        &config.Config{AdminUsername: "owner"},
			username:   "admin",
			setUser:    true,
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "falls back to admin when config is nil",
			cfg:        nil,
			username:   "admin",
			setUser:    true,
			wantOK:     true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "requires username claim",
			cfg:        &config.Config{AdminUsername: "admin"},
			setUser:    false,
			wantOK:     false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects malformed username claim",
			cfg:        &config.Config{AdminUsername: "admin"},
			username:   123,
			setUser:    true,
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			if tt.setUser {
				ctx.Set("username", tt.username)
			}

			handler := NewAuthHandler(nil, tt.cfg)
			if got := handler.requireAdmin(ctx); got != tt.wantOK {
				t.Fatalf("requireAdmin() = %v, want %v", got, tt.wantOK)
			}
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserPayloadIncludesConfiguredAdminFlag(t *testing.T) {
	handler := NewAuthHandler(nil, &config.Config{AdminUsername: "owner"})

	adminPayload := handler.userPayload(&models.User{ID: 1, Name: "Owner", Username: "owner"})
	if adminPayload["is_admin"] != true {
		body, _ := json.Marshal(adminPayload)
		t.Fatalf("admin payload = %s, want is_admin=true", body)
	}

	staffPayload := handler.userPayload(&models.User{ID: 2, Name: "Staff", Username: "admin"})
	if staffPayload["is_admin"] != false {
		body, _ := json.Marshal(staffPayload)
		t.Fatalf("staff payload = %s, want is_admin=false", body)
	}
}
