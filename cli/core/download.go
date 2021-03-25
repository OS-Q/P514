// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package core

import (
	"context"
	"os"

	"github.com/OS-Q/S04A/cli/errorcodes"
	"github.com/OS-Q/S04A/cli/feedback"
	"github.com/OS-Q/S04A/cli/globals"
	"github.com/OS-Q/S04A/cli/instance"
	"github.com/OS-Q/S04A/cli/output"
	"github.com/OS-Q/S04A/commands/core"
	rpc "github.com/OS-Q/S04A/rpc/commands"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initDownloadCommand() *cobra.Command {
	downloadCommand := &cobra.Command{
		Use:   "download [PACKAGER:ARCH[@VERSION]](S)",
		Short: "Downloads one or more cores and corresponding tool dependencies.",
		Long:  "Downloads one or more cores and corresponding tool dependencies.",
		Example: "" +
			"  " + os.Args[0] + " core download arduino:samd       # to download the latest version of Arduino SAMD core.\n" +
			"  " + os.Args[0] + " core download arduino:samd@1.6.9 # for a specific version (in this case 1.6.9).",
		Args: cobra.MinimumNArgs(1),
		Run:  runDownloadCommand,
	}
	return downloadCommand
}

func runDownloadCommand(cmd *cobra.Command, args []string) {
	inst, err := instance.CreateInstance()
	if err != nil {
		feedback.Errorf("Error downloading: %v", err)
		os.Exit(errorcodes.ErrGeneric)
	}

	logrus.Info("Executing `arduino core download`")

	platformsRefs, err := globals.ParseReferenceArgs(args, true)
	if err != nil {
		feedback.Errorf("Invalid argument passed: %v", err)
		os.Exit(errorcodes.ErrBadArgument)
	}

	for i, platformRef := range platformsRefs {
		platformDownloadreq := &rpc.PlatformDownloadReq{
			Instance:        inst,
			PlatformPackage: platformRef.PackageName,
			Architecture:    platformRef.Architecture,
			Version:         platformRef.Version,
		}
		_, err := core.PlatformDownload(context.Background(), platformDownloadreq, output.ProgressBar())
		if err != nil {
			feedback.Errorf("Error downloading %s: %v", args[i], err)
			os.Exit(errorcodes.ErrNetwork)
		}
	}
}
