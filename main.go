package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yxlimo/protoc-gen-go-asynq/internal/gen"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		gen.NewModule(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}
