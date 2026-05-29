// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// NewApp creates an app with sane defaults.
func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = ""
	app.Usage = usage
	app.Copyright = "Copyright 2013-2025 The go-ethereum Authors"
	app.Before = func(ctx *cli.Context) error {
		MigrateGlobalFlags(ctx)
		return nil
	}
	return app
}

var migrationApplied = map[*cli.Command]struct{}{}

// MigrateGlobalFlags makes all global flag values available in the
// context. This should be called as early as possible in app.Before.
//
// Example:
//
//	geth account new --keystore /tmp/mykeystore --lightkdf
//
// is equivalent after calling this method with:
//
//	geth --keystore /tmp/mykeystore --lightkdf account new
//
// i.e. in the subcommand Action function of 'account new', ctx.Bool("lightkdf)
// will return true even if --lightkdf is set as a global option.
//
// This function may become unnecessary when https://github.com/urfave/cli/pull/1245 is merged.
func MigrateGlobalFlags(ctx *cli.Context) {
	var iterate func(cs []*cli.Command, fn func(*cli.Command))
	iterate = func(cs []*cli.Command, fn func(*cli.Command)) {
		for _, cmd := range cs {
			if _, ok := migrationApplied[cmd]; ok {
				continue
			}
			migrationApplied[cmd] = struct{}{}
			fn(cmd)
			iterate(cmd.Subcommands, fn)
		}
	}

	// This iterates over all commands and wraps their action function.
	iterate(ctx.App.Commands, func(cmd *cli.Command) {
		if cmd.Action == nil {
			return
		}

		action := cmd.Action
		cmd.Action = func(ctx *cli.Context) error {
			doMigrateFlags(ctx)
			return action(ctx)
		}
	})
}

func doMigrateFlags(ctx *cli.Context) {
	// Figure out if there are any aliases of commands. If there are, we want
	// to ignore them when iterating over the flags.
	aliases := make(map[string]bool)
	for _, fl := range ctx.Command.Flags {
		for _, alias := range fl.Names()[1:] {
			aliases[alias] = true
		}
	}
	for _, name := range ctx.FlagNames() {
		for _, parent := range ctx.Lineage()[1:] {
			if parent.IsSet(name) {
				// When iterating across the lineage, we will be served both
				// the 'canon' and alias formats of all commands. In most cases,
				// it's fine to set it in the ctx multiple times (one for each
				// name), however, the Slice-flags are not fine.
				// The slice-flags accumulate, so if we set it once as
				// "foo" and once as alias "F", then both will be present in the slice.
				if _, isAlias := aliases[name]; isAlias {
					continue
				}
				// If it is a string-slice, we need to set it as
				// "alfa, beta, gamma" instead of "[alfa beta gamma]", in order
				// for the backing StringSlice to parse it properly.
				if result := parent.StringSlice(name); len(result) > 0 {
					ctx.Set(name, strings.Join(result, ","))
				} else {
					ctx.Set(name, parent.String(name))
				}
				break
			}
		}
	}
}

func init() {
	cli.FlagStringer = FlagString
}

// FlagString prints a single flag in help.
func FlagString(f cli.Flag) string {
	df, ok := f.(cli.DocGenerationFlag)
	if !ok {
		return ""
	}
	needsPlaceholder := df.TakesValue()
	placeholder := ""
	if needsPlaceholder {
		placeholder = "value"
	}

	namesText := cli.FlagNamePrefixer(df.Names(), placeholder)

	defaultValueString := ""
	if s := df.GetDefaultText(); s != "" {
		defaultValueString = " (default: " + s + ")"
	}
	envHint := strings.TrimSpace(cli.FlagEnvHinter(df.GetEnvVars(), ""))
	if envHint != "" {
		envHint = " (" + envHint[1:len(envHint)-1] + ")"
	}
	usage := strings.TrimSpace(df.GetUsage())
	usage = wordWrap(usage, 80)
	usage = indent(usage, 10)

	return fmt.Sprintf("\n    %-35s%-35s%s\n%s", namesText, defaultValueString, envHint, usage)
}

func indent(s string, nspace int) string {
	ind := strings.Repeat(" ", nspace)
	return ind + strings.ReplaceAll(s, "\n", "\n"+ind)
}

func wordWrap(s string, width int) string {
	var (
		output     strings.Builder
		lineLength = 0
	)

	for {
		sp := strings.IndexByte(s, ' ')
		var word string
		if sp == -1 {
			word = s
		} else {
			word = s[:sp]
		}
		wlen := len(word)
		over := lineLength+wlen >= width
		if over {
			output.WriteByte('\n')
			lineLength = 0
		} else {
			if lineLength != 0 {
				output.WriteByte(' ')
				lineLength++
			}
		}

		output.WriteString(word)
		lineLength += wlen

		if sp == -1 {
			break
		}
		s = s[wlen+1:]
	}

	return output.String()
}

// CheckExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func CheckExclusive(ctx *cli.Context, args ...any) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.Names()[0]

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag check, make sure value set doesn't conflict with passed in option
				if ctx.String(flag.Names()[0]) == option {
					name += "=" + option
					set = append(set, "--"+name)
				}
				// shift arguments and continue
				i++
				continue

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.IsSet(flag.Names()[0]) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		fmt.Fprintf(os.Stderr, "Flags %v can't be used at the same time", strings.Join(set, ", "))
		os.Exit(1)
	}
}
