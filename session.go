package session

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v7"
)

var (
	payloadSep = []byte{':'}

	// ErrSessionNotFound will be returned if the session key is not found in Redis
	ErrSessionNotFound = errors.New("Session not found")
	// ErrInvalidPayload suggests the payload stored in the session does not conform to the
	// expected format e.g. not JSON
	ErrInvalidPayload = errors.New("Invalid payload")
	// ErrEmptySession will be returned if the session is empty
	ErrEmptySession = errors.New("Empty session")
)

// PyRedisSessionClient provides GET and SET operations against Django Redis sessions
type PyRedisSessionClient struct {
	client *redis.Client
}

// NewPyRedisSessionClient returns a new instance of PyRedisSessionClient
func NewPyRedisSessionClient(options redis.Options) (*PyRedisSessionClient, error) {
	client := redis.NewClient(&options)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &PyRedisSessionClient{client}, nil
}

// Get returns the unmarshalled JSON stored in key's session
func (c *PyRedisSessionClient) Get(key string) (map[string]interface{}, error) {
	val, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}
	return c.parseSession(val)
}

func (c *PyRedisSessionClient) parseSession(val string) (map[string]interface{}, error) {
	if len(val) == 0 {
		return nil, ErrEmptySession
	}

	// decoded should be in the form: uuid:json
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	idx := bytes.Index(decoded, payloadSep)
	if idx == -1 {
		return nil, ErrInvalidPayload
	}

	// we only care about the json portion
	data := decoded[idx+1:]
	payload := make(map[string]interface{})
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("%v : %w", err, ErrInvalidPayload)
	}
	return payload, nil
}
