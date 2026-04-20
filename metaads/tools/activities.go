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

// RegisterActivities registers account activity and audit log tools.
func RegisterActivities(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "get_account_activities",
		Description: "Get recent activities for the ad account (budget changes, status changes, etc.).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"category":    {Type: "string", Description: "Filter by category (e.g. ACCOUNT, CAMPAIGN, AD, FUNDING, BID, BUDGET)."},
				"since":       {Type: "string", Description: "Start datetime (ISO 8601)."},
				"until":       {Type: "string", Description: "End datetime (ISO 8601)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			Category string `json:"category"`
			Since    string `json:"since"`
			Until    string `json:"until"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "event_time,event_type,object_id,object_name,actor_id,actor_name,extra_data,translated_event_type")
		if p.Category != "" {
			v.Set("category", p.Category)
		}
		if p.Since != "" {
			v.Set("since", p.Since)
		}
		if p.Until != "" {
			v.Set("until", p.Until)
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/activities", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_account_audit_log",
		Description: "Get the audit log for the ad account showing who made what changes.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"since":    {Type: "string", Description: "Start datetime (ISO 8601)."},
				"until":    {Type: "string", Description: "End datetime (ISO 8601)."},
				"actor_id": {Type: "string", Description: "Filter by actor (user) ID."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Since   string `json:"since"`
			Until   string `json:"until"`
			ActorID string `json:"actor_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "event_time,event_type,object_id,object_name,object_type,actor_id,actor_name,extra_data")
		if p.Since != "" {
			v.Set("since", p.Since)
		}
		if p.Until != "" {
			v.Set("until", p.Until)
		}
		if p.ActorID != "" {
			v.Set("uid", p.ActorID)
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/activities", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_activity",
		Description: "Get the activity log for a specific ad, ad set, or campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"object_id": {Type: "string", Description: "Ad, ad set, or campaign ID."},
				"since":     {Type: "string", Description: "Start datetime (ISO 8601)."},
				"until":     {Type: "string", Description: "End datetime (ISO 8601)."},
			}, paging.Properties()),
			Required: []string{"account", "object_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ObjectID string `json:"object_id"`
			Since    string `json:"since"`
			Until    string `json:"until"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "event_time,event_type,actor_id,actor_name,extra_data,translated_event_type")
		v.Set("object_id", p.ObjectID)
		if p.Since != "" {
			v.Set("since", p.Since)
		}
		if p.Until != "" {
			v.Set("until", p.Until)
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/activities", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})
}
