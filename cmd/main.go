package main

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func main() {
	rootCmd.AddCommand(restApiCmd)

	err := rootCmd.Execute()
	util.Panic(err)
}
