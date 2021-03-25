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

package upgrade

import (
	"context"
	"os"

	"github.com/OS-Q/S04A/cli/core"
	"github.com/OS-Q/S04A/cli/errorcodes"
	"github.com/OS-Q/S04A/cli/feedback"
	"github.com/OS-Q/S04A/cli/instance"
	"github.com/OS-Q/S04A/cli/output"
	"github.com/OS-Q/S04A/commands"
	rpc "github.com/OS-Q/S04A/rpc/commands"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand creates a new `upgrade` command
func NewCommand() *cobra.Command {
	upgradeCommand := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrades installed cores and libraries.",
		Long:    "Upgrades installed cores and libraries to latest version.",
		Example: "  " + os.Args[0] + " upgrade",
		Args:    cobra.NoArgs,
		Run:     runUpgradeCommand,
	}

	core.AddPostInstallFlagsToCommand(upgradeCommand)
	return upgradeCommand
}

func runUpgradeCommand(cmd *cobra.Command, args []string) {
	inst, err := instance.CreateInstance()
	if err != nil {
		feedback.Errorf("Error upgrading: %v", err)
		os.Exit(errorcodes.ErrGeneric)
	}

	logrus.Info("Executing `arduino upgrade`")

	err = commands.Upgrade(context.Background(), &rpc.UpgradeReq{
		Instance:        inst,
		SkipPostInstall: core.DetectSkipPostInstallValue(),
	}, output.NewDownloadProgressBarCB(), output.TaskProgress())

	if err != nil {
		feedback.Errorf("Error upgrading: %v", err)
	}

	logrus.Info("Done")
}
