package scope

import (
	"context"
	"sync"
	project "./project"
	"net/http"
	"encoding/json"
	"fmt"
)

const mName = "scope"

type Repository interface {
  UpSettings(ctx context.Context, module string, settings interface{}) error
  FindSettingsByModule(ctx context.Context, module string, settings interface{}) error
}

type Scope struct {
  rules	[]Rule
  repo	Repository
  mutex	sync.RWMutex
}

type Rule struct {
  URL    *regexp.Regexp
  Header Header
  Body   *regexp.Regexp
}

//@@@TODO: Header struct
