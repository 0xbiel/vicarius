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

type Target struct {
  repo Repository
  activeProj string
  projOpenFns []projOpen
  projCloseFns []projClose
  mutex sync.RWMutex
}

//@@@ TODO: error vars.

//@@@ TODO: newTarget func.
