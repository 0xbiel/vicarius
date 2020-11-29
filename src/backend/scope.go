package scope

import "context"

type Repository interface {
  UpSettings(ctx context.Context, module string, settings interface{}) error
  FindSettingsByModule(ctx context.Context, module string, settings interface{}) error
}
