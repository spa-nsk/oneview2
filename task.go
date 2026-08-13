package oneview

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	uriTasks = "/rest/tasks"
)

// Task is GET /rest/tasks/{id} (TaskResourceV2 and later).
type Task struct {
	Resource
	AssociatedResource      AssociatedResource `json:"associatedResource,omitempty"`
	AssociatedTaskURI       string             `json:"associatedTaskUri,omitempty"`
	CompletedSteps          int                `json:"completedSteps,omitempty"`
	ComputedPercentComplete int                `json:"computedPercentComplete,omitempty"`
	ExpectedDuration        int                `json:"expectedDuration,omitempty"`
	Hidden                  bool               `json:"hidden,omitempty"`
	Owner                   string             `json:"owner,omitempty"`
	ParentTaskURI           string             `json:"parentTaskUri,omitempty"`
	PercentComplete         int                `json:"percentComplete,omitempty"`
	ProgressUpdates         []ProgressUpdate   `json:"progressUpdates,omitempty"`
	TaskErrors              []APIError         `json:"taskErrors,omitempty"`
	TaskOutput              []string           `json:"taskOutput,omitempty"`
	TaskState               string             `json:"taskState,omitempty"`
	TaskStatus              string             `json:"taskStatus,omitempty"`
	TaskType                string             `json:"taskType,omitempty"`
	TotalSteps              int                `json:"totalSteps,omitempty"`
	UserInitiated           bool               `json:"userInitiated,omitempty"`
	Data                    map[string]any     `json:"data,omitempty"`
}

// ProgressUpdate is a task log line.
type ProgressUpdate struct {
	ID           int    `json:"id,omitempty"`
	StatusUpdate string `json:"statusUpdate,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
}

func (t *Task) terminal() bool {
	switch strings.ToLower(t.TaskState) {
	case "completed", "error", "killed", "terminated", "warning":
		return true
	default:
		return false
	}
}

func (t *Task) failed() bool {
	switch strings.ToLower(t.TaskState) {
	case "error", "killed", "terminated":
		return true
	default:
		return false
	}
}

// ListTasks returns tasks matching opts.
func (c *Client) ListTasks(ctx context.Context, opts ListOptions) (*Collection[Task], error) {
	return GetAll[Task](ctx, c, uriTasks, opts)
}

// GetTask returns a task by UUID or full URI.
func (c *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	var t Task
	if err := c.GetJSON(ctx, resourcePath(uriTasks, id), &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// WaitTask polls a task URI until it reaches a terminal state.
func (c *Client) WaitTask(ctx context.Context, taskURI string) (*Task, error) {
	return c.WaitTaskInterval(ctx, taskURI, 2*time.Second)
}

// WaitTaskInterval is WaitTask with a custom poll interval.
func (c *Client) WaitTaskInterval(ctx context.Context, taskURI string, interval time.Duration) (*Task, error) {
	if taskURI == "" {
		return nil, fmt.Errorf("oneview: empty task URI")
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		t, err := c.GetTask(ctx, taskURI)
		if err != nil {
			return nil, err
		}
		if t.terminal() {
			if t.failed() {
				msg := t.TaskStatus
				if msg == "" {
					msg = t.TaskState
				}
				if len(t.TaskErrors) > 0 && t.TaskErrors[0].Message != "" {
					msg = t.TaskErrors[0].Message
				}
				return t, fmt.Errorf("oneview: task %s %s: %s", t.URI, t.TaskState, msg)
			}
			return t, nil
		}
		select {
		case <-ctx.Done():
			return t, ctx.Err()
		case <-ticker.C:
		}
	}
}

// WaitResponse waits when the HTTP result is 202 Accepted with a Location task URI.
func (c *Client) WaitResponse(ctx context.Context, resp *Response) (*Task, error) {
	if resp == nil {
		return nil, nil
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	uri := resp.TaskURI()
	if uri == "" {
		return nil, nil
	}
	return c.WaitTask(ctx, uri)
}

func resourcePath(base, idOrURI string) string {
	if strings.HasPrefix(idOrURI, "/rest/") {
		return idOrURI
	}
	if strings.HasPrefix(idOrURI, "http://") || strings.HasPrefix(idOrURI, "https://") {
		return idOrURI
	}
	return joinPath(base, idOrURI)
}
