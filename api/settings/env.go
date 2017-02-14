package settings

import (
	"log"
	"os"
)

func SetPreproductionEnv() {

	err := os.Setenv("GO_ENV", "preproduction")
	if err != nil {
		log.Fatal(err)
	}
}
