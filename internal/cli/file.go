package cli

import (
	"fmt"
	"os"

	"github.com/ducla/cloud-agent/internal/cli/client"
	"github.com/spf13/cobra"
)

func NewFileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Manage files",
		Long:  "Commands for file operations on the agent",
	}

	cmd.AddCommand(newFileListCommand())
	cmd.AddCommand(newFileUploadCommand())
	cmd.AddCommand(newFileDownloadCommand())
	cmd.AddCommand(newFileDeleteCommand())

	return cmd
}

func newFileListCommand() *cobra.Command {
	var path string
	
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			files, err := c.ListFiles(cmd.Context(), path)
			if err != nil {
				return fmt.Errorf("failed to list files: %w", err)
			}

			return printOutput(files, globalFlags.Output)
		},
	}
	
	cmd.Flags().StringVar(&path, "path", "/", "Directory path")
	
	return cmd
}

func newFileUploadCommand() *cobra.Command {
	var remotePath string
	
	cmd := &cobra.Command{
		Use:   "upload [local-file]",
		Short: "Upload a file to the agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()
			
			if err := c.UploadFile(cmd.Context(), file, remotePath); err != nil {
				return fmt.Errorf("failed to upload file: %w", err)
			}

			fmt.Printf("File uploaded successfully to %s\n", remotePath)
			return nil
		},
	}
	
	cmd.Flags().StringVar(&remotePath, "remote-path", "", "Remote file path (required)")
	cmd.MarkFlagRequired("remote-path")
	
	return cmd
}

func newFileDownloadCommand() *cobra.Command {
	var localPath string
	
	cmd := &cobra.Command{
		Use:   "download [remote-file]",
		Short: "Download a file from the agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			file, err := os.Create(localPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer file.Close()
			
			if err := c.DownloadFile(cmd.Context(), args[0], file); err != nil {
				return fmt.Errorf("failed to download file: %w", err)
			}

			fmt.Printf("File downloaded successfully to %s\n", localPath)
			return nil
		},
	}
	
	cmd.Flags().StringVar(&localPath, "local-path", "", "Local file path (required)")
	cmd.MarkFlagRequired("local-path")
	
	return cmd
}

func newFileDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [remote-file]",
		Short: "Delete a file on the agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(globalFlags.AgentURL, globalFlags.Token)
			
			if err := c.DeleteFile(cmd.Context(), args[0]); err != nil {
				return fmt.Errorf("failed to delete file: %w", err)
			}

			fmt.Printf("File %s deleted successfully\n", args[0])
			return nil
		},
	}
}
