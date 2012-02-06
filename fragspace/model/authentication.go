package model

import (
  "crypto/rand"

  "appengine/datastore"

  "libs/passwordhash"
)

type Authentication struct {
  User *datastore.Key
  Token []byte
  Type string
}

func NewPasswordAuth(userKey *datastore.Key, password string) *Authentication {
  return &Authentication{
    userKey,
    passwordHash(password),
    "password",
  }
}
func (auth *Authentication) PasswordHash() *passwordhash.PasswordHash {
  return decodePassword(auth.Token)
}

func NewOAuth2(userKey *datastore.Key) *Authentication {
  buff := make([]byte, 256)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }

  return &Authentication{
    userKey,
    buff,
    "oauth2",
  }
}
