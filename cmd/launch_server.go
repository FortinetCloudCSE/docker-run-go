package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"docker-run-go/dockerinternal"
	"github.com/spf13/cobra"
)

func getFlagString(cmd *cobra.Command, flagName string) string {
	value, _ := cmd.Flags().GetString(flagName)
	return value
}

var launchServerCmd = &cobra.Command{
	Use:   "launch-server",
	Short: "Launch the Hugo server container",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dockerinternal.ServerConfig{
			DockerImage:   getFlagString(cmd, "docker-image"),
			HostPort:      getFlagString(cmd, "host-port"),
			ContainerPort: getFlagString(cmd, "container-port"),
			WatchDir:      getFlagString(cmd, "watch-dir"),
		}

		// Ensure the watch directory is absolute.
		abs, err := filepath.Abs(cfg.WatchDir)
		if err == nil {
			cfg.WatchDir = abs
		}

		ctx := context.Background()

		containerID, err := dockerinternal.StartContainer(ctx, dockerClient, cfg)
		if err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			os.Exit(1)
		}

		if err := dockerinternal.AttachContainer(ctx, dockerClient, containerID); err != nil {
			fmt.Printf("Error attaching container: %v\n", err)
			os.Exit(1)
		}

		// Setup signal handling.
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\nReceived shutdown signal. Stopping container.")
			dockerinternal.StopAndRemoveContainer(dockerClient, containerID)
			os.Exit(0)
		}()

		// Start file watcher.
		dockerinternal.WatchAndRestart(ctx, dockerClient, cfg, &containerID)
	},
}

func init() {
	rootCmd.AddCommand(launchServerCmd)
	launchServerCmd.Flags().String("docker-image", "fortinet-hugo:latest", "Docker image to use")
	launchServerCmd.Flags().String("host-port", "1313", "Host port to expose")
	launchServerCmd.Flags().String("container-port", "1313", "Container port to expose")
	launchServerCmd.Flags().String("watch-dir", ".", "Directory to watch for file changes")
}
