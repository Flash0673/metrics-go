package file

import (
	"fmt"
	"os"

	"github.com/Flash0673/metrics-go/pkg/logging/logger"
)

func Store(destination string, payload []byte) error {
	logger.Logger.Info(fmt.Sprintf("storing payload to file %s", destination))
	return os.WriteFile(destination, payload, 0666)
}

func Load(source string) ([]byte, error) {
	logger.Logger.Info(fmt.Sprintf("loading payload from file %s", source))
	return os.ReadFile(source)
}
