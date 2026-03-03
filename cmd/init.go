package cmd

import (
	"fmt"

	"github.com/dFragmented/labcheck/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure labcheck with your homelab details",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{}
		
		fmt.Print("Host (IP or hostname): ")
		fmt.Scanln(&cfg.Host)

		fmt.Print("Port (default 22): ")
		fmt.Scanln(&cfg.Port)
		if cfg.Port == 0 {
			cfg.Port = 22
		}

		fmt.Print("SSH User: ")
		fmt.Scanln(&cfg.User)

		fmt.Print("Path to SSH private key (e.g. ~/.ssh/id_rsa): ")
		fmt.Scanln(&cfg.KeyPath)

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		
		fmt.Println("Config saved!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}