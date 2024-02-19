package command

import "context"

func GetUserId(ctx context.Context) int32 {
	return ctx.Value("user").(int32)
}
