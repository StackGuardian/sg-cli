// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewChatCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Manage AI assistant chats (infra2code and kyro) and their messages",
		Long:  `Manage AI assistant chats (infra2code and kyro) and their messages`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a chat
  list                   List chats
  read                   Get a chat
  update                 Update a chat
  delete                 Archive a chat
  commit-bundle          Commit an uploaded git bundle
  create-pr              Push the generated code and open a pull request
  pr-status              Show the status of the chat's pull request
  pr-sync-status         Compare the local bundle with the pull request
  artifacts              Work with a chat's artifacts through presigned URLs
  messages               Send and read the messages of a kyro chat`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newListCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newCommitBundleCmd(c),
		newCreatePrCmd(c),
		newPrStatusCmd(c),
		newPrSyncStatusCmd(c),
		newArtifactsCmd(c),
		newMessagesCmd(c),
	)
	return cmd
}
func newArtifactsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Work with a chat's artifacts through presigned URLs",
		Long:  `Work with a chat's artifacts through presigned URLs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   Get a presigned URL that lists the artifacts
  get-url                Get a presigned download URL for an artifact
  upload-url             Get a presigned upload URL for an artifact`)
		},
	}
	cmd.PersistentFlags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkPersistentFlagRequired("chat-id")
	cmd.AddCommand(
		newArtifactsListCmd(c),
		newArtifactsGetUrlCmd(c),
		newArtifactsUploadUrlCmd(c),
	)
	return cmd
}
func newMessagesCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Send and read the messages of a kyro chat",
		Long:  `Send and read the messages of a kyro chat`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Send a message (use --text or a payload)
  list                   List the messages of a chat, oldest first
  read                   Get a message
  cancel                 Stop a pending or streaming assistant message
  retry                  Retry an errored or stopped assistant message`)
		},
	}
	cmd.PersistentFlags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkPersistentFlagRequired("chat-id")
	cmd.AddCommand(
		newMessagesCreateCmd(c),
		newMessagesListCmd(c),
		newMessagesReadCmd(c),
		newMessagesCancelCmd(c),
		newMessagesRetryCmd(c),
	)
	return cmd
}
