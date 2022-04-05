package middleware

import (
	"github.com/fusuwei/gspider/gspider"
)

func UA(ctx *gspider.Context) {
	ctx.Session.UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36"
	ctx.Next()
	//ctx.Abort()
}
