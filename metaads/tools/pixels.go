package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/internal/paging"
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterPixels registers pixel management tools.
func RegisterPixels(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_pixels",
		Description: "List Meta Pixels for the ad account via GET /{ad_account_id}/adspixels.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields (e.g. id,name,code,last_fired_time)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Fields  string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		if p.Fields != "" {
			v.Set("fields", p.Fields)
		} else {
			v.Set("fields", "id,name,code,creation_time,last_fired_time,is_unavailable")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/adspixels", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_pixel",
		Description: "Get details of a specific Meta Pixel.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"pixel_id": {Type: "string", Description: "Pixel ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "pixel_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PixelID string `json:"pixel_id"`
			Fields  string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		if p.Fields != "" {
			v.Set("fields", p.Fields)
		}
		return cl.Get(ctx, "/"+p.PixelID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_pixel_stats",
		Description: "Get event statistics for a Meta Pixel.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"pixel_id": {Type: "string", Description: "Pixel ID."},
				"start":    {Type: "string", Description: "Start date (YYYY-MM-DD). Defaults to last 7 days."},
				"end":      {Type: "string", Description: "End date (YYYY-MM-DD)."},
			},
			Required: []string{"account", "pixel_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PixelID string `json:"pixel_id"`
			Start   string `json:"start"`
			End     string `json:"end"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("aggregation", "event")
		if p.Start != "" {
			v.Set("start", p.Start)
		}
		if p.End != "" {
			v.Set("end", p.End)
		}
		return cl.Get(ctx, "/"+p.PixelID+"/stats", v)
	})
}
