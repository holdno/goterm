package set

import (
	"github.com/spf13/pflag"
)

type SetOptions struct {
	Project  string
	Category string
	Remark   string
	Command  string
	Mode     string
}

func NewOptions() *SetOptions {
	return &SetOptions{}
}

func (o *SetOptions) AddFlags(flagSet *pflag.FlagSet) {
	// Add flags for generic options
	flagSet.StringVarP(&o.Command, "command", "a", "", "Execute command")
	flagSet.StringVarP(&o.Mode, "mode", "m", "string", "If the -m option is string, then commands are read from string. If there are arguments after the string, they are assigned to the positional parameters, starting with $0")
	flagSet.StringVarP(&o.Project, "project", "p", "", "The project to which the command belongs")
	flagSet.StringVarP(&o.Category, "category", "c", "", "Set a category for the command, you can quickly find the command by category")
	flagSet.StringVarP(&o.Remark, "remark", "r", "", "Record the purpose of this command in case you forget it")
}
