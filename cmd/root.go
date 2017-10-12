// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

var gexp glob.Glob
var exp *regexp.Regexp

var red = color.New(color.FgRed).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "fd [OPTIONS] [<pattern>] [<path>]",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]
		path := args[1]

		// hidden, _ := cmd.Flags().GetBool("hidden")
		// cs, _ := cmd.Flags().GetBool("case-sensitive")
		// extension, _ := cmd.Flags().GetBool("extension")

		// create simple glob
		exp = regexp.MustCompile(pattern)
		gexp = glob.MustCompile(pattern)

		// 1. list files in path
		// 2. check if file is dir
		// 3. if so, do from 1
		// 4. else check for pattern match
		// 5. if matches print them as result
		// 6. stop once all files are visited in path

		find(path, walker(path))
	},
}

func find(path string, paths <-chan os.FileInfo) {
	for p := range paths {
		tp := filepath.Join(path, p.Name())
		if exp.MatchString(tp) || gexp.Match(tp) {
			fmt.Printf("%s\n", blue(tp))
			continue
		}
		if p.IsDir() {
			find(tp, walker(tp))
		}
	}
}

func walk(path string, ch chan os.FileInfo) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("%v\n", red(err))
		return
	}
	for _, f := range files {
		ch <- f
	}
}

func walker(path string) <-chan os.FileInfo {
	ch := make(chan os.FileInfo)
	go func() {
		walk(path, ch)
		close(ch)
	}()

	return ch
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(red(err))
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("hidden", "H", false, "Search hidden files and directories")
	RootCmd.Flags().BoolP("no-ignore", "I", false, "Do not respect .(git)ignore files")
	RootCmd.Flags().BoolP("case-sensitive", "s", false, "Case-sensitive search (default: smart case)")
	RootCmd.Flags().BoolP("absolute-path", "a", false, "Show absolute instead of relative paths")
	RootCmd.Flags().BoolP("follow", "L", false, "Follow symbolic links")
	RootCmd.Flags().BoolP("full-path", "p", false, "Search full path (default: file-/dirname only)")
	RootCmd.Flags().BoolP("print0", "0", false, "Separate results by the null character")
	RootCmd.Flags().BoolP("version", "V", false, "Prints version information")
	RootCmd.Flags().StringP("max-depth", "d", "none", "Set maximum search depth (default: none)")
	RootCmd.Flags().StringP("type", "t", "none", "Filter by type: f(ile), d(irectory), s(ymlink)")
	RootCmd.Flags().StringP("extension", "e", "none", "Filter by file extension")
	RootCmd.Flags().StringP("color", "c", "auto", "When to use colors: never, *auto*, always")
	RootCmd.Flags().StringP("threads", "j", "auto", "Set number of threads to use for searching")
}
