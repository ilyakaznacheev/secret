package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/secret/internal/models"
	"github.com/stretchr/testify/assert"
)

type testDB struct {
	secret *models.Secret
	err    error
	hash   string
}

func (db *testDB) GetSecret(hash string) (*models.Secret, error) {
	db.hash = hash
	return db.secret, db.err
}

func (db *testDB) CreateSecret(hash string, s models.Secret) error {
	db.hash = hash
	return db.err
}

func (db *testDB) DeleteSecret(hash string) error {
	db.hash = hash
	return db.err
}

func (db *testDB) UpdateSecret(hash string, s models.Secret) error {
	db.hash = hash
	return db.err
}

func TestSecretHandler_GetSecret(t *testing.T) {
	now := time.Now()
	future := time.Now().Add(time.Hour)

	tests := []struct {
		name     string
		secretID string
		respCode int
		respBody string
		db       testDB
	}{
		{
			name:     "simple",
			secretID: "12345",
			respCode: 200,
			respBody: "",
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     "test",
					},
					Version: 0,
				},
				err:  nil,
				hash: "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.Default()
			h := SecretHandler{
				db:      tt.db,
				nowFunc: func() time.Time { return now },
			}
			router.GET("/secret/:hash", h.GetSecret)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/secret/"+tt.secretID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.respCode, w.Code)
			assert.Equal(t, tt.respBody, w.Body.String())
			assert.Equal(t, tt.secret, tt.db.hash)
			// assert.JSONEq(t, expected, actual, msgAndArgs)

		})
	}
}
