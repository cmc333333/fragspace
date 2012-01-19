package rackspace

import (
  "http"
  "lib"
  "libs/conf"
  "appengine"
  "appengine/memcache"
  "appengine/urlfetch"
)

type Auth struct {
  Url string
  Token string
}
func retrieveAuth(context appengine.Context) *Auth {
  //  First, check memcache
  auth := Auth{}
  _, err := memcache.Gob.Get(context, "rackspace-auth", &auth)
  if err == memcache.ErrCacheMiss {
    return newAuth(context)
  } else if err != nil {
    context.Warningf("error retrieving item %v", err)
    return newAuth(context)
  }
  return &auth
}
func newAuth(context appengine.Context) *Auth {
  client := urlfetch.Client(context)
  c, err := conf.ReadConfigFile("config.ini")
  if err != nil {
    panic(lib.ServerError{"could not read config file: " + err.String()})
  }
  username, err := c.GetString("rackspace", "username")
  if err != nil {
    panic(lib.ServerError{"could not read username: " + err.String()})
  }
  key, err := c.GetString("rackspace", "key")
  if err != nil {
    panic(lib.ServerError{"could not read key: " + err.String()})
  }

  //  post an empty body
  req, err := http.NewRequest("GET", "https://auth.api.rackspacecloud.com/v1.0", nil)
  if err != nil {
    panic(lib.ServerError{"could not create request: " + err.String()})
  }
  req.Header.Add("x-Auth-Key", key)
  req.Header.Add("x-Auth-User", username)

  res, err := client.Do(req)
  if err != nil {
    panic(lib.ServerError{"problem contacting hosting provider: " + err.String()})
  }

  auth := Auth{res.Header.Get("X-Server-Management-Url"), res.Header.Get("X-Auth-Token")}
  err = memcache.Gob.Set(context, &memcache.Item{
    Key: "rackspace-auth",
    Object: auth,
  })
  if err != nil {
    context.Warningf("could not write to memcache: %v", err.String())
  }
  return &auth
}
