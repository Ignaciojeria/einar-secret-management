package main

import (
	_ "secret-management/app/shared/configuration"
	"secret-management/app/shared/constants"
	_ "embed"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	_ "secret-management/app/onload"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}
