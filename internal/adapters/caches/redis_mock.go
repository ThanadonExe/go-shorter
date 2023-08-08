package caches

import (
	"encoding/json"
	"errors"
	"time"
)

type RedisCacheMock struct {
	items map[string][]byte
}

func NewRedisCacheMock() (*RedisCacheMock, error) {
	return &RedisCacheMock{}, nil
}

func (r *RedisCacheMock) Get(key string, value interface{}) error {

	if data, ok := r.items[key]; ok {
		if err := json.Unmarshal([]byte(data), value); err != nil {
			return err
		}
	}

	return errors.New("no data")
}

func (r *RedisCacheMock) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	r.items[key] = data

	return nil
}

func (r *RedisCacheMock) Delete(key string) error {
	if _, ok := r.items[key]; ok {
		delete(r.items, key)
		return nil
	}

	return errors.New("not found")
}

func (r *RedisCacheMock) Close() error {
	return nil
}
