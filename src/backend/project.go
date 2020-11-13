package project

import (
  "sync"
  "context"
  "errors"
  "fmt"
  "log"
  "regexp"
)

type projOpen func(name string) error
type projClose func(name string) error

type Project struct {
  Name string
  IsActive bool
}

//@@@TODO: details struct.
