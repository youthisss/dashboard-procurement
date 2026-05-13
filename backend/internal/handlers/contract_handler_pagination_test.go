package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func paginationTestContext(rawQuery string) *gin.Context {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/contracts/dedicated-fix?"+rawQuery, nil)
	c.Request = req
	return c
}

func TestParsePaginationUsesDefaultLimitWhenEndMissing(t *testing.T) {
	limit, offset := parsePagination(paginationTestContext("_start=20"))
	if limit != contractDefaultPageLimit {
		t.Fatalf("limit = %d, want %d", limit, contractDefaultPageLimit)
	}
	if offset != 20 {
		t.Fatalf("offset = %d, want 20", offset)
	}
}

func TestParsePaginationCapsLargeRange(t *testing.T) {
	limit, offset := parsePagination(paginationTestContext("_start=0&_end=999999"))
	if limit != contractMaxPageLimit {
		t.Fatalf("limit = %d, want %d", limit, contractMaxPageLimit)
	}
	if offset != 0 {
		t.Fatalf("offset = %d, want 0", offset)
	}
}

func TestParsePaginationFallsBackWhenEndInvalid(t *testing.T) {
	limit, offset := parsePagination(paginationTestContext("_start=10&_end=abc"))
	if limit != contractDefaultPageLimit {
		t.Fatalf("limit = %d, want %d", limit, contractDefaultPageLimit)
	}
	if offset != 10 {
		t.Fatalf("offset = %d, want 10", offset)
	}
}
