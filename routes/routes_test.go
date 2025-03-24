package routes

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilhamtubagus/go-shorten-url/entity"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockShortenedService is a mock of the ShortenedService interface
type MockShortenedService struct {
	mock.Mock
}

func (m *MockShortenedService) ShortenURL(ctx context.Context, originalURL string) (*entity.ShortenedURL, error) {
	args := m.Called(ctx, originalURL)
	return args.Get(0).(*entity.ShortenedURL), args.Error(1)
}

func (m *MockShortenedService) GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error) {
	args := m.Called(ctx, shortcode)
	return args.Get(0).(*entity.ShortenedURL), args.Error(1)
}

func (m *MockShortenedService) ListShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error) {
	args := m.Called(ctx)
	return args.Get(0).(*[]entity.ShortenedURL), args.Error(1)
}

func (m *MockShortenedService) DeleteShortenedURL(ctx context.Context, shortcode string) error {
	args := m.Called(ctx, shortcode)
	return args.Error(0)
}

func (m *MockShortenedService) UpdateShortenedURL(ctx context.Context, shortcode string, originalURL string) (*entity.ShortenedURL, error) {
	args := m.Called(ctx, shortcode, originalURL)
	return args.Get(0).(*entity.ShortenedURL), args.Error(1)
}

func TestRoutes_Index(t *testing.T) {
	tmpl := template.Must(template.New("index").Parse("Index Page"))
	mockService := new(MockShortenedService)
	routes := NewRoutes(tmpl, mockService)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/", routes.Index())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Index Page", rr.Body.String())
}

func TestRoutes_NotFound(t *testing.T) {
	tmpl := template.Must(template.New("404.html").Parse("404 Not Found"))
	mockService := new(MockShortenedService)
	routes := NewRoutes(tmpl, mockService)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	rr := httptest.NewRecorder()

	handler := routes.NotFound()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "404 Not Found", rr.Body.String())
}

func TestRoutes_ShortenURL(t *testing.T) {
	tmpl := template.Must(template.New("shorten.html").Parse("Shortened: {{.ShortenedURL}}"))
	mockService := new(MockShortenedService)
	routes := NewRoutes(tmpl, mockService)

	mockService.On("ShortenURL", mock.Anything, "https://example.com").Return(&entity.ShortenedURL{
		OriginalURL:  "https://example.com",
		ShortCode:    "abc123",
		ShortenedURL: "http://short.url/abc123",
	}, nil)

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString("originalURL=https://example.com"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.POST("/shorten", routes.ShortenURL())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Shortened: http://short.url/abc123", rr.Body.String())
}

func TestRoutes_RedirectURL(t *testing.T) {
	tmpl := template.Must(template.New("404.html").Parse("404 Not Found"))
	mockService := new(MockShortenedService)
	routes := NewRoutes(tmpl, mockService)

	mockService.On("GetByShortCode", mock.Anything, "abc123").Return(&entity.ShortenedURL{
		OriginalURL:  "https://example.com",
		ShortCode:    "abc123",
		ShortenedURL: "http://short.url/abc123",
	}, nil)

	req, _ := http.NewRequest("GET", "/abc123", nil)
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/:shortCode", routes.RedirectURL())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "https://example.com", rr.Header().Get("Location"))
}

func TestRoutes_ListShortenedURLs(t *testing.T) {
	tmpl := template.Must(template.New("list.html").Parse("{{range .}}{{.ShortenedURL}}\n{{end}}"))
	mockService := new(MockShortenedService)
	routes := NewRoutes(tmpl, mockService)

	mockURLs := &[]entity.ShortenedURL{
		{OriginalURL: "https://example1.com", ShortCode: "abc123", ShortenedURL: "http://short.url/abc123"},
		{OriginalURL: "https://example2.com", ShortCode: "def456", ShortenedURL: "http://short.url/def456"},
	}
	mockService.On("ListShortenedURLs", mock.Anything).Return(mockURLs, nil)

	req, _ := http.NewRequest("GET", "/list", nil)
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/list", routes.ListShortenedURLs())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "http://short.url/abc123\nhttp://short.url/def456\n", rr.Body.String())
}

func TestRoutes_DeleteShortenedURL(t *testing.T) {
	mockService := new(MockShortenedService)
	routes := NewRoutes(nil, mockService)

	mockService.On("DeleteShortenedURL", mock.Anything, "abc123").Return(nil)

	req, _ := http.NewRequest("DELETE", "/abc123", nil)
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.DELETE("/:shortCode", routes.DeleteShortenedURL())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRoutes_UpdateShortenedURL(t *testing.T) {
	mockService := new(MockShortenedService)
	routes := NewRoutes(nil, mockService)

	mockService.On("UpdateShortenedURL", mock.Anything, "abc123", "https://newexample.com").Return(&entity.ShortenedURL{
		OriginalURL:  "https://newexample.com",
		ShortCode:    "abc123",
		ShortenedURL: "http://short.url/abc123",
	}, nil)

	req, _ := http.NewRequest("PUT", "/abc123", bytes.NewBufferString("newOriginalURL=https://newexample.com"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.PUT("/:shortCode", routes.UpdateShortenedURL())
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
