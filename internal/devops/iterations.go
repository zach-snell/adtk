package devops

import (
	"encoding/json"
	"fmt"
	"time"
)

// ListIterations lists all iterations (sprints) for a project/team.
func (c *Client) ListIterations(project, team string) ([]Iteration, error) {
	path := "/work/teamsettings/iterations"
	var data []byte
	var err error

	if team != "" {
		// Team-scoped URL: {project}/{team}/_apis/work/teamsettings/iterations
		teamProject := fmt.Sprintf("%s/%s", project, team)
		data, err = c.Get(teamProject, path, nil)
	} else {
		data, err = c.Get(project, path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("listing iterations: %w", err)
	}

	var result IterationList
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return result.Value, nil
}

// GetIteration gets a single iteration by ID.
func (c *Client) GetIteration(project, team, iterationID string) (*Iteration, error) {
	path := fmt.Sprintf("/work/teamsettings/iterations/%s", iterationID)
	var data []byte
	var err error

	if team != "" {
		teamProject := fmt.Sprintf("%s/%s", project, team)
		data, err = c.Get(teamProject, path, nil)
	} else {
		data, err = c.Get(project, path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("getting iteration: %w", err)
	}

	var result Iteration
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCurrentIteration returns the current (active) iteration for a team.
func (c *Client) GetCurrentIteration(project, team string) (*Iteration, error) {
	iterations, err := c.ListIterations(project, team)
	if err != nil {
		return nil, err
	}

	for _, iter := range iterations {
		if iter.Attributes != nil && iter.Attributes.TimeFrame == "current" {
			return &iter, nil
		}
	}

	return nil, fmt.Errorf("no current iteration found")
}

// ListBoards lists all boards for a project/team.
func (c *Client) ListBoards(project, team string) ([]Board, error) {
	path := "/work/boards"
	var data []byte
	var err error

	if team != "" {
		teamProject := fmt.Sprintf("%s/%s", project, team)
		data, err = c.Get(teamProject, path, nil)
	} else {
		data, err = c.Get(project, path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("listing boards: %w", err)
	}

	var result BoardList
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return result.Value, nil
}

// GetBoard gets a single board by ID.
func (c *Client) GetBoard(project, team, boardID string) ([]byte, error) {
	path := fmt.Sprintf("/work/boards/%s", boardID)
	if team != "" {
		teamProject := fmt.Sprintf("%s/%s", project, team)
		return c.Get(teamProject, path, nil)
	}
	return c.Get(project, path, nil)
}

// GetBoardColumns gets the columns for a board.
func (c *Client) GetBoardColumns(project, team, boardID string) ([]BoardColumn, error) {
	path := fmt.Sprintf("/work/boards/%s/columns", boardID)
	var data []byte
	var err error

	if team != "" {
		teamProject := fmt.Sprintf("%s/%s", project, team)
		data, err = c.Get(teamProject, path, nil)
	} else {
		data, err = c.Get(project, path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("getting board columns: %w", err)
	}

	var result struct {
		Count int           `json:"count"`
		Value []BoardColumn `json:"value"`
	}
	if err := unmarshalJSON(data, &result); err != nil {
		return nil, err
	}
	return result.Value, nil
}

// CreateIteration creates a new iteration under a project's classification nodes.
func (c *Client) CreateIteration(project, name string, startDate, finishDate *time.Time) error {
	path := "/wit/classificationnodes/Iterations"
	body := map[string]interface{}{
		"name": name,
	}
	if startDate != nil || finishDate != nil {
		attrs := map[string]interface{}{}
		if startDate != nil {
			attrs["startDate"] = startDate.Format(time.RFC3339)
		}
		if finishDate != nil {
			attrs["finishDate"] = finishDate.Format(time.RFC3339)
		}
		body["attributes"] = attrs
	}
	_, err := c.Post(project, path, body)
	if err != nil {
		return fmt.Errorf("creating iteration: %w", err)
	}
	return nil
}

// GetTeamSettings gets team settings including default iteration and area path.
func (c *Client) GetTeamSettings(project, team string) (map[string]interface{}, error) {
	path := "/work/teamsettings"
	var data []byte
	var err error

	if team != "" {
		teamProject := fmt.Sprintf("%s/%s", project, team)
		data, err = c.Get(teamProject, path, nil)
	} else {
		data, err = c.Get(project, path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("getting team settings: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling team settings: %w", err)
	}
	return result, nil
}
