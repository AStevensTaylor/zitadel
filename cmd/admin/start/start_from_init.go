package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/admin/initialise"
	"github.com/zitadel/zitadel/cmd/admin/key"
	"github.com/zitadel/zitadel/cmd/admin/setup"
)

func NewStartFromInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-from-init",
		Short: "cold starts zitadel",
		Long: `cold starts ZITADEL.
First the minimum requirements to start ZITADEL are set up.
Second the initial events are created.
Last ZITADEL starts.

Requirements:
- cockroachdb`,
		Run: func(cmd *cobra.Command, args []string) {
			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			initialise.InitAll(initialise.MustNewConfig(viper.GetViper()))

			setupConfig := setup.MustNewConfig(viper.GetViper())
			setupSteps := setup.MustNewSteps(viper.New())
			setup.Setup(setupConfig, setupSteps, masterKey)

			startConfig := MustNewConfig(viper.GetViper())

			err = startZitadel(startConfig, masterKey)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
