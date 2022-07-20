/*
Copyright © 2022 Himanshu Shekhar <himanshu.kiit@gmail.com>
Code ownership is with Himanshu Shekhar. Use without modifications.
*/
package main

import (
	"fmt"

	_ "embed"

	"github.com/imhshekhar47/hs-taskmaster/skills-api/cmd"
)

//go:embed LICENSE
var LICENSE string

func main() {
	fmt.Println(LICENSE)
	cmd.Execute()
}
