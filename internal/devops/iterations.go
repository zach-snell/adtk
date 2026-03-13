package devops

import (
	"fmt"
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
