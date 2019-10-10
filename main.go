package main

import (
	"fmt"
	"goterm/cmd/delete"
	"goterm/cmd/exec"
	"goterm/cmd/list"
	"goterm/cmd/set"
	"goterm/db"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:     "goterm",
		Version: fmt.Sprintf("v1.0"),
	}
	root.AddCommand(
		exec.NewCommand(),
		set.NewCommand(),
		list.NewCommand(),
		delete.NewCommand(),
	)

	f := root.UsageFunc()
	root.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println("Config Path:")
		fmt.Println("  " + db.GetFilePath())
		fmt.Printf("  export %s=[WHERE YOU WANT TO SAVE]\n", db.PathEnvKey)
		fmt.Println("")
		return f(command)
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
