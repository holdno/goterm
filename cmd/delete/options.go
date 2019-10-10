package delete

import "github.com/spf13/pflag"

type DeleteOptions struct {
	Project  string
	Category string
}

func NewOptions() *DeleteOptions {
	return &DeleteOptions{}
}

func (o *DeleteOptions) AddFlags(flagSet *pflag.FlagSet) {
	// Add flags for generic options
	flagSet.StringVarP(&o.Project, "project", "p", "", "get command under the project")
	flagSet.StringVarP(&o.Category, "category", "c", "", "get command under the category")
}
