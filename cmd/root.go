package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zsbahtiar/lionparcel-test/config"
)

var cfg *config.Config
var rootCmd = &cobra.Command{
	Use:   "lionparcel test - @zsbahtiar",
	Short: "this project for lionparcel - platform engineer",
}

func init() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd())

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
