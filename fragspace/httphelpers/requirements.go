package httphelpers

import (
  "http"

  "appengine"
  "libs/conf"

  fhttp "fragspace/http"
  "fragspace/model"
  "fragspace/oauth2"
)

func ReqToken(req *http.Request, success func(string) fhttp.Response) fhttp.Response {
  token := oauth2.DecodeToken(req)
  if token == nil {
    return fhttp.UserError("invalid_token")
  }
  return success(token.User)
}

func ReqUser(req *http.Request, success func(*model.User) fhttp.Response) fhttp.Response {
  return ReqToken(req, func(keyStr string) fhttp.Response {
    user, err := model.UserFromKey(keyStr, appengine.NewContext(req))
    if err != nil {
      return fhttp.UserError("invalid_token")
    }
    return success(user)
  })
}

func ReqTrustedClient(req *http.Request, success func(string) fhttp.Response) fhttp.Response {
  config, err := conf.ReadConfigFile("config.ini")
  if err != nil {
    panic(err)
  }
  clientId, err := config.GetString("webclient", "clientId")
  if err != nil {
    panic(err)
  }
  token := oauth2.DecodeToken(req)
  if token == nil || token.Client != clientId {
    return fhttp.UserError("invalid_token")
  }
  return success(token.User)
}
