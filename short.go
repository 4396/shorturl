package main

import (
	"crypto/md5"
	"errors"
	"fmt"

	"gopkg.in/redis.v3"
)

// Short type
type Short struct {
	addr string
	cli  *redis.Client
}

// NewShort new short object
func NewShort(addr string) *Short {
	return &Short{addr: addr}
}

// Close redis client
func (s *Short) Close() {
	if s.cli != nil {
		s.cli.Close()
		s.cli = nil
	}
}

func (s *Short) db() *redis.Client {
	if s.cli == nil {
		s.cli = redis.NewClient(&redis.Options{Addr: s.addr})
	}
	return s.cli
}

// Encode url, return a short url
func (s *Short) Encode(url string) (val string, err error) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	val, err = s.db().Get("h:" + hash).Result()
	if err == nil && val != "" {
		return
	}

	for i := 1; i <= len(hash); i++ {
		val, err = s.db().Get("s:" + hash[:i]).Result()
		if err == nil && val != "" {
			continue
		}
		val = hash[:i]
		err = nil
		break
	}
	if val == "" {
		err = errors.New("url encode failed")
		return
	}

	s.db().Set("s:"+val, url, 0)
	s.db().Set("h:"+hash, val, 0)
	return
}

// Decode val, return real url
func (s *Short) Decode(val string) (url string, err error) {
	url, err = s.db().Get("s:" + val).Result()
	return
}
