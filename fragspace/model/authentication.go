package model

import (
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
