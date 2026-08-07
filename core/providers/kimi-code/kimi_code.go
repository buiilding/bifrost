// Package kimicode implements the Kimi Code provider.
//
// Kimi Code exposes an OpenAI-compatible Chat Completions API, but it does not
// expose the OpenAI Responses or input-token-counting endpoints. This adapter
// reuses the OpenAI implementation with an internal operation policy so
// Responses requests are converted to Chat Completions and unsupported
// operations are rejected before they reach Kimi.
package kimicode

import (
	"strings"

	"github.com/maximhq/bifrost/core/providers/openai"
	"github.com/maximhq/bifrost/core/schemas"
)

const defaultBaseURL = "https://api.kimi.com/coding"

// KimiCodeProvider is the native Bifrost adapter for Kimi Code.
//
// The embedded OpenAI provider supplies the shared OpenAI-compatible request,
// streaming, response conversion, and tool-call behavior. Its internal custom
// provider policy is deliberately private to this adapter: it enables model
// listing and Chat Completions, enables Responses-to-Chat fallback, and rejects
// Count Tokens because Kimi Code does not implement that endpoint.
type KimiCodeProvider struct {
	*openai.OpenAIProvider
}

// NewKimiCodeProvider creates a Kimi Code provider with its managed endpoint
// and the operation policy supported by the upstream API.
func NewKimiCodeProvider(config *schemas.ProviderConfig, logger schemas.Logger) *KimiCodeProvider {
	providerConfig := *config
	providerConfig.CustomProviderConfig = &schemas.CustomProviderConfig{
		CustomProviderKey: string(schemas.KimiCode),
		BaseProviderType:  schemas.OpenAI,
		AllowedRequests: &schemas.AllowedRequests{
			ListModels:           true,
			ChatCompletion:       true,
			ChatCompletionStream: true,
		},
	}

	if strings.TrimSpace(providerConfig.NetworkConfig.BaseURL) == "" {
		providerConfig.NetworkConfig.BaseURL = defaultBaseURL
	}

	return &KimiCodeProvider{
		OpenAIProvider: openai.NewOpenAIProvider(&providerConfig, logger),
	}
}
