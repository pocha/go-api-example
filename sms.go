package models

type (
  SMS struct {
    from, to uint64
    message string
  }
)
