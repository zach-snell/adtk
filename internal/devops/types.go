package devops

import (
	"fmt"
	"strings"
	"time"
)

// ADOTime handles Azure DevOps time fields that may lack timezone suffixes.
// ADO sometimes returns "0001-01-01T00:00:00" (no "Z") for unset dates.
type ADOTime struct {
	time.Time
}

// UnmarshalJSON handles both "2006-01-02T15:04:05Z" and "0001-01-01T00:00:00" formats.
func (t *ADOTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "" || s == "null" {
		t.Time = time.Time{}
		return nil
	}
	// Try standard RFC3339 first
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		t.Time = parsed
		return nil
	}
	// Try without timezone (ADO sometimes omits Z on zero dates)
	parsed, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		t.Time = parsed
		return nil
	}
	// Try with fractional seconds
	parsed, err = time.Parse("2006-01-02T15:04:05.9999999Z", s)
	if err == nil {
		t.Time = parsed
		return nil
	}
	return fmt.Errorf("cannot parse ADO time %q", s)
}

// MarshalJSON outputs in RFC3339 format.
func (t ADOTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}

// --- Work Items ---

// WorkItem represents an Azure DevOps work item.
type WorkItem struct {
	ID     int                    `json:"id"`
	Rev    int                    `json:"rev"`
	Fields map[string]interface{} `json:"fields"`
	URL    string                 `json:"url"`
	Links  map[string]interface{} `json:"_links,omitempty"`
}

// WorkItemList is the response from batch or list work item queries.
type WorkItemList struct {
	Count int        `json:"count"`
	Value []WorkItem `json:"value"`
}

// WorkItemRelation represents a link/relation on a work item.
type WorkItemRelation struct {
	Rel        string                 `json:"rel"`
	URL        string                 `json:"url"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// WorkItemComment represents a single comment on a work item.
type WorkItemComment struct {
	ID           int       `json:"id"`
	WorkItemID   int       `json:"workItemId"`
	Text         string    `json:"text"`
	CreatedBy    Identity  `json:"createdBy"`
	CreatedDate  time.Time `json:"createdDate"`
	ModifiedBy   Identity  `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

// WorkItemCommentList is the paginated list of comments.
type WorkItemCommentList struct {
	TotalCount int               `json:"totalCount"`
	Count      int               `json:"count"`
	Comments   []WorkItemComment `json:"comments"`
}

// WorkItemType represents a work item type definition.
type WorkItemType struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        struct {
		URL string `json:"url,omitempty"`
	} `json:"icon,omitempty"`
}

// WorkItemTypeList is the response for listing work item types.
type WorkItemTypeList struct {
	Count int            `json:"count"`
	Value []WorkItemType `json:"value"`
}

// JSONPatchOp represents a single JSON Patch operation for work item updates.
type JSONPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
}

// BatchWorkItemUpdate holds one item's ID and its patch operations for batch updates.
type BatchWorkItemUpdate struct {
	ID  int
	Ops []JSONPatchOp
}

// IterationWorkItems is the response for getting work items in an iteration.
type IterationWorkItems struct {
	WorkItemRelations []IterationWorkItemRelation `json:"workItemRelations"`
}

// IterationWorkItemRelation represents a work item reference within an iteration.
type IterationWorkItemRelation struct {
	Target *WorkItemReference `json:"target,omitempty"`
	Source *WorkItemReference `json:"source,omitempty"`
	Rel    string             `json:"rel,omitempty"`
}

// WorkItemReference is a lightweight reference to a work item (ID + URL).
type WorkItemReference struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// WIQLResult is the response from a WIQL query.
type WIQLResult struct {
	QueryType       string `json:"queryType"`
	QueryResultType string `json:"queryResultType"`
	AsOf            string `json:"asOf"`
	WorkItems       []struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	} `json:"workItems"`
}

// --- Identity ---

// Identity represents an Azure DevOps user/identity.
type Identity struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
	ImageURL    string `json:"imageUrl,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ConnectionData is the response from /_apis/connectionData.
type ConnectionData struct {
	AuthenticatedUser Identity `json:"authenticatedUser"`
	AuthorizedUser    Identity `json:"authorizedUser"`
	InstanceID        string   `json:"instanceId"`
}

// IdentityList is the response from identity search.
type IdentityList struct {
	Count int        `json:"count"`
	Value []Identity `json:"value"`
}

// --- Projects ---

// Project represents an Azure DevOps project.
type Project struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description,omitempty"`
	URL            string  `json:"url"`
	State          string  `json:"state"`
	Revision       int     `json:"revision"`
	Visibility     string  `json:"visibility"`
	LastUpdateTime ADOTime `json:"lastUpdateTime"`
}

// ProjectList is the response for listing projects.
type ProjectList struct {
	Count int       `json:"count"`
	Value []Project `json:"value"`
}

// Team represents an Azure DevOps team.
type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	ProjectName string `json:"projectName,omitempty"`
}

// TeamList is the response for listing teams.
type TeamList struct {
	Count int    `json:"count"`
	Value []Team `json:"value"`
}

// --- Repositories ---

// GitRepository represents a Git repository in Azure DevOps.
type GitRepository struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	URL           string  `json:"url"`
	RemoteURL     string  `json:"remoteUrl"`
	SSHURL        string  `json:"sshUrl,omitempty"`
	WebURL        string  `json:"webUrl,omitempty"`
	DefaultBranch string  `json:"defaultBranch,omitempty"`
	Size          int64   `json:"size"`
	Project       Project `json:"project"`
}

// GitRepositoryList is the response for listing repositories.
type GitRepositoryList struct {
	Count int             `json:"count"`
	Value []GitRepository `json:"value"`
}

// GitRef represents a Git reference (branch/tag).
type GitRef struct {
	Name     string   `json:"name"`
	ObjectID string   `json:"objectId"`
	Creator  Identity `json:"creator,omitempty"`
}

// GitRefList is the response for listing refs/branches.
type GitRefList struct {
	Count int      `json:"count"`
	Value []GitRef `json:"value"`
}

// GitItem represents an item (file/folder) in a Git repository.
type GitItem struct {
	ObjectID      string `json:"objectId"`
	GitObjectType string `json:"gitObjectType"`
	CommitID      string `json:"commitId"`
	Path          string `json:"path"`
	URL           string `json:"url"`
	Content       string `json:"content,omitempty"`
}

// GitTreeRef represents a tree reference in a repository.
type GitTreeRef struct {
	ObjectID  string    `json:"objectId"`
	URL       string    `json:"url"`
	TreeEntry []GitItem `json:"treeEntries"`
}

// --- Pull Requests ---

// PullRequest represents a pull request.
type PullRequest struct {
	PullRequestID int           `json:"pullRequestId"`
	Title         string        `json:"title"`
	Description   string        `json:"description,omitempty"`
	Status        string        `json:"status"`
	CreatedBy     Identity      `json:"createdBy"`
	CreationDate  ADOTime       `json:"creationDate"`
	ClosedDate    ADOTime       `json:"closedDate,omitempty"`
	SourceRefName string        `json:"sourceRefName"`
	TargetRefName string        `json:"targetRefName"`
	MergeStatus   string        `json:"mergeStatus,omitempty"`
	IsDraft       bool          `json:"isDraft"`
	Repository    GitRepository `json:"repository"`
	Reviewers     []Reviewer    `json:"reviewers,omitempty"`
	URL           string        `json:"url"`
}

// PullRequestList is the response for listing pull requests.
type PullRequestList struct {
	Count int           `json:"count"`
	Value []PullRequest `json:"value"`
}

// Reviewer represents a PR reviewer with their vote.
type Reviewer struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
	Vote        int    `json:"vote"` // 10=approved, 5=approved with suggestions, 0=no vote, -5=waiting, -10=rejected
	IsRequired  bool   `json:"isRequired,omitempty"`
}

// PRComment represents a comment on a pull request.
type PRComment struct {
	ID              int       `json:"id"`
	Content         string    `json:"content"`
	Author          Identity  `json:"author"`
	PublishedDate   time.Time `json:"publishedDate"`
	LastUpdatedDate time.Time `json:"lastUpdatedDate"`
	CommentType     string    `json:"commentType"`
	ParentCommentID int       `json:"parentCommentId,omitempty"`
}

// PRThread represents a comment thread on a pull request.
type PRThread struct {
	ID         int         `json:"id"`
	Comments   []PRComment `json:"comments"`
	Status     string      `json:"status"`
	Properties interface{} `json:"properties,omitempty"`
}

// PRThreadList is the response for listing PR comment threads.
type PRThreadList struct {
	Count int        `json:"count"`
	Value []PRThread `json:"value"`
}

// --- Iterations/Sprints ---

// Iteration represents a sprint/iteration.
type Iteration struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Path       string               `json:"path"`
	URL        string               `json:"url"`
	Attributes *IterationAttributes `json:"attributes,omitempty"`
}

// IterationAttributes holds start/finish dates for an iteration.
type IterationAttributes struct {
	StartDate  *time.Time `json:"startDate,omitempty"`
	FinishDate *time.Time `json:"finishDate,omitempty"`
	TimeFrame  string     `json:"timeFrame,omitempty"`
}

// IterationList is the response for listing iterations.
type IterationList struct {
	Count int         `json:"count"`
	Value []Iteration `json:"value"`
}

// --- Boards ---

// Board represents a Kanban board.
type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// BoardList is the response for listing boards.
type BoardList struct {
	Count int     `json:"count"`
	Value []Board `json:"value"`
}

// BoardColumn represents a column on a board.
type BoardColumn struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	ItemLimit     int               `json:"itemLimit"`
	StateMappings map[string]string `json:"stateMappings,omitempty"`
	ColumnType    string            `json:"columnType,omitempty"`
}

// --- Wiki ---

// WikiPage represents a wiki page.
type WikiPage struct {
	ID           int        `json:"id,omitempty"`
	Path         string     `json:"path"`
	Content      string     `json:"content,omitempty"`
	GitItemPath  string     `json:"gitItemPath,omitempty"`
	URL          string     `json:"url,omitempty"`
	RemoteURL    string     `json:"remoteUrl,omitempty"`
	Order        int        `json:"order,omitempty"`
	IsParentPage bool       `json:"isParentPage,omitempty"`
	SubPages     []WikiPage `json:"subPages,omitempty"`
}

// Wiki represents a wiki in a project.
type Wiki struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // "projectWiki" or "codeWiki"
	URL          string `json:"url"`
	RemoteURL    string `json:"remoteUrl,omitempty"`
	ProjectID    string `json:"projectId,omitempty"`
	RepositoryID string `json:"repositoryId,omitempty"`
}

// WikiList is the response for listing wikis.
type WikiList struct {
	Count int    `json:"count"`
	Value []Wiki `json:"value"`
}

// --- Pipelines ---

// Pipeline represents a pipeline definition.
type Pipeline struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Folder string `json:"folder,omitempty"`
	URL    string `json:"url"`
}

// PipelineList is the response for listing pipelines.
type PipelineList struct {
	Count int        `json:"count"`
	Value []Pipeline `json:"value"`
}

// PipelineRun represents a pipeline run.
type PipelineRun struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	State        string   `json:"state"`
	Result       string   `json:"result,omitempty"`
	CreatedDate  ADOTime  `json:"createdDate"`
	FinishedDate ADOTime  `json:"finishedDate,omitempty"`
	URL          string   `json:"url"`
	Pipeline     Pipeline `json:"pipeline"`
}

// PipelineRunList is the response for listing pipeline runs.
type PipelineRunList struct {
	Count int           `json:"count"`
	Value []PipelineRun `json:"value"`
}

// --- Search ---

// SearchResult is a generic search result.
type SearchResult struct {
	Count   int           `json:"count"`
	Results []interface{} `json:"results"`
}

// WorkItemSearchResult represents a work item search hit.
type WorkItemSearchResult struct {
	Count   int `json:"count"`
	Results []struct {
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
		Fields map[string]string `json:"fields"`
	} `json:"results"`
}

// CodeSearchResult represents a code search hit.
type CodeSearchResult struct {
	Count   int `json:"count"`
	Results []struct {
		FileName   string `json:"fileName"`
		Path       string `json:"path"`
		Repository struct {
			Name string `json:"name"`
		} `json:"repository"`
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
		Matches map[string][]struct {
			CharOffset int `json:"charOffset"`
			Length     int `json:"length"`
		} `json:"matches,omitempty"`
	} `json:"results"`
}

// --- Generic list wrapper ---

// ListResponse is a generic wrapper for paginated list responses.
type ListResponse[T any] struct {
	Count int `json:"count"`
	Value []T `json:"value"`
}
