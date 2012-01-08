package lib

type UserError struct {
  Msg string
}
func (err *UserError) String() string {
  return err.Msg
}

type ServerError struct {
  Msg string
}
func (err *ServerError) String() string {
  return err.Msg
}
