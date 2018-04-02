package main

import (
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
	minio "github.com/minio/minio/cmd"
)

// credSyncer syncs and manages funcinsts.
type credSyncer struct {
	pool *redis.Pool
}

var _ minio.CredentialProvider = (*credSyncer)(nil)

func (r *credSyncer) Get(key string) (cred minio.CustomCredentials, errCode minio.APIErrorCode) {
	c := r.pool.Get()
	defer c.Close()

	reply, err := c.Do("GET", key)
	if err != nil {
		log.Printf("[WARN] failed to get credential for %s, %v", key, err)
		errCode = minio.ErrInvalidAccessKeyID
		return
	}
	bts, ok := reply.([]byte)
	if !ok {
		log.Printf("failed to get credential bytes for %s", key)
		errCode = minio.ErrInvalidAccessKeyID
		return
	}
	if err = json.Unmarshal(bts, &cred); err != nil {
		log.Printf("failed to unmarshal credential for %s, %v", key, err)
		errCode = minio.ErrInvalidAccessKeyID
		return
	}
	return
}
