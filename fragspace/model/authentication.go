package model

import (
  "bytes"
  "gob"

  "appengine/datastore"
  "libs/passwordhash"
)

type Authentication struct {
  User *datastore.Key
  Token []byte
  Type string
}

func NewPasswordAuth(userKey *datastore.Key, password string) *Authentication {
  hashBuffer := bytes.NewBuffer(make([]byte, 0, 1000000))
  gob.NewEncoder(hashBuffer).Encode(passwordhash.New(password))
  return &Authentication{
    userKey,
    hashBuffer.Bytes(),
    "password",
  }
}
