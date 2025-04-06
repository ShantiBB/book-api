package routers

import (
	"strconv"
)

func parseStringID(strID string) (*int, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
