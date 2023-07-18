package drivers

import (
	"fmt"
	"strings"
	"tugboat/internal/driver"
	"tugboat/internal/drivers/docker"

	log "github.com/sirupsen/logrus"
)

func NewDriver(driverName string, opts driver.DriverOptions) (driver.Driver, error) {
	switch strings.ToLower(driverName) {
	case "docker":
		return docker.NewDockerDriver(opts)
	case "auto":
		return autoDiscover(opts)
	default:
		return nil, fmt.Errorf("unsupported driver name: %s", driverName)
	}
}

// Attempt to auto discover what container engine is being utilized
func autoDiscover(opts driver.DriverOptions) (driver.Driver, error) {
	log.Debug("Attempting to match a driver")

	// we only support docker for the moment, no need to look for it
	log.Debug("Initializing the docker driver")
	driver, err := docker.NewDockerDriver(opts)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
