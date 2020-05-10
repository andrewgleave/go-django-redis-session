package drsession

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v7"
)

var (
	// ErrSessionNotFound will be returned if the session is not found in Redis
	ErrSessionNotFound = errors.New("Session not found")
	// ErrInvalidPayload suggests the payload stored in the session does not conform to the
	// expected format i.e. not JSON
	ErrInvalidPayload = errors.New("Invalid payload")
	// ErrEmptySession will be returned if the session is empty (no value)
	ErrEmptySession = errors.New("Empty session")
)

// SessionClient provides provides simple retrieval and deserialisation support for
// JSON serialised Djagno sessions stored in Redis
type SessionClient struct {
	client *redis.Client
}

// NewSessionClient returns a new instance of SessionClient
func NewSessionClient(options redis.Options) (*SessionClient, error) {
	c := redis.NewClient(&options)
	_, err := c.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &SessionClient{c}, nil
}

// Get returns the unmarshalled JSON stored in key's session
func (c *SessionClient) Get(key string) (map[string]interface{}, error) {
	val, err := c.client.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}
	return c.parse(val)
}

func (c *SessionClient) parse(val string) (map[string]interface{}, error) {
	if len(val) == 0 {
		return nil, ErrEmptySession
	}

	// decoded should be in the form: uuid:json
	dec, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	sep := []byte{':'}
	idx := bytes.Index(dec, sep)
	if idx == -1 {
		return nil, ErrInvalidPayload
	}

	// we only care about the json...
	data := dec[idx+1:]
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("%v : %w", err, ErrInvalidPayload)
	}
	return payload, nil
}
