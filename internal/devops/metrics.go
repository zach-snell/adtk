package devops

import (
	"fmt"
	"time"
)

// WorkItemMetrics holds computed lifecycle metrics for a work item.
type WorkItemMetrics struct {
	CycleTime         time.Duration            `json:"cycle_time"`
	LeadTime          time.Duration            `json:"lead_time"`
	TimeInStatus      map[string]time.Duration `json:"time_in_status"`
	CurrentStatus     string                   `json:"current_status"`
	StatusTransitions []StatusTransition       `json:"status_transitions"`
}

// StatusTransition records a state change.
type StatusTransition struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	At   time.Time `json:"at"`
}

// ComputeWorkItemMetrics calculates cycle time, lead time, and time-in-status
// from the revision history of a work item.
// Uses: GET /{project}/_apis/wit/workitems/{id}/updates
func (c *Client) ComputeWorkItemMetrics(project string, id int) (*WorkItemMetrics, error) {
	updates, err := c.GetWorkItemUpdates(project, id)
	if err != nil {
		return nil, fmt.Errorf("getting updates for metrics: %w", err)
	}

	transitions, createdAt, currentStatus := extractTransitions(updates)

	metrics := &WorkItemMetrics{
		CurrentStatus:     currentStatus,
		StatusTransitions: transitions,
		TimeInStatus:      computeTimeInStatus(transitions, currentStatus, createdAt),
	}

	metrics.LeadTime = computeLeadTime(transitions, createdAt)
	metrics.CycleTime = computeCycleTime(transitions)

	return metrics, nil
}

// extractTransitions walks the update history and pulls out state changes.
func extractTransitions(updates []map[string]interface{}) ([]StatusTransition, time.Time, string) {
	var transitions []StatusTransition
	var createdAt time.Time
	currentStatus := ""

	for _, u := range updates {
		fields, ok := u["fields"].(map[string]interface{})
		if !ok {
			continue
		}

		// Extract revision date
		revDate := extractRevisionDate(u)

		// Capture created date from the first revision
		if createdAt.IsZero() && !revDate.IsZero() {
			createdAt = revDate
		}

		stateChange, ok := fields["System.State"].(map[string]interface{})
		if !ok {
			continue
		}

		oldVal, _ := stateChange["oldValue"].(string)
		newVal, _ := stateChange["newValue"].(string)
		if newVal == "" {
			continue
		}

		currentStatus = newVal
		transitions = append(transitions, StatusTransition{
			From: oldVal,
			To:   newVal,
			At:   revDate,
		})
	}

	return transitions, createdAt, currentStatus
}

// extractRevisionDate parses the revisedDate from an update entry.
func extractRevisionDate(update map[string]interface{}) time.Time {
	dateStr, _ := update["revisedDate"].(string)
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try alternate format
		t, err = time.Parse("2006-01-02T15:04:05.999Z", dateStr)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}

// computeLeadTime returns Created → first Closed/Done transition.
func computeLeadTime(transitions []StatusTransition, createdAt time.Time) time.Duration {
	if createdAt.IsZero() {
		return 0
	}
	for _, t := range transitions {
		if isClosedState(t.To) && !t.At.IsZero() {
			return t.At.Sub(createdAt)
		}
	}
	return 0
}

// computeCycleTime returns first Active → last Closed/Done transition.
func computeCycleTime(transitions []StatusTransition) time.Duration {
	var firstActive time.Time
	var lastClosed time.Time

	for _, t := range transitions {
		if isActiveState(t.To) && firstActive.IsZero() {
			firstActive = t.At
		}
		if isClosedState(t.To) {
			lastClosed = t.At
		}
	}

	if firstActive.IsZero() || lastClosed.IsZero() {
		return 0
	}
	return lastClosed.Sub(firstActive)
}

// computeTimeInStatus calculates how long the work item spent in each status.
func computeTimeInStatus(transitions []StatusTransition, currentStatus string, createdAt time.Time) map[string]time.Duration {
	result := make(map[string]time.Duration)
	if len(transitions) == 0 {
		return result
	}

	// Walk transitions pairwise
	for i := 0; i < len(transitions); i++ {
		status := transitions[i].To
		start := transitions[i].At
		if start.IsZero() {
			continue
		}

		var end time.Time
		if i+1 < len(transitions) {
			end = transitions[i+1].At
		} else {
			// Still in this status
			end = time.Now()
		}
		if !end.IsZero() {
			result[status] += end.Sub(start)
		}
	}

	return result
}

func isClosedState(state string) bool {
	return state == "Closed" || state == "Done" || state == "Resolved" || state == "Removed"
}

func isActiveState(state string) bool {
	return state == "Active" || state == "In Progress" || state == "Committed"
}
