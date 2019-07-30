package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/ilyakaznacheev/secret/internal/models"
)

var (
	// ErrSecretOutdated secret in not valid anymore
	ErrSecretOutdated = errors.New("secret isn't valid anymore")
)

// SecretHandler is a REST API handler for secrete service
type SecretHandler struct {
	db      Database
	nowFunc func() time.Time
}

// NewSecretHandler creates a new API handler
func NewSecretHandler(db Database) *SecretHandler {
	return &SecretHandler{
		db:      db,
		nowFunc: time.Now,
	}
}

// GetSecret returns a secret if possible
func (h *SecretHandler) GetSecret(c *gin.Context) {
	now := h.nowFunc()

	hash := c.Param("hash")
	s, err := h.db.GetSecret(hash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// validity checks
	if s.ExpiresAt != nil && now.After(*s.ExpiresAt) {
		log.Printf("secret %s expire time is out", hash)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": ErrSecretOutdated.Error()})
		defer h.db.DeleteSecret(hash)
		return
	}

	if s.RemainingViews > 0 {
		s.RemainingViews--
		if err := h.db.UpdateSecret(hash, *s); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	} else {
		log.Printf("secret %s expire counter is out", hash)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": ErrSecretOutdated.Error()})
		defer h.db.DeleteSecret(hash)
		return
	}

	// decrypt secret
	encSecret, err := decryptSecret(hash, s.SecretText)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var expFormatted *strfmt.DateTime
	if s.ExpiresAt != nil {
		exp := strfmt.DateTime(*s.ExpiresAt)
		expFormatted = &exp
	}

	// prepare response structure
	res := models.SecretResponse{
		CreatedAt:      strfmt.DateTime(s.CreatedAt),
		ExpiresAt:      expFormatted,
		Hash:           hash,
		RemainingViews: s.RemainingViews,
		SecretText:     encSecret,
	}

	getResponseFunc(c)(&res)
}

// PostSecret creates a new secret
func (h *SecretHandler) PostSecret(c *gin.Context) {
	// read and parse parameters
	secret := c.PostForm("secret")
	expireCounterStr := c.PostForm("expireAfterViews")
	expireTimeoutStr := c.PostForm("expireAfter")

	expireCounter, err := strconv.Atoi(expireCounterStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
		return
	}

	expireTimeout, err := strconv.Atoi(expireTimeoutStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
		return
	}

	// validity checks
	if expireCounter <= 0 {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": fmt.Errorf("wrong expireAfterViews value %d", expireCounter)})
		return
	}

	var expTime *time.Time
	if expireTimeout > 0 {
		exp := h.nowFunc().Add(time.Minute * time.Duration(expireTimeout))
		expTime = &exp
	}

	// encrypt secret
	key := generateKey()
	encSecret, err := encryptSecret(key, secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
		return
	}

	// fill db model data
	s := models.Secret{
		SecretBase: models.SecretBase{
			CreatedAt:      h.nowFunc(),
			ExpiresAt:      expTime,
			RemainingViews: int32(expireCounter),
			SecretText:     encSecret,
		},
	}

	// save to database
	if err := h.db.CreateSecret(key, s); err != nil {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
		return
	}

	// if nil set empty time as default
	var expFormatted *strfmt.DateTime
	if expTime != nil {
		exp := strfmt.DateTime(*expTime)
		expFormatted = &exp
	}

	// prepare response structure
	res := models.SecretResponse{
		CreatedAt:      strfmt.DateTime(s.CreatedAt),
		ExpiresAt:      expFormatted,
		Hash:           key,
		RemainingViews: s.RemainingViews,
		SecretText:     secret,
	}

	log.Printf("key %s was issued for IP %s", key, c.Request.Host)

	getResponseFunc(c)(&res)
}

// RedirectTo redirects to provided url
func RedirectTo(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, url)
		c.Abort()
	}
}

// getResponseFunc returns data marshalling function for accepted MIME type
func getResponseFunc(c *gin.Context) func(interface{}) {
	mimeTypes := strings.Split(strings.Replace(c.GetHeader("Accept"), " ", "", -1), ",")
	for _, mime := range mimeTypes {
		switch mime {
		case "application/json":
			return func(v interface{}) {
				c.JSON(http.StatusOK, v)
			}
		case "application/xml":
			return func(v interface{}) {
				c.XML(http.StatusOK, v)
			}
		}
	}
	// default json
	return func(v interface{}) {
		c.JSON(http.StatusOK, v)
	}
}

// Database is a database layer interface
type Database interface {
	GetSecret(hash string) (*models.Secret, error)
	CreateSecret(hash string, s models.Secret) error
	DeleteSecret(hash string) error
	UpdateSecret(hash string, s models.Secret) error
}
