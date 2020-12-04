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

type Header struct {
  Key	*regexp.Regexp
  Value	*regexp.Regexp
}

func New(repo Repository, projService *proj.Service) *Scope {
  scope := &Scope{
	repo: repo,
  }

  projService.OnProjectOpen(func(_ string) error {
	err := s.load(context.Background())

	if err == proj.ErrNoSettings {
	  return nil
	}
	if err != nil {
	  return fmt.Errorf("Error: could not load scope: %v", err)
	}
		return nil
	})

	projService.OnProjectClose(func(_ string) error {
		scope.unload()
		return nil
	})

	return scope
}

func (scope *Scope) Rules() []Rule {
  defer scope.mutex.RUnlock()
  return scope.rules
}

//@@@TODO: LoadScope function.
