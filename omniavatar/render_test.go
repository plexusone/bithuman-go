package omniavatar

import (
	"testing"

	"github.com/plexusone/bithuman-go/api"

	"github.com/plexusone/omniavatar-core/render"
)

func TestMapVideoJobState(t *testing.T) {
	tests := []struct {
		in   api.VideoJobStatus
		want render.JobState
	}{
		{api.VideoJobStatusPending, render.JobStatePending},
		{api.VideoJobStatusProcessing, render.JobStateProcessing},
		{api.VideoJobStatusCompleted, render.JobStateCompleted},
		{api.VideoJobStatusFailed, render.JobStateFailed},
		{api.VideoJobStatus("unknown"), render.JobStateProcessing},
	}
	for _, tt := range tests {
		if got := mapVideoJobState(tt.in); got != tt.want {
			t.Errorf("mapVideoJobState(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestVideoJobToStatus(t *testing.T) {
	job := &api.VideoJob{
		ID:     "vid-1",
		Status: api.VideoJobStatusFailed,
		Error:  api.NewOptString("boom"),
	}
	status := videoJobToStatus(job)
	if status.ID != "vid-1" || status.State != render.JobStateFailed {
		t.Errorf("videoJobToStatus() = %+v, want failed vid-1", status)
	}
	if status.ErrorMsg != "boom" {
		t.Errorf("ErrorMsg = %q, want %q", status.ErrorMsg, "boom")
	}
	if status.RawStatus != "failed" {
		t.Errorf("RawStatus = %q, want %q", status.RawStatus, "failed")
	}
}
