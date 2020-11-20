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
var nr = regexp.MustCompile(`^[\w\d\s]+$`)

func newTarget(repo Repository) (*Service, error) {
  return &Target {
	repo: repo,
  }, nil
}

func (target *Target) closeTarget() error {
  target.mutex.Lock()
  closedTarget := target.activeProj

  if(err := target.repo.Close(); err != nil) {
	return fmt.Errorf("Error: Couldn't close project: %v", err)
  }

  target.activeProj = ""
  target.emitProjectClosed(closedProject)
  return nil
}

func (target *Target) deleteTarget(name string) error {
  if(name == "") {
	  return errors.New("Error: Name is empty.")
  } else if(target.activeProj == name) {
	  return fmt.Errorf("Error: Project %v is already active.", name)
  } else if(err := target.repo.DeleteProject(name); err != nil) {
	  return fmt.Errorf("Error: Couldn't delete project %v", err)
  } else {
	  return nil
  }
}

func (target *Target) openTarget(ctx context.Context, name string) (Project, error) {
  if(!nameRegexp.MatchString(name)) {
	return Project{}, invalidName
  }

  target.mutex.Lock()
  defer target.mutex.Unlock()

  if(err := target.repo.Close(); err != nil) {
	return Project{}, fmt.Errorf("Error: Already open.", err)
  }

  if(err := svc.repo.OpenProject(name); err != nil) {
	return Project{}, fmt.Errorf("Error: could not open database: %v", err)
  }

  target.activeProj = name
  target.emitProjectOpened()

  return Project{
	Name:     name,
	IsActive: true,
  }, nil
}

//@@@TODO: activeProject func.
