package main

import (
	"github.com/szksh-lab/docfresh/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("docfresh", version, cli.Run)
}
