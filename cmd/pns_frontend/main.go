package main

import (
	"github.com/L-LYR/pns/internal/admin/frontend"
)

func main() {
	frontend.MustRegisterFrontendRouters()
	frontend.Run()
}
