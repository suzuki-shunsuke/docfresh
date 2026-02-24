package main

import (
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/szksh-lab/docfresh/pkg/cli"
)

var version = ""

func main() {
	urfave.Main("docfresh", version, cli.Run)
}
