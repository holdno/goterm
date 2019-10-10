package exec

import "github.com/spf13/pflag"

// SetupOptions init RDD options
type ExecOptions struct {
	Project string
}

// New a function to return a inited SetupOptions
func NewOptions() *ExecOptions {
	return &ExecOptions{}
}

// AddFlags related to RDD options
func (o *ExecOptions) AddFlags(flagSet *pflag.FlagSet) {
	// Add flags for generic options
	flagSet.StringVarP(&o.Project, "project", "p", "", "execute command under the project")
}
