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

// RegisterCreatives registers ad creative management tools.
func RegisterCreatives(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── list_creatives ─────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_creatives",
		Description: "List all ad creatives in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated fields."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		acctID := p.AdAccountID
		if acctID == "" {
			acctID = client.AdAccountID()
		}
		qp := url.Values{}
		if p.Fields != "" {
			qp.Set("fields", p.Fields)
		} else {
			qp.Set("fields", "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,status,created_time")
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/adcreatives", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── create_creative ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name: "create_creative",
		Description: "Create a new ad creative. Pass the full creative body as JSON for maximum flexibility " +
			"(supports object_story_spec, asset_feed_spec, link_data, video_data, etc.).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"name":          {Type: "string", Description: "Creative name."},
				"body":          {Type: "string", Description: "Full creative JSON body. Must be valid JSON containing the creative spec fields (object_story_spec, asset_feed_spec, etc.)."},
			},
			Required: []string{"account", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string          `json:"account"`
			AdAccountID string          `json:"ad_account_id"`
			Name        string          `json:"name"`
			Body        json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		acctID := p.AdAccountID
		if acctID == "" {
			acctID = client.AdAccountID()
		}

		// Parse the body JSON and post as form params
		var bodyMap map[string]json.RawMessage
		if err := json.Unmarshal(p.Body, &bodyMap); err != nil {
			return nil, fmt.Errorf("body must be valid JSON: %w", err)
		}
		fp := url.Values{}
		if p.Name != "" {
			fp.Set("name", p.Name)
		}
		for k, v := range bodyMap {
			// Unquote string values, keep objects as JSON strings
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				fp.Set(k, s)
			} else {
				fp.Set(k, string(v))
			}
		}
		return client.Post(ctx, fmt.Sprintf("/%s/adcreatives", acctID), fp)
	})

	// ── get_creative ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_creative",
		Description: "Get details of a specific ad creative by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"creative_id": {Type: "string", Description: "Creative ID."},
				"fields":      {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "creative_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CreativeID string `json:"creative_id"`
			Fields     string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		qp := url.Values{}
		if p.Fields != "" {
			qp.Set("fields", p.Fields)
		} else {
			qp.Set("fields", "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,status,object_type,created_time")
		}
		return client.Get(ctx, fmt.Sprintf("/%s", p.CreativeID), qp)
	})

	// ── update_creative ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "update_creative",
		Description: "Update an existing ad creative. Only provided fields in body are changed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"creative_id": {Type: "string", Description: "Creative ID to update."},
				"name":        {Type: "string", Description: "New creative name."},
				"body":        {Type: "string", Description: "JSON with fields to update (e.g. object_story_spec changes)."},
			},
			Required: []string{"account", "creative_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			CreativeID string          `json:"creative_id"`
			Name       string          `json:"name"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		if p.Name != "" {
			fp.Set("name", p.Name)
		}
		if len(p.Body) > 0 && string(p.Body) != "null" {
			var bodyMap map[string]json.RawMessage
			if err := json.Unmarshal(p.Body, &bodyMap); err == nil {
				for k, v := range bodyMap {
					var s string
					if err := json.Unmarshal(v, &s); err == nil {
						fp.Set(k, s)
					} else {
						fp.Set(k, string(v))
					}
				}
			}
		}
		return client.Post(ctx, fmt.Sprintf("/%s", p.CreativeID), fp)
	})

	// ── delete_creative ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "delete_creative",
		Description: "Delete an ad creative by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"creative_id": {Type: "string", Description: "Creative ID to delete."},
			},
			Required: []string{"account", "creative_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CreativeID string `json:"creative_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/%s", p.CreativeID))
	})
}
