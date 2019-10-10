package delete

import (
	"fmt"
	"goterm/db"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	opt := NewOptions()
	cmd := &cobra.Command{
		Use:              "delete [OPTIONS] [COMMANDS]",
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

func Run(args []string, opts DeleteOptions) error {
	var (
		deleteAll = false
		name      string
	)

	if len(args) > 0 {
		if args[0] == "all" {
			deleteAll = true
		} else {
			name = args[0]
		}
	} else {
		fmt.Println("Please enter the name of the task you want to delete")
		fmt.Printf("Run goterm delete 'name'")
		return nil
	}

	oldConfigMap, err := db.Load()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	newConfigMap := make(map[string]db.CommandDetail)

	if !deleteAll {
		for k, v := range oldConfigMap {
			if name == k && v.Category != "" && v.Category == opts.Category {
				continue
			}

			_name := fmt.Sprintf("%s#%s", v.Project, name)
			fmt.Println(_name, k)
			if _name == k {
				continue
			}

			newConfigMap[k] = v
		}
	}

	if err = db.Save(newConfigMap); err != nil {
		fmt.Println("Failed to delete command", err)
	}

	return nil
}
