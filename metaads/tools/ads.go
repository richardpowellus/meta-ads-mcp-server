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

// RegisterAds registers ad management tools.
func RegisterAds(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── list_ads ───────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_ads",
		Description: "List all ads in the ad account, optionally filtered by adset or campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"adset_id":      {Type: "string", Description: "Filter by ad set ID. Lists ads under this adset."},
				"campaign_id":   {Type: "string", Description: "Filter by campaign ID. Lists ads under this campaign."},
				"fields":        {Type: "string", Description: "Comma-separated fields."},
				"status":        {Type: "string", Description: "Filter by effective_status: ACTIVE, PAUSED, etc."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			AdsetID     string `json:"adset_id"`
			CampaignID  string `json:"campaign_id"`
			Fields      string `json:"fields"`
			Status      string `json:"status"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		qp := url.Values{}
		if p.Fields != "" {
			qp.Set("fields", p.Fields)
		} else {
			qp.Set("fields", "id,name,adset_id,campaign_id,status,effective_status,creative{id,name,thumbnail_url},created_time,updated_time")
		}
		if p.Status != "" {
			qp.Set("effective_status", fmt.Sprintf("[%q]", p.Status))
		}
		// Determine parent: adset > campaign > account
		parent := p.AdsetID
		if parent == "" {
			parent = p.CampaignID
		}
		if parent == "" {
			parent = p.AdAccountID
			if parent == "" {
				parent = client.AdAccountID()
			}
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/ads", parent), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── create_ad ──────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "create_ad",
		Description: "Create a new ad. Requires name, adset_id, creative, and status.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"ad_account_id":  {Type: "string", Description: "Ad account ID (overrides default)."},
				"name":           {Type: "string", Description: "Ad name."},
				"adset_id":       {Type: "string", Description: "Ad set ID to place this ad in."},
				"creative":       {Type: "string", Description: "JSON creative spec: {\"creative_id\":\"123\"} for existing creative, or inline spec."},
				"status":         {Type: "string", Description: "Initial status: ACTIVE or PAUSED."},
				"tracking_specs": {Type: "string", Description: "JSON tracking specs for conversion tracking."},
			},
			Required: []string{"account", "name", "adset_id", "creative", "status"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			AdAccountID   string `json:"ad_account_id"`
			Name          string `json:"name"`
			AdsetID       string `json:"adset_id"`
			Creative      string `json:"creative"`
			Status        string `json:"status"`
			TrackingSpecs string `json:"tracking_specs"`
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
		fp := url.Values{}
		fp.Set("name", p.Name)
		fp.Set("adset_id", p.AdsetID)
		fp.Set("creative", p.Creative)
		fp.Set("status", p.Status)
		setOptional(fp, "tracking_specs", p.TrackingSpecs)
		return client.Post(ctx, fmt.Sprintf("/%s/ads", acctID), fp)
	})

	// ── get_ad ─────────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad",
		Description: "Get details of a specific ad by ID.",
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		qp := url.Values{}
		if p.Fields != "" {
			qp.Set("fields", p.Fields)
		} else {
			qp.Set("fields", "id,name,adset_id,campaign_id,status,effective_status,creative{id,name,body,title,thumbnail_url,object_story_spec},tracking_specs,created_time,updated_time")
		}
		return client.Get(ctx, fmt.Sprintf("/%s", p.AdID), qp)
	})

	// ── update_ad ──────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "update_ad",
		Description: "Update an existing ad. Only provided fields are changed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"ad_id":          {Type: "string", Description: "Ad ID to update."},
				"name":           {Type: "string", Description: "New ad name."},
				"status":         {Type: "string", Description: "ACTIVE, PAUSED, or ARCHIVED."},
				"creative":       {Type: "string", Description: "JSON creative spec to replace the current creative."},
				"tracking_specs": {Type: "string", Description: "JSON tracking specs."},
				"adset_id":       {Type: "string", Description: "Move ad to a different adset."},
			},
			Required: []string{"account", "ad_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			AdID          string `json:"ad_id"`
			Name          string `json:"name"`
			Status        string `json:"status"`
			Creative      string `json:"creative"`
			TrackingSpecs string `json:"tracking_specs"`
			AdsetID       string `json:"adset_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		setOptional(fp, "name", p.Name)
		setOptional(fp, "status", p.Status)
		setOptional(fp, "creative", p.Creative)
		setOptional(fp, "tracking_specs", p.TrackingSpecs)
		setOptional(fp, "adset_id", p.AdsetID)
		return client.Post(ctx, fmt.Sprintf("/%s", p.AdID), fp)
	})

	// ── delete_ad ──────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "delete_ad",
		Description: "Delete an ad by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"ad_id":   {Type: "string", Description: "Ad ID to delete."},
			},
			Required: []string{"account", "ad_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			AdID    string `json:"ad_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/%s", p.AdID))
	})

	// ── get_ad_insights ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_insights",
		Description: "Get performance insights/metrics for a specific ad.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"ad_id":       {Type: "string", Description: "Ad ID."},
				"fields":      {Type: "string", Description: "Comma-separated metrics."},
				"date_preset": {Type: "string", Description: "Date preset: today, yesterday, last_7d, last_30d, etc."},
				"time_range":  {Type: "string", Description: "JSON {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}."},
				"breakdowns":  {Type: "string", Description: "Comma-separated breakdowns."},
			},
			Required: []string{"account", "ad_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AdID       string `json:"ad_id"`
			Fields     string `json:"fields"`
			DatePreset string `json:"date_preset"`
			TimeRange  string `json:"time_range"`
			Breakdowns string `json:"breakdowns"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		qp := insightParams(p.Fields, p.DatePreset, p.TimeRange, p.Breakdowns, "", "", "")
		return client.Get(ctx, fmt.Sprintf("/%s/insights", p.AdID), qp)
	})

	// ── get_ad_previews ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_previews",
		Description: "Get previews for an existing ad in various formats. Returns HTML for rendering.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"ad_id":     {Type: "string", Description: "Ad ID."},
				"ad_format": {Type: "string", Description: "Preview format: DESKTOP_FEED_STANDARD, MOBILE_FEED_STANDARD, INSTAGRAM_STANDARD, RIGHT_COLUMN_STANDARD."},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		qp := url.Values{}
		qp.Set("ad_format", p.AdFormat)
		return client.Get(ctx, fmt.Sprintf("/%s/previews", p.AdID), qp)
	})

	// ── copy_ad ────────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "copy_ad",
		Description: "Duplicate an ad, optionally to a different adset.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"ad_id":          {Type: "string", Description: "Ad ID to copy."},
				"adset_id":       {Type: "string", Description: "Target adset ID (optional, defaults to same adset)."},
				"status_option":  {Type: "string", Description: "Status for the copy: PAUSED or INHERITED (default: PAUSED)."},
				"rename_options": {Type: "string", Description: "JSON rename options: {\"rename_strategy\":\"DEEP_RENAME\",\"rename_prefix\":\"Copy of\"}."},
			},
			Required: []string{"account", "ad_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			AdID          string `json:"ad_id"`
			AdsetID       string `json:"adset_id"`
			StatusOption  string `json:"status_option"`
			RenameOptions string `json:"rename_options"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		if p.StatusOption != "" {
			fp.Set("status_option", p.StatusOption)
		} else {
			fp.Set("status_option", "PAUSED")
		}
		setOptional(fp, "adset_id", p.AdsetID)
		setOptional(fp, "rename_options", p.RenameOptions)
		return client.Post(ctx, fmt.Sprintf("/%s/copies", p.AdID), fp)
	})
}
