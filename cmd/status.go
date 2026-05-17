package cmd

import (
	"fmt"

	"github.com/rjdailey/labcheck/config"
	labssh "github.com/rjdailey/labcheck/ssh"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker container status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("run 'labcheck init' first: %w", err)
		}

		client, err := labssh.Connect(cfg.Host, cfg.Port, cfg.User, cfg.KeyPath)
		if err != nil {
			return err
		}
		defer client.Close()

		out, err := client.Run(`docker ps -a --format "{{.Names}}\t{{.Status}}\t{{.Image}}\t{{.Label \"com.docker.compose.project\"}}"`)
		if err != nil {
			return fmt.Errorf("failed to get container status: %w", err)
		}

		containers := parseContainers(out)
		return runTUI(containers)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
