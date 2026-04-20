package tools

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterAuth registers token management and authentication tools.
func RegisterAuth(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "exchange_token",
		Description: "Exchange a short-lived token for a long-lived token via POST /oauth/access_token.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"app_id":       {Type: "string", Description: "Facebook App ID."},
				"app_secret":   {Type: "string", Description: "Facebook App Secret."},
				"short_token":  {Type: "string", Description: "Short-lived access token to exchange."},
				"grant_type":   {Type: "string", Description: "Grant type (default: fb_exchange_token).", Default: "fb_exchange_token"},
			},
			Required: []string{"account", "app_id", "app_secret", "short_token"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AppID      string `json:"app_id"`
			AppSecret  string `json:"app_secret"`
			ShortToken string `json:"short_token"`
			GrantType  string `json:"grant_type"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("client_id", p.AppID)
		v.Set("client_secret", p.AppSecret)
		v.Set("fb_exchange_token", p.ShortToken)
		gt := p.GrantType
		if gt == "" {
			gt = "fb_exchange_token"
		}
		v.Set("grant_type", gt)
		return cl.Post(ctx, "/oauth/access_token", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "debug_token",
		Description: "Debug an access token to inspect its properties via GET /debug_token.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"input_token": {Type: "string", Description: "Token to debug/inspect."},
			},
			Required: []string{"account", "input_token"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			InputToken string `json:"input_token"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("input_token", p.InputToken)
		return cl.Get(ctx, "/debug_token", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_token_info",
		Description: "Get information about the currently configured access token (scopes, expiry, app).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "id,name,email")
		return cl.Get(ctx, "/me", v)
	})
}
