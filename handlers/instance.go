package handlers

import (
  "libs/conf"
  "bytes"
  "http"
  "lib"
  "fmt"
  "appengine"
  "appengine/urlfetch"
)
func init() {
  i := new(instanceHandler)
  i.Self = i
  http.Handle("/instance", i)
}

type instanceHandler struct { lib.BaseHandler }

func (handler *instanceHandler) Post(w lib.JsonResponse, r *lib.JsonRequest) {
  context := appengine.NewContext((*http.Request)(r))
  client := urlfetch.Client(context)
  c, _ := conf.ReadConfigFile("config.ini")
  username, _ := c.GetString("rackspace", "username")
  key, _ := c.GetString("rackspace", "key")

  //  post an empty body
  req, _ := http.NewRequest("POST", "https://auth.api.rackspacecloud.com/v1.0", bytes.NewBuffer([]byte{}))
  req.Header.Add("x-Auth-Key", key)
  req.Header.Add("x-Auth-User", username)
  req.ContentLength = 0

  fmt.Print(client)
  fmt.Print(req)
  res, _ := client.Do(req)

  fmt.Print(res)
}
