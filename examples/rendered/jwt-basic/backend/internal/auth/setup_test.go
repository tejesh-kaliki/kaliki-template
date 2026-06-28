package auth_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/example/jwt-basic-app/backend/internal/auth"
	"github.com/example/jwt-basic-app/backend/internal/config"
	"github.com/example/jwt-basic-app/backend/internal/testsupport"
)

var (
	testDB *testsupport.TestDB
	router *gin.Engine
)

func TestMain(m *testing.M) {
	testDB = testsupport.Connect("test_auth")

	r, api := testsupport.NewRouter()
	auth.New(testDB.DB, config.TokenConfig{Secret: "test-secret", ExpiryHours: 1}).Register(api)
	router = r

	os.Exit(m.Run())
}

func setupTest(t *testing.T) {
	t.Helper()
	testDB.Truncate(t)
}
