package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	app "github.com/datachainlab/cosmos-sdk-interchain-dns/example/simapp"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	encodingConfig := simapp.MakeEncodingConfig()
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	authclient.Codec = appCodec

	clientContext := client.Context{}.
		WithCodec(cdc).
		WithJSONMarshaler(appCodec).
		WithTxGenerator(encodingConfig.TxGenerator)

	rootCmd := &cobra.Command{
		Use:   "simappcli",
		Short: "Command line interface for interacting with simappd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// NOTE: StatusCommand always use FlagKeyringBackend
	rpccmd := rpc.StatusCommand()
	rpccmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "keyring backend")

	// Construct Root Command
	rootCmd.AddCommand(
		rpccmd,
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(clientContext),
		txCmd(clientContext),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(clientCtx client.Context) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(clientCtx.Codec),
		flags.LineBreak,
		rpc.ValidatorCommand(clientCtx.Codec),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(clientCtx.Codec),
		authcmd.QueryTxCmd(clientCtx.Codec),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, clientCtx)

	return queryCmd
}

func txCmd(clientCtx client.Context) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.NewSendTxCmd(clientCtx),
		flags.LineBreak,
		authcmd.GetSignCommand(clientCtx),
		authcmd.GetMultiSignCommand(clientCtx),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(clientCtx),
		authcmd.GetEncodeCommand(clientCtx),
		authcmd.GetDecodeCommand(clientCtx),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, clientCtx)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == authtypes.ModuleName || cmd.Use == banktypes.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
//func registerRoutes(rs *lcd.RestServer) {
//	clienttypes.RegisterRoutes(rs.CliCtx, rs.Mux)
//	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
//	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
//}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}