package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/ilyakaznacheev/secret/internal/models"
)

const (
	hashSecretKey  = "secret"
	hashVersionKey = "secret_version"
)

// RedisDB is a database interaction manager for Redis
type RedisDB struct {
	client *redis.Client
}

// NewRedisDB creates a new database connection to Redis
func NewRedisDB(address string) (*RedisDB, error) {
	opts, err := redis.ParseURL(address)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	pong, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	log.Println(pong)
	return &RedisDB{
		client: client,
	}, nil
}

// GetSecret returns a secret or errors
func (r *RedisDB) GetSecret(hash string) (*models.Secret, error) {
	version, err := r.client.Get(versionKey(hash)).Int64()
	if err != nil {
		return nil, err
	}

	// get secret from the hash map
	str, err := r.client.HGet(hashSecretKey, hash).Result()
	if err != nil {
		return nil, err
	}

	var sec models.SecretBase
	if err := json.Unmarshal([]byte(str), &sec); err != nil {
		return nil, err
	}
	return &models.Secret{
		SecretBase: sec,
		Version:    version,
	}, nil
}

// CreateSecret creates a new secret
func (r *RedisDB) CreateSecret(hash string, s models.Secret) error {
	str, err := json.Marshal(s.SecretBase)
	if err != nil {
		return err
	}

	// open transaction
	tx := r.client.TxPipeline()
	defer tx.Discard()

	// set data
	if err := tx.HSet(hashSecretKey, hash, str).Err(); err != nil {
		return err
	}
	// set version
	if err := tx.Set(versionKey(hash), s.Version, 0).Err(); err != nil {
		return err
	}
	// execute transaction
	_, err = tx.Exec()
	return err
}

// DeleteSecret removes existing secret
func (r *RedisDB) DeleteSecret(hash string) error {
	if err := r.client.HDel(hash).Err(); err != nil {
		return err
	}
	r.client.Del(versionKey(hash))
	return nil
}

// UpdateSecret decreases secret view counter
func (r *RedisDB) UpdateSecret(hash string, s models.Secret) error {
	str, err := json.Marshal(s.SecretBase)
	if err != nil {
		return err
	}

	// execute transaction by watching at version id
	return r.client.Watch(func(tx *redis.Tx) error {
		// get version id
		// it must be the same as version id in the incoming data set
		// otherwise the data was changed by concurrent session
		if versionCurrent, err := tx.Get(versionKey(hash)).Int64(); err != nil {
			return err
		} else if versionCurrent != s.Version {
			return errors.New("secret was modified from another session. Try again")
		}

		// change the data
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			// set data
			if err := pipe.HSet(hashSecretKey, hash, str).Err(); err != nil {
				return err
			}

			// increment version
			if err := pipe.Incr(versionKey(hash)).Err(); err != nil {
				return err
			}
			return nil
		})
		return err
	}, versionKey(hash))
}

// versionKey returns a Redis version counter key
func versionKey(hash string) string {
	return fmt.Sprintf("%s:%s", hashVersionKey, hash)
}
