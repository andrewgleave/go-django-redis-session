package session

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v7"
)

func TestParse(t *testing.T) {
	value := "YWFiYjp7ImEiOiJhIn0="
	c := PyRedisSessionClient{}
	p, err := c.parseSession(value)
	if err != nil {
		t.Errorf("failed to parse session: %v", err)
	}

	wantJSON := "{\"a\":\"a\"}"
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(wantJSON), &m); err != nil {
		t.Errorf("failed to unmarshal JSON: %v", err)
	}

	if !reflect.DeepEqual(m, p) {
		t.Errorf("wanted: %+v got: %+v", m, p)
	}
}

func TestGetSession(t *testing.T) {
	options := redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	c, err := NewPyRedisSessionClient(options)
	if err != nil {
		t.Errorf("failed to connect to redis: %v", err)
	}

	_, err = c.Get("3cyqr81yqbsprhtwdb7j9nbw8z76pxsv")
	if err != nil {
		t.Errorf("failed to get session: %v", err)
	}
}

func TestGetNoSession(t *testing.T) {
	options := redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	c, err := NewPyRedisSessionClient(options)
	if err != nil {
		t.Errorf("failed to connect to redis: %v", err)
	}

	_, err = c.Get("aa")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("Should have received ErrSessionNotFound: %v", err)
	}
}
