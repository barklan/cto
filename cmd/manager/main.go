package main

import (
	"flag"

	manager "github.com/barklan/cto/pkg/manager"
)

var generalCommand = flag.String("cmd", "none", "What command to execute.")

func main() {
	image := "registry.gitlab.com/nftgalleryx/nftgallery_backend/backend"
	flag.Parse()
	switch *generalCommand {
	case "deployStag":
		manager.Deploy("stag", image)
	case "deployProd":
		manager.Deploy("prod", image)
	}
}
