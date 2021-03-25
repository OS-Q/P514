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

package update

import (
	"context"
	"os"

	"github.com/OS-Q/S04A/cli/errorcodes"
	"github.com/OS-Q/S04A/cli/feedback"
	"github.com/OS-Q/S04A/cli/instance"
	"github.com/OS-Q/S04A/cli/output"
	"github.com/OS-Q/S04A/commands"
	rpc "github.com/OS-Q/S04A/rpc/commands"
	"github.com/OS-Q/S04A/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand creates a new `update` command
func NewCommand() *cobra.Command {
	updateCommand := &cobra.Command{
		Use:     "update",
		Short:   "Updates the index of cores and libraries",
		Long:    "Updates the index of cores and libraries to the latest versions.",
		Example: "  " + os.Args[0] + " update",
		Args:    cobra.NoArgs,
		Run:     runUpdateCommand,
	}
	updateCommand.Flags().BoolVar(&updateFlags.showOutdated, "show-outdated", false, "Show outdated cores and libraries after index update")
	return updateCommand
}

var updateFlags struct {
	showOutdated bool
}

func runUpdateCommand(cmd *cobra.Command, args []string) {
	instance := instance.CreateInstanceIgnorePlatformIndexErrors()

	logrus.Info("Executing `arduino update`")

	err := commands.UpdateCoreLibrariesIndex(context.Background(), &rpc.UpdateCoreLibrariesIndexReq{
		Instance: instance,
	}, output.ProgressBar())
	if err != nil {
		feedback.Errorf("Error updating core and libraries index: %v", err)
		os.Exit(errorcodes.ErrGeneric)
	}

	if updateFlags.showOutdated {
		outdatedResp, err := commands.Outdated(context.Background(), &rpc.OutdatedReq{
			Instance: instance,
		})
		if err != nil {
			feedback.Errorf("Error retrieving outdated cores and libraries: %v", err)
		}

		// Prints outdated cores
		tab := table.New()
		tab.SetHeader("Core name", "Installed version", "New version")
		if len(outdatedResp.OutdatedPlatform) > 0 {
			for _, p := range outdatedResp.OutdatedPlatform {
				tab.AddRow(p.Name, p.Installed, p.Latest)
			}
			feedback.Print(tab.Render())
		}

		// Prints outdated libraries
		tab = table.New()
		tab.SetHeader("Library name", "Installed version", "New version")
		if len(outdatedResp.OutdatedLibrary) > 0 {
			for _, l := range outdatedResp.OutdatedLibrary {
				tab.AddRow(l.Library.Name, l.Library.Version, l.Release.Version)
			}
			feedback.Print(tab.Render())
		}
	}

	logrus.Info("Done")
}
