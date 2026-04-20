package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterDiagnostics registers delivery diagnostics and recommendation tools.
func RegisterDiagnostics(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "get_delivery_insights",
		Description: "Get delivery diagnostics for the ad account, including delivery issues and cost insights.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
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
			v.Set("fields", "ad_object_id,name,results,cost_per_result,spend,impressions,delivery_info")
		}
		return cl.Get(ctx, fmt.Sprintf("/%s/delivery_estimate", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_recommendations",
		Description: "Get optimization recommendations for a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID."},
			},
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Get(ctx, "/"+p.CampaignID+"/adrules_governed", nil)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_adset_recommendations",
		Description: "Get optimization recommendations for an ad set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"adset_id": {Type: "string", Description: "Ad set ID."},
			},
			Required: []string{"account", "adset_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			AdsetID string `json:"adset_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "id,name,delivery_estimate")
		return cl.Get(ctx, "/"+p.AdsetID+"/delivery_estimate", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_recommendations",
		Description: "Get recommendations and relevance diagnostics for an ad.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"ad_id":   {Type: "string", Description: "Ad ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "ad_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			AdID    string `json:"ad_id"`
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
			v.Set("fields", "quality_ranking,engagement_rate_ranking,conversion_rate_ranking")
		}
		return cl.Get(ctx, "/"+p.AdID+"/insights", v)
	})
}
