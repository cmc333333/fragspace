package model

import (
  "crypto/rand"
  "encoding/base64"

  "appengine"
  "appengine/datastore"
)

type Session struct {
  Key string
  Errors []string
}

func NewSession(context appengine.Context) *Session {
  buff := make([]byte, 128)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  key := base64.StdEncoding.EncodeToString(buff)
  session := &Session{key, []string{}}
  dsKey := datastore.NewKey(context, "Session", key, 0, nil)
  if _, err := datastore.Put(context, dsKey, session); err != nil {
    panic(err)
  }
  return session
}
