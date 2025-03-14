package cmd

import (
	"fmt"
	"os"

	"docker-run-go/dockerinternal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// buildImageCmd represents the `build-image` command using Docker SDK
var buildImageCmd = &cobra.Command{
	Use:   "build-image [author-dev | admin-dev]",
	Short: "Builds a Docker image programmatically using the Docker SDK",
	Long: `Builds a Docker image with the specified environment.

Example:
  docker-run-go build-image author-dev
  docker-run-go build-image admin-dev
`,
	Args: cobra.ExactArgs(1), // Require exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		envArg := args[0]

		// Map provided argument to actual Docker build target
		envMap := map[string]string{
			"author-dev": "prod",
			"admin-dev":  "dev",
		}
		env, exists := envMap[envArg]
		if !exists {
			fmt.Println("Usage: docker-run-go build-image [author-dev | admin-dev]")
			os.Exit(1)
		}

		// Determine the corresponding container name
		containerMap := map[string]string{
			"prod": "fortinet-hugo",
			"dev":  "hugotester",
		}
		containerName := containerMap[env]

		// Initialize Docker client
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Printf("Error creating Docker client: %v\n", err)
			os.Exit(1)
		}

		// Build the Docker image
		err = dockerinternal.BuildDockerImage(cli, containerName, env, envArg)
		if err != nil {
			fmt.Printf("Error building Docker image: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("**** Built a %s container named: %s ****\n", envArg, containerName)
	},
}

func init() {
	rootCmd.AddCommand(buildImageCmd)
}
