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

var noTarget = errors.New("Error: No target opened.")
var noSettings = errors.New("Error: Settings not found.")
var invalidName = errors.New("Error: Invalid name. Please, don't use special chars.")

//@@@ TODO: newTarget func.
