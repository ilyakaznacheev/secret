package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecretHandler_GetSecret(t *testing.T) {
	type fields struct {
		db Database
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &SecretHandler{
				db: tt.fields.db,
			}
			h.GetSecret(tt.args.c)
		})
	}
}
