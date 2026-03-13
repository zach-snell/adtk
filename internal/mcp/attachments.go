package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

// ManageAttachmentsInput defines the input schema for the manage_attachments tool.
type ManageAttachmentsInput struct {
	Action     string `json:"action" jsonschema:"Action to perform: 'list', 'upload', 'download'"`
	ProjectKey string `json:"project_key,omitempty" jsonschema:"Project name"`
	WorkItemID int    `json:"work_item_id,omitempty" jsonschema:"Work item ID (required for list, upload)"`
	FilePath   string `json:"file_path,omitempty" jsonschema:"Absolute path to the file to upload (required for upload). Note: paths refer to the MCP server's filesystem."`
	Comment    string `json:"comment,omitempty" jsonschema:"Optional comment for the attachment (for upload)"`
	URL        string `json:"url,omitempty" jsonschema:"Attachment URL (for download)"`
}

// ManageAttachmentsHandler returns the handler for the manage_attachments tool.
func ManageAttachmentsHandler(c *devops.Client, enableWrites bool) func(context.Context, *sdkmcp.CallToolRequest, ManageAttachmentsInput) (*sdkmcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *sdkmcp.CallToolRequest, input ManageAttachmentsInput) (*sdkmcp.CallToolResult, any, error) {
		switch input.Action {
		case "list":
			if input.WorkItemID == 0 {
				return resultError("work_item_id is required for 'list' action")
			}
			rels, err := c.ListWorkItemAttachments(input.ProjectKey, input.WorkItemID)
			if err != nil {
				return resultError(fmt.Sprintf("listing attachments: %v", err))
			}
			return resultJSON(rels)

		case "upload":
			if !enableWrites {
				return resultError("upload action requires ADTK_ENABLE_WRITES=true")
			}
			if input.WorkItemID == 0 || input.FilePath == "" {
				return resultError("work_item_id and file_path are required for 'upload' action")
			}

			content, err := os.ReadFile(input.FilePath)
			if err != nil {
				return resultError(fmt.Sprintf("reading file %q: %v", input.FilePath, err))
			}

			fileName := filepath.Base(input.FilePath)
			ref, err := c.UploadAttachment(input.ProjectKey, fileName, content)
			if err != nil {
				return resultError(fmt.Sprintf("uploading attachment: %v", err))
			}

			comment := input.Comment
			if comment == "" {
				comment = "Uploaded via adtk"
			}
			_, err = c.LinkAttachmentToWorkItem(input.ProjectKey, input.WorkItemID, ref.URL, comment)
			if err != nil {
				return resultError(fmt.Sprintf("linking attachment to work item: %v", err))
			}

			return resultJSON(map[string]interface{}{
				"id":           ref.ID,
				"url":          ref.URL,
				"file_name":    fileName,
				"size":         len(content),
				"work_item_id": input.WorkItemID,
				"status":       "uploaded",
			})

		case "download":
			if input.URL == "" {
				return resultError("url is required for 'download' action")
			}
			content, err := c.DownloadAttachment(input.URL)
			if err != nil {
				return resultError(fmt.Sprintf("downloading attachment: %v", err))
			}

			downloadDir := filepath.Join(os.TempDir(), "adtk-downloads")
			if err := os.MkdirAll(downloadDir, 0o755); err != nil {
				return resultError(fmt.Sprintf("creating download directory: %v", err))
			}

			// Extract filename from URL or use generic
			fileName := "attachment"
			savePath := filepath.Join(downloadDir, fileName)
			if err := os.WriteFile(savePath, content, 0o600); err != nil {
				return resultError(fmt.Sprintf("saving attachment: %v", err))
			}

			return resultJSON(map[string]interface{}{
				"size":     len(content),
				"saved_to": savePath,
				"status":   "downloaded",
			})

		default:
			return resultError(fmt.Sprintf("unknown action: %s", input.Action))
		}
	}
}
