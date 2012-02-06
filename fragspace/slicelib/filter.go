package slicelib

func Filter (slice []string, fn func(string) bool) (goodList []string) {
  for _, str := range slice {
    if fn(str) {
      goodList = append(goodList, str)
    }
  }
  return goodList
}

func IsNonEmpty (str string) bool {
  if len(str) > 0 {
    return true
  }
  return false
}
