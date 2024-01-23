//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

var (
	binaryName string = "burningmoe"
	buildDir   string = "bin"
	cmd        string = fmt.Sprintf(".%[1]ccmd%[1]cweb", os.PathSeparator)
	output     string = fmt.Sprintf(".%[1]c%[2]s%[1]c%[3]s", os.PathSeparator, buildDir, binaryName)
)

func Build() error {
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", output, cmd)
}

func Run() error {
	fmt.Println("Running...")
	return sh.Run("go", "run", cmd)
}

func RunBinary() error {
	Build()
	fmt.Println("Running binary...")
	return sh.Run(output)
}

func Clean() error {
	fmt.Println("Cleaning...")
	return os.Remove(output)
}
