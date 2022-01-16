package porter

import (
	"time"

	"github.com/barklan/cto/pkg/security"
	"go.uber.org/zap"
)

func makeUserIntegrationPass(base *Base, email string) string {
	pass := security.CharString(12)
	err := base.Cache.Set(pass, email, 72*time.Hour)
	if err != nil {
		base.Log.Error("failed to set new integration pass to cache", zap.Error(err))
	}
	return pass
}
