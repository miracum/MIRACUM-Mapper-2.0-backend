package server

import (
	"log"
	"miracummapper/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestStuff(database *gorm.DB) (router *gin.Engine, w *httptest.ResponseRecorder) {
	router = gin.Default()
	svr := CreateServer(database, nil, nil)
	strictHandler := api.NewStrictHandler(svr, nil)
	api.RegisterHandlers(router, strictHandler)

	w = httptest.NewRecorder()

	return router, w
}

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestPingRoute(t *testing.T) {
	// gormDB, mock := NewMockDB()
	gormDB, _ := NewMockDB()
	router, w := setupTestStuff(gormDB)

	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}\n", w.Body.String())
}
