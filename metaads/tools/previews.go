package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterPreviews registers ad preview tools.
func RegisterPreviews(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_preview",
		Description: "Get an HTML preview of an existing ad via GET /{ad_id}/previews.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"ad_id":       {Type: "string", Description: "Ad ID to preview."},
				"ad_format":   {Type: "string", Description: "Preview format (e.g. DESKTOP_FEED_STANDARD, MOBILE_FEED_STANDARD, RIGHT_COLUMN_STANDARD, INSTAGRAM_STANDARD)."},
			},
			Required: []string{"account", "ad_id", "ad_format"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			AdID     string `json:"ad_id"`
			AdFormat string `json:"ad_format"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("ad_format", p.AdFormat)
		return cl.Get(ctx, "/"+p.AdID+"/previews", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "generate_preview",
		Description: "Generate an ad preview from a creative specification without creating an ad.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"creative":  {Type: "string", Description: "JSON creative spec object."},
				"ad_format": {Type: "string", Description: "Preview format (e.g. DESKTOP_FEED_STANDARD, MOBILE_FEED_STANDARD)."},
			},
			Required: []string{"account", "creative", "ad_format"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			Creative string `json:"creative"`
			AdFormat string `json:"ad_format"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("creative", p.Creative)
		v.Set("ad_format", p.AdFormat)
		return cl.Get(ctx, fmt.Sprintf("/%s/generatepreviews", cl.AdAccountID()), v)
	})
}
