package handlers

import (
  "fmt"
  "http"
  "models"

  "appengine"
  "appengine/datastore"

  fhttp "fragspace/http"
)

func init() {
  u := new(userHandler)
  u.Self = u
  http.Handle("/user", u)
}

type userHandler struct { fhttp.BaseHandler }

type postResponse struct {
  Authentication string `json:"authentication"`
}

func (handler *userHandler) Post(r *fhttp.JsonRequest) fhttp.Response {
  user := new(models.User)
  err := r.Extract(user)
  if err != nil {
    return fhttp.UserError("invalid json")
  }
  if user.Nickname == "" {
    return fhttp.UserError("nickname cannot be empty")
  }

  context := appengine.NewContext((*http.Request)(r))
  userKey, err := datastore.Put(context, datastore.NewIncompleteKey(context, "User", nil), user)
  if err != nil {
    return fhttp.ServerError(err.String())
  }

  auth := models.NewAuth(userKey)
  _, err = datastore.Put(context, datastore.NewIncompleteKey(context, "Auth", nil), auth)
  if err != nil {
    return fhttp.ServerError(err.String())
  }
  return fhttp.JsonResponse{
    &postResponse{
      fmt.Sprintf("%x", auth.Public),
    },
  }
}

/*
func userHandler(w http.ResponseWriter, r *http.Request) {
  //context := appengine.NewContext(r)
  user := models.NewUser()
  //_, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Device", nil), device)
  fmt.Fprintf(w, "%x", user.Public)
}
*/

