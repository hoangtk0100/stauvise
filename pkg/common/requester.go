package common

import (
	"context"
	"strconv"

	"github.com/hoangtk0100/app-context/core"
)

func GetRequesterID(ctx context.Context) (uint64, error) {
	requester := core.GetRequester(ctx)
	return strconv.ParseUint(requester.GetID(), 10, 64)
}
