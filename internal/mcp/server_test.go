package mcp

import (
	"context"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zach-snell/adtk/internal/devops"
)

func TestNewServer_RegistersTools(t *testing.T) {
	t.Parallel()
	c := devops.NewTestClient("http://localhost:0")
	s := newServer(c)
	if s == nil {
		t.Fatal("newServer returned nil")
	}
}

func TestNewServer_WithRealConstructor(t *testing.T) {
	t.Parallel()
	s := New("test-org", "test-pat")
	if s == nil {
		t.Fatal("New returned nil")
	}
}

func makePromptReq(args map[string]string) *sdkmcp.GetPromptRequest {
	return &sdkmcp.GetPromptRequest{
		Params: &sdkmcp.GetPromptParams{
			Arguments: args,
		},
	}
}

func TestSprintSummaryPrompt(t *testing.T) {
	t.Parallel()
	result, err := sprintSummaryPromptHandler(context.Background(), makePromptReq(map[string]string{"project": "MyProject", "team": "Alpha"}))
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty prompt result")
	}
	tc, ok := result.Messages[0].Content.(*sdkmcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Messages[0].Content)
	}
	if tc.Text == "" {
		t.Error("expected non-empty prompt text")
	}
}

func TestSprintSummaryPrompt_MissingProject(t *testing.T) {
	t.Parallel()
	_, err := sprintSummaryPromptHandler(context.Background(), makePromptReq(map[string]string{}))
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestPRReviewDigestPrompt(t *testing.T) {
	t.Parallel()
	result, err := prReviewDigestPromptHandler(context.Background(), makePromptReq(map[string]string{"project": "MyProject", "repo": "my-repo"}))
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty prompt result")
	}
}

func TestPRReviewDigestPrompt_MissingFields(t *testing.T) {
	t.Parallel()
	_, err := prReviewDigestPromptHandler(context.Background(), makePromptReq(map[string]string{"project": "MyProject"}))
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestPipelineHealthPrompt(t *testing.T) {
	t.Parallel()
	result, err := pipelineHealthPromptHandler(context.Background(), makePromptReq(map[string]string{"project": "MyProject", "pipeline_id": "42"}))
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty prompt result")
	}
}

func TestPipelineHealthPrompt_MissingProject(t *testing.T) {
	t.Parallel()
	_, err := pipelineHealthPromptHandler(context.Background(), makePromptReq(map[string]string{}))
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestReleaseReadinessPrompt(t *testing.T) {
	t.Parallel()
	result, err := releaseReadinessPromptHandler(context.Background(), makePromptReq(map[string]string{"project": "MyProject", "iteration": "Sprint 5"}))
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatal("expected non-empty prompt result")
	}
}

func TestReleaseReadinessPrompt_MissingProject(t *testing.T) {
	t.Parallel()
	_, err := releaseReadinessPromptHandler(context.Background(), makePromptReq(map[string]string{}))
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}
