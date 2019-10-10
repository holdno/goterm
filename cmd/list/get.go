package list

import (
	"bufio"
	"fmt"
	"goterm/cmd/exec"
	"goterm/db"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type exectureInfo struct {
	cmdName string
	project string
}

var (
	counter  = 0
	indexMap = make(map[string]exectureInfo)
)

func NewCommand() *cobra.Command {
	opt := NewOptions()
	cmd := &cobra.Command{
		Use:              "list [OPTIONS]",
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

func Run(args []string, opts GetOptions) error {
	counter = 0
	var (
		listAll bool
	)
	if len(args) > 0 {
		if args[0] == "all" {
			listAll = true
		} else {
			return fmt.Errorf("The OPTIONS only support 'all' or empty")
		}
	}

	configMap, err := db.Load()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var names []string
	for name := range configMap {
		names = append(names, name)
	}
	sort.Strings(names)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Remark", "Name", "Command", "Project", "Category", "mode"})
	table.SetBorder(false)

	if listAll || opts.Project == "" && opts.Category == "" {
		for _, name := range names {
			table.Append(parseCommandDetailToSlice(name, configMap[name]))
		}
	} else {
		for _, name := range names {
			v := configMap[name]
			if opts.Project != "" && v.Project != opts.Project {
				continue
			}
			if opts.Category != "" && v.Category != opts.Category {
				continue
			}

			table.Append(parseCommandDetailToSlice(name, configMap[name]))
		}
	}

	table.Render()

	fmt.Println("\nPlease input the command ID to execute or input 'q' to exit")
	fmt.Print("Input: ")

	input := bufio.NewScanner(os.Stdin)
	// 逐行扫描
	for input.Scan() {
		// 获取用户输入
		line := input.Text()

		switch line {
		case "exit":
			return nil
		case "q":
			return nil
		default:
			if cmd, exist := indexMap[line]; exist {
				exec.Run([]string{cmd.cmdName}, exec.ExecOptions{Project: cmd.project})
				goto FINISH
			}
			fmt.Println("Please input the right command ID")
			fmt.Print("Input: ")
		}
	}
FINISH:
	return nil
}

func parseCommandDetailToSlice(name string, v db.CommandDetail) []string {
	counter++
	indexStr := fmt.Sprintf("%d", counter)
	cmdName := strings.Split(name, "#")[1]
	indexMap[indexStr] = exectureInfo{
		cmdName: cmdName,
		project: v.Project,
	}
	return []string{indexStr, v.Remark, cmdName, strings.Join(v.Command, " "), v.Project, v.Category, v.Mode}
}
