package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	hash                    string
	newSecret               models.Secret
	errGetSecret            error
	errCreateSecret         error
	errDeleteSecret         error
	errUpdateSecret         error
	callCounterGetSecret    int
	callCounterCreateSecret int
	callCounterDeleteSecret int
	callCounterUpdateSecret int
}

func (db *testDB) GetSecret(hash string) (*models.Secret, error) {
	db.callCounterGetSecret++
	db.hash = hash
	return db.secret, db.errGetSecret
}

func (db *testDB) CreateSecret(hash string, s models.Secret) error {
	db.callCounterCreateSecret++
	db.hash = hash
	db.newSecret = s
	return db.errCreateSecret
}

func (db *testDB) DeleteSecret(hash string) error {
	db.callCounterDeleteSecret++
	db.hash = hash
	return db.errDeleteSecret
}

func (db *testDB) UpdateSecret(hash string, s models.Secret) error {
	db.callCounterUpdateSecret++
	db.hash = hash
	return db.errUpdateSecret
}

func TestSecretHandler_GetSecret(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-01T10:10:10Z")
	future, _ := time.Parse(time.RFC3339, "2020-03-01T10:10:10Z")
	past, _ := time.Parse(time.RFC3339, "2020-01-01T10:10:10Z")

	encTestSecret := func(key string) string {
		res, _ := encryptSecret(key, "test_secret")
		return res
	}

	errTest := &testError{"test error"}

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
			respBody: `{"error":"test error"}`,
			db: &testDB{
				secret:       nil,
				errGetSecret: errTest,
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 0,
		},

		{
			name:     "update error",
			secretID: "5621caf61d79545957a49c7d",
			respCode: 404,
			respBody: `{"error":"test error"}`,
			db: &testDB{
				secret: &models.Secret{
					SecretBase: models.SecretBase{
						CreatedAt:      now,
						ExpiresAt:      &future,
						RemainingViews: 100,
						SecretText:     encTestSecret("5621caf61d79545957a49c7d"),
					},
				},
				hash:            "5621caf61d79545957a49c7d",
				errUpdateSecret: errTest,
			},
			callCounterGetSecret:    1,
			callCounterDeleteSecret: 0,
			callCounterUpdateSecret: 1,
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

func TestSecretHandler_PostSecret(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-01T10:10:10Z")
	future, _ := time.Parse(time.RFC3339, "2020-02-01T10:20:10Z")

	encTestSecret := func(key string) string {
		res, _ := encryptSecret(key, "test_secret")
		return res
	}

	errTest := &testError{"test error"}

	testKey := "5621caf61d79545957a49c7d"

	tests := []struct {
		name                    string
		respCode                int
		respBody                string
		headers                 map[string]string
		postFields              map[string]string
		db                      *testDB
		secret                  models.Secret
		callCounterCreateSecret int
	}{
		{
			name:     "normal creation",
			respCode: 200,
			respBody: `{"createdAt":"2020-02-01T10:10:10.000Z","expiresAt":"2020-02-01T10:20:10.000Z","hash":"5621caf61d79545957a49c7d","remainingViews":10,"secretText":"test_secret"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "10",
				"expireAfter":      "10",
			},
			db: &testDB{},
			secret: models.Secret{
				SecretBase: models.SecretBase{
					CreatedAt:      now,
					ExpiresAt:      &future,
					RemainingViews: 10,
					SecretText:     encTestSecret(testKey),
				},
			},
			callCounterCreateSecret: 1,
		},

		{
			name:     "bad view counter",
			respCode: 405,
			respBody: `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "abc",
				"expireAfter":      "10",
			},
			db:                      &testDB{},
			callCounterCreateSecret: 0,
		},

		{
			name:     "bad expiration time",
			respCode: 405,
			respBody: `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "10",
				"expireAfter":      "abc",
			},
			db:                      &testDB{},
			callCounterCreateSecret: 0,
		},

		{
			name:     "zero view counter",
			respCode: 405,
			respBody: `{"error":"wrong expireAfterViews value 0"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "0",
				"expireAfter":      "10",
			},
			db:                      &testDB{},
			callCounterCreateSecret: 0,
		},

		{
			name:     "infinite expiration time",
			respCode: 200,
			respBody: `{"createdAt":"2020-02-01T10:10:10.000Z","hash":"5621caf61d79545957a49c7d","remainingViews":10,"secretText":"test_secret"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "10",
				"expireAfter":      "0",
			},
			db: &testDB{},
			secret: models.Secret{
				SecretBase: models.SecretBase{
					CreatedAt:      now,
					ExpiresAt:      nil,
					RemainingViews: 10,
					SecretText:     encTestSecret(testKey),
				},
			},
			callCounterCreateSecret: 1,
		},

		{
			name:     "creation error",
			respCode: 405,
			respBody: `{"error":"test error"}`,
			postFields: map[string]string{
				"secret":           "test_secret",
				"expireAfterViews": "10",
				"expireAfter":      "10",
			},
			db: &testDB{
				errCreateSecret: errTest,
			},
			secret: models.Secret{
				SecretBase: models.SecretBase{
					CreatedAt:      now,
					ExpiresAt:      &future,
					RemainingViews: 10,
					SecretText:     encTestSecret(testKey),
				},
			},
			callCounterCreateSecret: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			h := SecretHandler{
				db:      tt.db,
				nowFunc: func() time.Time { return now },
				keygen:  func() string { return testKey },
			}
			router.POST("/secret", h.PostSecret)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/secret", nil)
			for key, value := range tt.headers {
				req.Header.Add(key, value)
			}
			req.PostForm = make(url.Values)
			for key, value := range tt.postFields {
				req.PostForm.Add(key, value)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.respCode, w.Code)
			assert.Equal(t, tt.respBody, w.Body.String())

			assert.Equal(t, tt.secret.CreatedAt, tt.db.newSecret.CreatedAt)
			assert.Equal(t, tt.secret.ExpiresAt, tt.db.newSecret.ExpiresAt)
			assert.Equal(t, tt.secret.RemainingViews, tt.db.newSecret.RemainingViews)

			assert.Equal(t, tt.callCounterCreateSecret, tt.db.callCounterCreateSecret)
		})
	}
}
