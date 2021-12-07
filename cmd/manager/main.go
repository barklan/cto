package main

import (
	"flag"

	manager "github.com/barklan/cto/pkg/manager"
)

func main() {
	image := "registry.gitlab.com/nftgalleryx/nftgallery_backend/backend"
	generalCommand := flag.String("cmd", "none", "What command to execute.")
	flag.Parse()
	switch *generalCommand {
	case "deployStag":
		manager.Deploy("stag", image)
	case "deployProd":
		manager.Deploy("prod", image)
	}
}
