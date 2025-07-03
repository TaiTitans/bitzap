package main

import (
	_ "shortener-svc/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"shortener-svc/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
