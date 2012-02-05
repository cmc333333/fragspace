package frontend

import (
  "crypto/rand"
  "encoding/base64"
  "http"

  "appengine"
  "appengine/datastore"
)

type Session struct {
  Key string
  Errors []string
  Changed bool `datastore:"-"`
}

func NewSession(context appengine.Context) *Session {
  buff := make([]byte, 128)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  key := base64.StdEncoding.EncodeToString(buff)
  session := &Session{key, []string{}, true}
  return session
}

func CurrentSession(w http.ResponseWriter, r *http.Request, context appengine.Context) *Session {
  cookie, err := r.Cookie("session")
  if err != nil { //  no such cookie yet
    session := NewSession(context)
    cookie := &http.Cookie{
      Name: "session",
      Value: session.Key,
    }
    w.Header().Add("Set-Cookie", cookie.String())
    return session
  }
  dsKey := datastore.NewKey(context, "Session", cookie.Value, 0, nil)
  var session Session
  err = datastore.Get(context, dsKey, &session)
  if err != nil {
    panic(err)
  }
  return &session
}

func WithSession(w http.ResponseWriter, r *http.Request, fn func(http.ResponseWriter, *http.Request, *Session)) {
  context := appengine.NewContext(r)
  session := CurrentSession(w, r, context)
  fn(w, r, session)
  if session.Changed {
    dsKey := datastore.NewKey(context, "Session", session.Key, 0, nil)
    if _, err := datastore.Put(context, dsKey, session); err != nil {
      panic(err)
    }
  }
}
