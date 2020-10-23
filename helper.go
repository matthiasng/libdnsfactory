package libdnsfactory

import (
	"fmt"
	"strconv"
)

func getValueString(key string, required bool, config map[string]string) (string, error) {
	val, ok := config[key]
	if !ok && required {
		return "", fmt.Errorf(`"%s" not set`, key)
	}

	return val, nil
}

func getValueInt(key string, required bool, config map[string]string) (int, error) {
	val, err := getValueString(key, required, config)
	if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return intVal, nil
}
