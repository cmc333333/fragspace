package oauth2

import (
  "crypto/rand"
  "encoding/hex"
  "time"

  "appengine"
  "appengine/memcache"
)

type codeStruct struct {
  ExpiresAt int64
  User string
  Client string
}

func newCodeKey(userKey string, clientKey string, context appengine.Context) string {
  buff := make([]byte, 32)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  key := hex.EncodeToString(buff)

  err := memcache.Gob.Set(context, &memcache.Item{
    Key: "oauth-code-" + key,
    Object: &codeStruct{
      time.Seconds() + 60,  //  lasts for 1 min
      userKey,
      clientKey,
    },
  })
  if err != nil {
    context.Warningf("could not write to memcache: %v", err.String())
  }
  return key
}
