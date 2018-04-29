package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/opb/seekrits/kvp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "seekrits <secretsNames...> -- <command> [<arg...>]",
	Short: "Grab secrets from AWS SecretsManager and load into the env",
	Long: `Use seekrits as a prefix for a command to load secrets from
AWS SecretsManager into the Env variable context for 
the specified command`,
	Args: func(cmd *cobra.Command, args []string) error {
		dashIdx := cmd.ArgsLenAtDash()
		if dashIdx == -1 {
			return errors.New("please separate secretsNames and command with '--'. See usage")
		}
		if err := cobra.MinimumNArgs(1)(cmd, args[:dashIdx]); err != nil {
			return errors.Wrap(err, "at least one secretsNames must be specified")
		}
		if err := cobra.MinimumNArgs(1)(cmd, args[dashIdx:]); err != nil {
			return errors.Wrap(err, "must specify command to run. See usage")
		}
		return nil
	},
	Run: execRun,
}

func execRun(cmd *cobra.Command, args []string) {
	dashIdx := cmd.ArgsLenAtDash()
	secretsNames := args[:dashIdx]
	commandAndArgs := args[dashIdx:]
	command := commandAndArgs[0]

	var commandArgs []string
	if len(commandAndArgs) > 2 {
		commandArgs = commandAndArgs[2:]
	}

	envs, err := kvp.EnvPairs(secretsNames)
	if err != nil {
		log.Fatalln(err)
	}

	execPlatform(command, commandArgs, envs)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}