package bithuman

import (
	"context"
	"fmt"

	"github.com/plexusone/bithuman-go/api"
)

// AgentsService handles avatar agent management.
type AgentsService struct {
	client *api.Client
}

// List returns all agents.
func (s *AgentsService) List(ctx context.Context, params api.ListAgentsParams) (*api.AgentList, error) {
	res, err := s.client.ListAgents(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.AgentList); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Create creates a new avatar agent.
func (s *AgentsService) Create(ctx context.Context, req *api.CreateAgentRequest) (*api.Agent, error) {
	res, err := s.client.CreateAgent(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Agent); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific agent.
func (s *AgentsService) Get(ctx context.Context, agentID string) (*api.Agent, error) {
	res, err := s.client.GetAgent(ctx, api.GetAgentParams{AgentId: agentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Agent); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates an agent.
func (s *AgentsService) Update(ctx context.Context, agentID string, req *api.UpdateAgentRequest) (*api.Agent, error) {
	res, err := s.client.UpdateAgent(ctx, req, api.UpdateAgentParams{AgentId: agentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Agent); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes an agent.
func (s *AgentsService) Delete(ctx context.Context, agentID string) error {
	_, err := s.client.DeleteAgent(ctx, api.DeleteAgentParams{AgentId: agentID})
	return err
}

// Speak makes an agent speak the given text.
func (s *AgentsService) Speak(ctx context.Context, agentID string, req *api.SpeakRequest) (*api.SpeakResponse, error) {
	res, err := s.client.Speak(ctx, req, api.SpeakParams{AgentId: agentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.SpeakResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// SessionsService handles real-time conversation sessions.
type SessionsService struct {
	client *api.Client
}

// Create creates a new real-time session with an agent.
// Returns connection details for WebRTC/LiveKit.
func (s *SessionsService) Create(ctx context.Context, req *api.CreateSessionRequest) (*api.Session, error) {
	res, err := s.client.CreateSession(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Session); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific session.
func (s *SessionsService) Get(ctx context.Context, sessionID string) (*api.Session, error) {
	res, err := s.client.GetSession(ctx, api.GetSessionParams{SessionId: sessionID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Session); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// End ends an active session.
func (s *SessionsService) End(ctx context.Context, sessionID string) error {
	_, err := s.client.EndSession(ctx, api.EndSessionParams{SessionId: sessionID})
	return err
}

// CreateEmbedToken creates a short-lived JWT for website embedding.
func (s *SessionsService) CreateEmbedToken(ctx context.Context, req *api.CreateEmbedTokenRequest) (*api.EmbedToken, error) {
	res, err := s.client.CreateEmbedToken(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.EmbedToken); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// TTSService handles text-to-speech synthesis.
type TTSService struct {
	client *api.Client
}

// Synthesize converts text to audio.
// Returns audio data as bytes (audio/mpeg format).
func (s *TTSService) Synthesize(ctx context.Context, req *api.TTSRequest) (*api.TextToSpeechOK, error) {
	res, err := s.client.TextToSpeech(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.TextToSpeechOK); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// ListVoices returns available TTS voices.
func (s *TTSService) ListVoices(ctx context.Context) (*api.VoiceList, error) {
	res, err := s.client.ListVoices(ctx)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.VoiceList); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// VideosService handles talking video generation.
type VideosService struct {
	client *api.Client
}

// Create starts a new video generation job.
// This is async; poll Get() to check status.
func (s *VideosService) Create(ctx context.Context, req *api.CreateVideoRequest) (*api.VideoJob, error) {
	res, err := s.client.CreateVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.VideoJob); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns video job status and download URL when complete.
func (s *VideosService) Get(ctx context.Context, videoID string) (*api.VideoJob, error) {
	res, err := s.client.GetVideo(ctx, api.GetVideoParams{VideoId: videoID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.VideoJob); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// FilesService handles file uploads.
type FilesService struct {
	client *api.Client
}

// Upload uploads a file (image, video, audio, or document).
func (s *FilesService) Upload(ctx context.Context, req *api.UploadFileRequest) (*api.File, error) {
	res, err := s.client.UploadFile(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.File); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// BillingService handles account balance and usage.
type BillingService struct {
	client *api.Client
}

// GetBalance returns the account credit balance.
func (s *BillingService) GetBalance(ctx context.Context) (*api.Balance, error) {
	res, err := s.client.GetBalance(ctx)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Balance); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// WebhooksService handles webhook event notifications.
type WebhooksService struct {
	client *api.Client
}

// List returns all registered webhooks.
func (s *WebhooksService) List(ctx context.Context) (*api.WebhookList, error) {
	res, err := s.client.ListWebhooks(ctx)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.WebhookList); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Create registers a new webhook endpoint.
func (s *WebhooksService) Create(ctx context.Context, req *api.CreateWebhookRequest) (*api.Webhook, error) {
	res, err := s.client.CreateWebhook(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Webhook); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete removes a webhook.
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	_, err := s.client.DeleteWebhook(ctx, api.DeleteWebhookParams{WebhookId: webhookID})
	return err
}
