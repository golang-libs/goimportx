package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/golang-libs/goimportx/pkg/importx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "goimportx",
	Short:   "sort and group go imports",
	Example: `goimportx --file /path/to/file.go --group "system,local,third"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := importx.InitGroup(group); err != nil {
			return err
		}

		var list []string
		for _, file := range files {
			list = append(list, strings.FieldsFunc(file, func(r rune) bool {
				return r == ',' || r == '|'
			})...)
		}

		for _, file := range list {
			result, err := importx.Sort(file, nil)
			if err != nil {
				return err
			}

			if write {
				_ = os.WriteFile(file, result, 0644)
			} else {
				_, _ = fmt.Fprint(os.Stdout, string(result))
			}
		}

		return nil
	},
}

var files []string
var group string
var write bool

func init() {
	rootCmd.Flags().StringSliceVarP(&files, "file", "f", nil, "file path")
	rootCmd.Flags().StringVarP(&group, "group", "g", "system,local,third", "group rule, split by comma, only supports [system,local,third,others]")
	rootCmd.Flags().BoolVarP(&write, "write", "w", false, "write result to (source) file instead of stdout")
}

func Execute() {
	rootCmd.Version = fmt.Sprintf(
		"%s %s/%s", "v0.0.1",
		runtime.GOOS, runtime.GOARCH)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
