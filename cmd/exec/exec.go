package exec

import (
	"bufio"
	"fmt"
	"goterm/db"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputIndex = 0
	inputCache []string
)

// NewCommand of RDD controller
func NewCommand() *cobra.Command {
	opt := NewOptions()
	cmd := &cobra.Command{
		Use:              "exec [OPTIONS] [COMMANDS]",
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

// Run will start running RDD controller and returns error if something unexpected
func Run(args []string, opts ExecOptions) error {
	configMap, err := db.Load()
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return nil
	}

	if len(args) != 1 {
		return fmt.Errorf("error command:%s", strings.Join(args, " "))
	}

	cmdName := args[0]

	if configMap == nil || len(configMap) == 0 {
		fmt.Printf("Command:%s is not exist\n", cmdName)
		fmt.Println("You must run 'goterm set' first to init command config")
		fmt.Println("Run 'goterm set -h' for usage.")
		return nil
	}

	var _cmd *db.CommandDetail

	for name, v := range configMap {
		if opts.Project != "" {
			if name == fmt.Sprintf("%s#%s", opts.Project, cmdName) {
				_cmd = &v
				break
			}
		} else {
			if strings.Contains(name, cmdName) {
				_cmd = &v
				break
			}
		}
	}

	if _cmd == nil {
		fmt.Printf("Command:%s is not exist\n", cmdName)
		fmt.Println("You can use:")
		for k, v := range configMap {
			if v.Project != "" {
				fmt.Printf("  %s -p %s\n", strings.Split(k, "_")[1], v.Project)
			} else {
				fmt.Printf("  %s\n", strings.Split(k, "_")[1])
			}
		}
		return nil
	}

	var cmd *exec.Cmd
	command := strings.Join(_cmd.Command, " ")
	cmd = exec.Command("sh", "-c", command)
	if _cmd.Mode == "" || _cmd.Mode == "notstring" {
		input := bufio.NewScanner(os.Stdin)
		// 逐行扫描
		for input.Scan() {
			// 获取用户输入
			line := input.Text()
			if line == "exit" {
				return nil
			}
			pushCommand(line)
			cmd1 := exec.Command("sh", "-c", line)
			result, err := cmd1.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(result))
		}
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()

	return nil
}

func pushCommand(command string) {
	if len(inputCache) >= 10 {
		inputCache = inputCache[1:]
	}
	inputCache = append(inputCache, command)
	inputIndex = len(inputCache)
}
