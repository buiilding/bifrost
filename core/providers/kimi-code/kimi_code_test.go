package kimicode

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
)

var _ schemas.Provider = (*KimiCodeProvider)(nil)

func TestNewKimiCodeProvider(t *testing.T) {
	provider := NewKimiCodeProvider(&schemas.ProviderConfig{}, nil)

	if provider.GetProviderKey() != schemas.KimiCode {
		t.Fatalf("expected provider key %s, got %s", schemas.KimiCode, provider.GetProviderKey())
	}
}

func TestKimiCodeRejectsCountTokens(t *testing.T) {
	provider := NewKimiCodeProvider(&schemas.ProviderConfig{}, nil)

	_, err := provider.CountTokens(nil, schemas.Key{}, nil)
	if err == nil {
		t.Fatal("expected Count Tokens to be rejected")
	}
	if err.Error == nil || err.Error.Code == nil || *err.Error.Code != "unsupported_operation" {
		t.Fatalf("expected unsupported_operation, got %#v", err)
	}
	if err.ExtraFields.Provider != schemas.KimiCode {
		t.Fatalf("expected error provider %s, got %s", schemas.KimiCode, err.ExtraFields.Provider)
	}
}

func TestKimiCodeListModelsUsesOpenAICompatibleEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Errorf("expected /v1/models, got %s", r.URL.Path)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("expected bearer authorization, got %q", got)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"data":[{"id":"k3-256k","object":"model"}]}`)
	}))
	defer server.Close()

	provider := NewKimiCodeProvider(&schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{BaseURL: server.URL},
	}, nil)
	ctx := schemas.NewBifrostContext(context.Background(), schemas.NoDeadline)
	response, err := provider.ListModels(ctx, []schemas.Key{{
		Value:  schemas.SecretVar{Val: "test-key"},
		Models: schemas.WhiteList{"*"},
	}}, &schemas.BifrostListModelsRequest{Unfiltered: true})
	if err != nil {
		t.Fatalf("ListModels returned error: %v", err)
	}
	if len(response.Data) != 1 || response.Data[0].ID != "kimi-code/k3-256k" {
		t.Fatalf("unexpected model response: %#v", response.Data)
	}
}
