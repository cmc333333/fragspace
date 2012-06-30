package api

import (
  "http"

  "appengine"
  "appengine/datastore"

  fhttp "fragspace/http"
  "fragspace/httphelpers"
  "fragspace/model"
)

func init() {
  c := new(chargeHandler)
  c.Self = c
  http.Handle("/api/charge", c)
}

type chargeHandler struct { fhttp.BaseHandler }

type chargeReq struct {
  CardToken string `json:"card_token"`
}

func (handler *chargeHandler) Post(r *fhttp.JsonRequest) fhttp.Response {
  return httphelpers.ReqTrustedClient((*http.Request)(r), func(userId string) fhttp.Response {
    post := new(chargeReq)
    if err := r.Extract(post); err != nil || post.CardToken == "" {
      return fhttp.UserError("invalid json")
    }
    context := appengine.NewContext((*http.Request)(r))
    chargeKey := datastore.NewKey(context, "Charge", post.CardToken, 0, nil)
    //  Check if this charge already exists
    charge := new(model.Charge)
    if err := datastore.Get(context, chargeKey, charge); err != datastore.ErrNoSuchEntity {
      return fhttp.UserError("charge already exists")
    }
    charge = model.NewCharge(datastore.NewKey(context, "User", userId, 0, nil))
    if _, err := datastore.Put(context, chargeKey, charge); err != nil {
      return fhttp.ServerError(err.String())
    }
    return fhttp.Success{}
  })
}
