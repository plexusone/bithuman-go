package omniavatar

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	bithumansdk "github.com/plexusone/bithuman-go"
	"github.com/plexusone/bithuman-go/api"

	"github.com/plexusone/omniavatar-core/render"
)

// RenderConfig configures the bitHuman render (batch video generation) provider.
type RenderConfig struct {
	// APIKey is the bitHuman API key.
	// Required.
	APIKey string

	// BaseURL is the bitHuman API base URL.
	// Default: https://api.bithuman.ai
	BaseURL string

	// AgentID is the default bitHuman agent used when
	// GenerateRequest.AvatarID is empty.
	AgentID string

	// HTTPClient is an optional custom HTTP client, used for both API
	// calls and video downloads.
	HTTPClient *http.Client
}

// RenderProvider implements render.Provider for bitHuman video generation.
// It also implements render.AudioUploader via the bitHuman file API.
type RenderProvider struct {
	sdk        *bithumansdk.Client
	agentID    string
	httpClient *http.Client
}

// Compile-time interface checks.
var (
	_ render.Provider      = (*RenderProvider)(nil)
	_ render.AudioUploader = (*RenderProvider)(nil)
)

// NewRenderProvider creates a bitHuman render provider.
func NewRenderProvider(cfg RenderConfig) (*RenderProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("%w: APIKey is required", render.ErrInvalidConfig)
	}

	opts := []bithumansdk.Option{bithumansdk.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, bithumansdk.WithBaseURL(cfg.BaseURL))
	}
	if cfg.HTTPClient != nil {
		opts = append(opts, bithumansdk.WithHTTPClient(cfg.HTTPClient))
	} else {
		opts = append(opts, bithumansdk.WithTimeout(30*time.Second))
	}

	sdk, err := bithumansdk.NewClient(opts...)
	if err != nil {
		return nil, render.NewProviderError("bithuman", "new_render_provider", err)
	}

	return &RenderProvider{
		sdk:        sdk,
		agentID:    cfg.AgentID,
		httpClient: cfg.HTTPClient,
	}, nil
}

// Name returns the provider name.
func (p *RenderProvider) Name() string { return "bithuman" }

// Generate submits a video generation job to bitHuman.
//
// GenerateRequest.AvatarID maps to the bitHuman agent ID. Width, Height,
// and Background are not supported by the bitHuman API and are ignored.
// The "voice_id" extension selects a TTS voice for Script input.
func (p *RenderProvider) Generate(ctx context.Context, req render.GenerateRequest) (*render.Job, error) {
	if req.AvatarID == "" {
		req.AvatarID = p.agentID
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	apiReq := &api.CreateVideoRequest{
		AgentID: req.AvatarID,
	}
	if req.AudioURL != "" {
		u, err := url.Parse(req.AudioURL)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid AudioURL: %w", render.ErrInvalidRequest, err)
		}
		apiReq.AudioURL = api.NewOptURI(*u)
	} else {
		apiReq.Text = api.NewOptString(req.Script)
		if voiceID := req.GetString("voice_id", ""); voiceID != "" {
			apiReq.VoiceID = api.NewOptString(voiceID)
		}
	}

	job, err := p.sdk.Videos().Create(ctx, apiReq)
	if err != nil {
		return nil, render.NewProviderError("bithuman", "generate", err)
	}

	return &render.Job{ID: job.ID, Provider: "bithuman"}, nil
}

// Status returns the current status of a generation job.
func (p *RenderProvider) Status(ctx context.Context, jobID string) (*render.JobStatus, error) {
	job, err := p.sdk.Videos().Get(ctx, jobID)
	if err != nil {
		return nil, render.NewProviderError("bithuman", "status", err)
	}
	return videoJobToStatus(job), nil
}

// Download streams the completed video to dst.
func (p *RenderProvider) Download(ctx context.Context, jobID string, dst io.Writer) error {
	status, err := p.Status(ctx, jobID)
	if err != nil {
		return err
	}
	if status.State != render.JobStateCompleted || status.VideoURL == "" {
		return fmt.Errorf("%w: job %s is %s", render.ErrJobNotCompleted, jobID, status.State)
	}

	if err := render.DownloadURL(ctx, p.httpClient, status.VideoURL, dst); err != nil {
		return render.NewProviderError("bithuman", "download", err)
	}
	return nil
}

// UploadAudio uploads audio content via the bitHuman file API and returns
// a hosted URL usable as GenerateRequest.AudioURL.
func (p *RenderProvider) UploadAudio(ctx context.Context, filename string, r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", render.NewProviderError("bithuman", "upload_audio", err)
	}

	apiReq := &api.UploadFileRequest{
		Base64:      api.NewOptString(base64.StdEncoding.EncodeToString(data)),
		Filename:    api.NewOptString(filename),
		ContentType: api.NewOptString(render.AudioContentType(filename)),
	}

	file, err := p.sdk.Files().Upload(ctx, apiReq)
	if err != nil {
		return "", render.NewProviderError("bithuman", "upload_audio", err)
	}
	if !file.URL.Set {
		return "", render.NewProviderError("bithuman", "upload_audio",
			fmt.Errorf("%w: upload response missing URL", render.ErrProviderUnavailable))
	}
	return file.URL.Value.String(), nil
}

// videoJobToStatus converts a bitHuman VideoJob to a normalized JobStatus.
func videoJobToStatus(job *api.VideoJob) *render.JobStatus {
	status := &render.JobStatus{
		ID:        job.ID,
		State:     mapVideoJobState(job.Status),
		RawStatus: string(job.Status),
	}
	if job.VideoURL.Set {
		status.VideoURL = job.VideoURL.Value.String()
	}
	if job.DurationSeconds.Set {
		status.Duration = job.DurationSeconds.Value
	}
	if job.Error.Set {
		status.ErrorMsg = job.Error.Value
	}
	return status
}

// mapVideoJobState maps bitHuman job statuses to normalized states.
func mapVideoJobState(s api.VideoJobStatus) render.JobState {
	switch s {
	case api.VideoJobStatusPending:
		return render.JobStatePending
	case api.VideoJobStatusProcessing:
		return render.JobStateProcessing
	case api.VideoJobStatusCompleted:
		return render.JobStateCompleted
	case api.VideoJobStatusFailed:
		return render.JobStateFailed
	default:
		// Unknown states stay non-terminal so pollers keep waiting.
		return render.JobStateProcessing
	}
}
