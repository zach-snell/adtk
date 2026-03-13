package devops

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// Attachment represents a work item attachment.
type Attachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Size     int64  `json:"size,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

// AttachmentReference is the reference returned after uploading.
type AttachmentReference struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ListWorkItemAttachments returns the attachments (relations of type AttachedFile) on a work item.
func (c *Client) ListWorkItemAttachments(project string, id int) ([]WorkItemRelation, error) {
	wi, err := c.GetWorkItem(project, id, "Relations")
	if err != nil {
		return nil, fmt.Errorf("getting work item relations: %w", err)
	}

	if wi.Fields == nil {
		return nil, nil
	}

	// Relations come as an array in the raw JSON at "relations"
	relData, err := c.Get(project, fmt.Sprintf("/wit/workitems/%d", id), url.Values{"$expand": {"Relations"}})
	if err != nil {
		return nil, fmt.Errorf("getting work item with relations: %w", err)
	}

	var full struct {
		Relations []WorkItemRelation `json:"relations"`
	}
	if err := json.Unmarshal(relData, &full); err != nil {
		return nil, fmt.Errorf("unmarshaling relations: %w", err)
	}

	var attachments []WorkItemRelation
	for _, rel := range full.Relations {
		if rel.Rel == "AttachedFile" {
			attachments = append(attachments, rel)
		}
	}

	return attachments, nil
}

// UploadAttachment uploads a file as an attachment and returns the reference.
func (c *Client) UploadAttachment(project, fileName string, content []byte) (*AttachmentReference, error) {
	query := url.Values{}
	query.Set("fileName", fileName)

	requestURL := c.buildURL(HostMain, project, "/wit/attachments", query)
	resp, err := c.do("POST", requestURL, content, "application/octet-stream")
	if err != nil {
		return nil, fmt.Errorf("uploading attachment: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading upload response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, data)
	}

	var ref AttachmentReference
	if err := json.Unmarshal(data, &ref); err != nil {
		return nil, fmt.Errorf("unmarshaling attachment reference: %w", err)
	}

	return &ref, nil
}

// LinkAttachmentToWorkItem adds a previously uploaded attachment to a work item.
func (c *Client) LinkAttachmentToWorkItem(project string, workItemID int, attachmentURL, comment string) (*WorkItem, error) {
	ops := []JSONPatchOp{
		{
			Op:   "add",
			Path: "/relations/-",
			Value: map[string]interface{}{
				"rel": "AttachedFile",
				"url": attachmentURL,
				"attributes": map[string]string{
					"comment": comment,
				},
			},
		},
	}
	return c.UpdateWorkItem(project, workItemID, ops)
}

// DownloadAttachment downloads an attachment by its URL.
func (c *Client) DownloadAttachment(attachmentURL string) ([]byte, error) {
	resp, err := c.GetAbsolute(attachmentURL)
	if err != nil {
		return nil, fmt.Errorf("downloading attachment: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading attachment content: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, data)
	}

	return data, nil
}
