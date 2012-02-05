package oauth2

import (
  "crypto/rand"
  "encoding/hex"

  "appengine"
  "appengine/datastore"
  "appengine/memcache"
)

func newCodeKey(userKey *datastore.Key, context appengine.Context) string {
  buff := make([]byte, 32)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  key := hex.EncodeToString(buff)

  err := memcache.Gob.Set(context, &memcache.Item{
    Key: "oauth-code-" + key,
    Object: userKey,
  })
  if err != nil {
    context.Warningf("could not write to memcache: %v", err.String())
  }
  return key
}
