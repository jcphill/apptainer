// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cli

import (
	"os"
	"os/user"
	"path"
	"strings"
	"text/template"

	"github.com/singularityware/singularity/src/docs"
	"github.com/singularityware/singularity/src/pkg/sylog"
	"github.com/singularityware/singularity/src/pkg/util/auth"
	"github.com/spf13/cobra"
)

// Global variables for singularity CLI
var (
	debug   bool
	silent  bool
	verbose bool
	quiet   bool
)

var (
	// TokenFile holds the path to the sylabs auth token file
	tokenFile string
	// authToken holds the sylabs auth token
	authToken, authWarning string
)

func init() {
	SingularityCmd.Flags().SetInterspersed(false)
	SingularityCmd.PersistentFlags().SetInterspersed(false)

	templateFuncs := template.FuncMap{
		"TraverseParentsUses": TraverseParentsUses,
	}
	cobra.AddTemplateFuncs(templateFuncs)

	SingularityCmd.SetHelpTemplate(docs.HelpTemplate)
	SingularityCmd.SetUsageTemplate(docs.UseTemplate)

	SingularityCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Print debugging information")
	SingularityCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Only print errors")
	SingularityCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all normal output")
	SingularityCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity +1")

	usr, err := user.Current()
	if err != nil {
		sylog.Fatalf("Couldn't determine user home directory: %v", err)
	}
	defaultTokenFile := path.Join(usr.HomeDir, ".singularity", "sylabs-token")
	// authToken priority default_file < env < file_flag
	SingularityCmd.Flags().StringVar(&tokenFile, "tokenfile", defaultTokenFile, "path to the file holding your sylabs authentication token")
	if val := os.Getenv("SYLABS_TOKEN"); val != "" {
		authToken = val
	}
	if i := strings.Compare(tokenFile, defaultTokenFile); i != 0 {
		authToken, authWarning = auth.ReadToken(tokenFile)
	}
	if authToken == "" {
		authToken, authWarning = auth.ReadToken(defaultTokenFile)
	}
}

// SingularityCmd is the base command when called without any subcommands
var SingularityCmd = &cobra.Command{
	TraverseChildren:      true,
	DisableFlagsInUseLine: true,
	Run: nil,

	Use:     docs.SingularityUse,
	Short:   docs.SingularityShort,
	Long:    docs.SingularityLong,
	Example: docs.SingularityExample,
}

// ExecuteSingularity adds all child commands to the root command and sets
// flags appropriately. This is called by main.main(). It only needs to happen
// once to the root command (singularity).
func ExecuteSingularity() {
	if err := SingularityCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// TraverseParentsUses walks the parent commands and outputs a properly formatted use string
func TraverseParentsUses(cmd *cobra.Command) string {
	if cmd.HasParent() {
		return TraverseParentsUses(cmd.Parent()) + cmd.Use + " "
	}

	return cmd.Use + " "
}
