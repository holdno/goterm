package set

import (
	"errors"
	"fmt"
	"goterm/db"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	opt := NewOptions()
	cmd := &cobra.Command{
		Use:              "set [COMMAND NAME]",
		TraverseChildren: true,
		Short:            "execute your command",
		RunE: func(c *cobra.Command, args []string) error {
			if err := Run(args, *opt); err != nil {
				return err
			}
			return nil
		},
	}
	opt.AddFlags(cmd.Flags())
	return cmd
}

func Run(args []string, opts SetOptions) error {
	if len(args) != 1 {
		return errors.New("empty command name")
	}

	if strings.Contains(args[0], "#") {
		fmt.Println("Please rename again, without #")
		return nil
	}

	commandName := fmt.Sprintf("%s#%s", opts.Project, args[0])

	if opts.Command == "" {
		return fmt.Errorf("must set a command for %s", args[0])
	}

	configMap, err := db.Load()
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return nil
	}

	if configMap == nil {
		configMap = make(map[string]db.CommandDetail)
	}

	configMap[commandName] = db.CommandDetail{
		Command:  strings.Fields(opts.Command),
		Project:  opts.Project,
		Category: opts.Category,
		Remark:   opts.Remark,
		Mode:     opts.Mode,
	}

	if err = db.Save(configMap); err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
