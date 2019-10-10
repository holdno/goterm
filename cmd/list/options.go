package list

import "github.com/spf13/pflag"

type GetOptions struct {
	Project  string
	Category string
}

func NewOptions() *GetOptions {
	return &GetOptions{}
}

func (o *GetOptions) AddFlags(flagSet *pflag.FlagSet) {
	// Add flags for generic options
	flagSet.StringVarP(&o.Project, "project", "p", "", "get command under the project")
	flagSet.StringVarP(&o.Category, "category", "c", "", "get command under the category")

}
