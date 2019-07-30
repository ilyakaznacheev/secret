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

type testError struct {
	ErrorText string `json:"error"`
}

func (e testError) Error() string {
	return e.ErrorText
}

type testDB struct {
	secret                  *models.Secret
	err                     error
	hash                    string
	callCounterGetSecret    int
	callCounterCreateSecret int
	callCounterDeleteSecret int
	callCounterUpdateSecret int
}

func (db *testDB) GetSecret(hash string) (*models.Secret, error) {
	db.callCounterGetSecret++
	db.hash = hash
	return db.secret, db.err
}

func (db *testDB) CreateSecret(hash string, s models.Secret) error {
	db.callCounterCreateSecret++
	db.hash = hash
	return db.err
}

func (db *testDB) DeleteSecret(hash string) error {
	db.callCounterDeleteSecret++
	db.hash = hash
	return db.err
}

func (db *testDB) UpdateSecret(hash string, s models.Secret) error {
	db.callCounterUpdateSecret++
	db.hash = hash
	return db.err
}

func TestSecretHandler_GetSecret(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-01T10:10:10Z")
	future, _ := time.Parse(time.RFC3339, "2020-03-01T10:10:10Z")
	past, _ := time.Parse(time.RFC3339, "2020-01-01T10:10:10Z")
	_ = past

	encTestSecret := func(key string) string {
		res, _ := encryptSecret(key, "test_secret")
		return res
	}

	errNotFound := &testError{"test not found"}

	tests := []struct {
		name                    string
		secretID                string
		respCode                int
		respBody                string
		headers                 map[string]string
		db                      *testDB
		callCounterGetSecret    int
		callCounterDeleteSecret int
		callCounterUpdateSecret int
	}{
		{
			name:     "simple",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 200,
			respBody: `{"createdAt":"2020-02-01T10:10:10.000Z","expiresAt":"2020-03-01T10:10:10.000Z","hash":"5621caf61d79545957a49c7d","remainingViews":99,"secretText":"test_secret"}`,
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 1,
		},

		{
			name:     "no data",
			secretID: "12345",
			respCode: 404,
			respBody: `{"error":"test not found"}`,
			db: &testDB{
				secret: nil,
				err:    errNotFound,
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 0,
		},

		{
			name:     "expiration error",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 404,
			respBody: `{"error":"secret isn't valid anymore"}`,
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      past,
						ExpiresAt:      &past,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 1,
			callCounterUpdateSecret: 0,
		},

		{
			name:     "counter error",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 404,
			respBody: `{"error":"secret isn't valid anymore"}`,
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 0,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 1,
			callCounterUpdateSecret: 0,
		},

		{
			name:     "accept json",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 200,
			respBody: `{"createdAt":"2020-02-01T10:10:10.000Z","expiresAt":"2020-03-01T10:10:10.000Z","hash":"5621caf61d79545957a49c7d","remainingViews":99,"secretText":"test_secret"}`,
			headers: map[string]string{
				"Accept": "application/json",
			},
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 1,
		},

		{
			name:     "accept xml",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 200,
			respBody: `<Secret><createdAt>2020-02-01T10:10:10.000Z</createdAt><expiresAt>2020-03-01T10:10:10.000Z</expiresAt><hash>5621caf61d79545957a49c7d</hash><remainingViews>99</remainingViews><secretText>test_secret</secretText></Secret>`,
			headers: map[string]string{
				"Accept": "application/xml",
			},
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 1,
		},

		{
			name:     "accept unknown -> json",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 200,
			respBody: `{"createdAt":"2020-02-01T10:10:10.000Z","expiresAt":"2020-03-01T10:10:10.000Z","hash":"5621caf61d79545957a49c7d","remainingViews":99,"secretText":"test_secret"}`,
			headers: map[string]string{
				"Accept": "application/noname",
			},
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash: "5621caf61d79545957a49c7d",
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			h := SecretHandler{
				db:      tt.db,
				nowFunc: func() time.Time { return now },
			}
			router.GET("/secret/:hash", h.GetSecret)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/secret/"+tt.secretID, nil)
			for key, value := range tt.headers {
				req.Header.Add(key, value)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.respCode, w.Code)
			assert.Equal(t, tt.respBody, w.Body.String())

			assert.Equal(t, tt.secretID, tt.db.hash)
			assert.Equal(t, tt.callCounterGetSecret, tt.db.callCounterGetSecret)
			assert.Equal(t, tt.callCounterDeleteSecret, tt.db.callCounterDeleteSecret)
			assert.Equal(t, tt.callCounterUpdateSecret, tt.db.callCounterUpdateSecret)
		})
	}
}
