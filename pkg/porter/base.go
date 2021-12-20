package porter

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/caching"
	"github.com/barklan/cto/pkg/storage"
	"github.com/jmoiron/sqlx"
)

type Base struct {
	Config    *storage.InternalConfig
	MediaPath string
	R         *sqlx.DB
	Cache     caching.Cache
}

func InitBase(config *storage.InternalConfig, db *sqlx.DB) *Base {
	base := Base{}
	base.Config = config

	configEnvironment, ok := os.LookupEnv("CONFIG_ENV")
	if !ok {
		log.Panic("Config environment variable CONFIG_ENV must be specified.")
	}
	if configEnvironment == "dev" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}

		base.MediaPath = currentDir + "/.cache/media"
	} else {
		base.MediaPath = "/app/media"
	}

	return &base
}

// CreateMediaDirIfNotExists creates the directory in default media path.
// It can accept nested directory path, but all parent directories must
// exist. Returns full directory path.
func (b *Base) CreateMediaDirIfNotExists(dirname string) string {
	fullDirname := b.MediaPath + "/" + dirname
	_, err := os.Stat(fullDirname)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fullDirname, 0755)
		if errDir != nil {
			log.Panic(err)
		}
	}

	return fullDirname
}
