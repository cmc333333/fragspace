package model

import (
  "time"

  "appengine/datastore"
)

type Charge struct {
  Created datastore.Time
  User *datastore.Key
  Amount int
  Completed datastore.Time
}
func NewCharge(userKey *datastore.Key) *Charge {
  return &Charge{
    datastore.SecondsToTime(time.Seconds()),
    userKey,
    0,
    0,
  }
}
