// Package omniavatar provides an OmniAvatar render provider backed by the
// bitHuman API.
//
// It implements render.Provider (Generate/Status/Download) for
// audio-driven talking-head generation using bitHuman agents, and
// render.AudioUploader via the bitHuman file API.
//
// The adapter is constructor-based and depends only on omniavatar-core.
// The batteries omniavatar package registers it; import that to use it by
// name:
//
//	renderer, err := omniavatar.GetRenderProvider("bithuman",
//	    omniavatar.WithAPIKey(os.Getenv("BITHUMAN_API_KEY")),
//	    omniavatar.WithExtension("agent_id", agentID))
//
// The real-time (live) bitHuman provider lives in the batteries omniavatar
// package (github.com/plexusone/omniavatar/providers/bithuman), because
// its LiveKit integration depends on that package.
package omniavatar
