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

// RegisterAudiences registers custom audience and lookalike tools.
func RegisterAudiences(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_custom_audiences",
		Description: "List custom audiences in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields (default: id,name,subtype,approximate_count,delivery_status,operation_status)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Fields  string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/customaudiences"
		vals := url.Values{}
		fields := p.Fields
		if fields == "" {
			fields = "id,name,subtype,approximate_count,delivery_status,operation_status"
		}
		vals.Set("fields", fields)
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_custom_audience",
		Description: "Create a custom audience in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"name":        {Type: "string", Description: "Audience name."},
				"subtype":     {Type: "string", Description: "Subtype: CUSTOM, WEBSITE, APP, OFFLINE_CONVERSION, CLAIM, ENGAGEMENT, etc."},
				"description": {Type: "string", Description: "Optional description."},
				"rule":        {Type: "string", Description: "Optional JSON rule definition for website/app audiences."},
				"customer_file_source": {Type: "string", Description: "Source for customer file audiences: USER_PROVIDED_ONLY, PARTNER_PROVIDED_ONLY, BOTH_USER_AND_PARTNER_PROVIDED."},
			},
			Required: []string{"account", "name", "subtype"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account            string `json:"account"`
			Name               string `json:"name"`
			Subtype            string `json:"subtype"`
			Description        string `json:"description"`
			Rule               string `json:"rule"`
			CustomerFileSource string `json:"customer_file_source"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/customaudiences"
		vals := url.Values{}
		vals.Set("name", p.Name)
		vals.Set("subtype", p.Subtype)
		if p.Description != "" {
			vals.Set("description", p.Description)
		}
		if p.Rule != "" {
			vals.Set("rule", p.Rule)
		}
		if p.CustomerFileSource != "" {
			vals.Set("customer_file_source", p.CustomerFileSource)
		}
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_custom_audience",
		Description: "Get details of a specific custom audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
				"fields":      {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
			Fields     string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.AudienceID
		vals := url.Values{}
		if p.Fields != "" {
			vals.Set("fields", p.Fields)
		}
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_custom_audience",
		Description: "Update an existing custom audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
				"name":        {Type: "string", Description: "New name."},
				"description": {Type: "string", Description: "New description."},
				"rule":        {Type: "string", Description: "Updated JSON rule definition."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AudienceID  string `json:"audience_id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Rule        string `json:"rule"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		vals := url.Values{}
		if p.Name != "" {
			vals.Set("name", p.Name)
		}
		if p.Description != "" {
			vals.Set("description", p.Description)
		}
		if p.Rule != "" {
			vals.Set("rule", p.Rule)
		}
		return client.Post(ctx, "/"+p.AudienceID, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_custom_audience",
		Description: "Delete a custom audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID to delete."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, "/"+p.AudienceID)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "add_audience_users",
		Description: "Add users to a custom audience. Provide a JSON payload with schema and data arrays containing hashed user identifiers.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
				"payload":     {Type: "string", Description: "JSON payload: {\"schema\":[\"EMAIL\"],\"data\":[[\"hash1\"],[\"hash2\"]]}. Values must be SHA-256 hashed."},
			},
			Required: []string{"account", "audience_id", "payload"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
			Payload    string `json:"payload"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.AudienceID + "/users"
		vals := url.Values{}
		vals.Set("payload", p.Payload)
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "remove_audience_users",
		Description: "Remove users from a custom audience. Provide a JSON payload with schema and data arrays.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
				"payload":     {Type: "string", Description: "JSON payload: {\"schema\":[\"EMAIL\"],\"data\":[[\"hash1\"],[\"hash2\"]]}. Values must be SHA-256 hashed."},
			},
			Required: []string{"account", "audience_id", "payload"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
			Payload    string `json:"payload"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.AudienceID + "/users"
		vals := url.Values{}
		vals.Set("payload", p.Payload)
		// DELETE with payload — Graph API supports this as a POST with method=delete
		vals.Set("method", "delete")
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_lookalike_audience",
		Description: "Create a lookalike audience based on an existing custom audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":            {Type: "string", Description: "Account name."},
				"name":               {Type: "string", Description: "Name for the lookalike audience."},
				"origin_audience_id": {Type: "string", Description: "Source custom audience ID."},
				"ratio":              {Type: "string", Description: "Lookalike ratio (0.01 to 0.20). E.g. 0.01 = top 1%."},
				"country":            {Type: "string", Description: "Two-letter country code for the lookalike (e.g. US)."},
			},
			Required: []string{"account", "name", "origin_audience_id", "ratio", "country"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			Name             string `json:"name"`
			OriginAudienceID string `json:"origin_audience_id"`
			Ratio            string `json:"ratio"`
			Country          string `json:"country"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/customaudiences"
		vals := url.Values{}
		vals.Set("name", p.Name)
		vals.Set("subtype", "LOOKALIKE")
		vals.Set("origin_audience_id", p.OriginAudienceID)
		vals.Set("lookalike_spec", fmt.Sprintf(
			`{"ratio":%s,"country":"%s"}`,
			p.Ratio, p.Country,
		))
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_health",
		Description: "Get health metrics for a custom audience (size, operation status).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		vals := url.Values{}
		vals.Set("fields", "id,name,approximate_count,operation_status,delivery_status,permission_for_actions,data_source")
		return client.Get(ctx, "/"+p.AudienceID, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "share_custom_audience",
		Description: "Share a custom audience with another ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"audience_id":      {Type: "string", Description: "The custom audience ID to share."},
				"target_account_id": {Type: "string", Description: "The ad account ID to share with (with or without act_ prefix)."},
			},
			Required: []string{"account", "audience_id", "target_account_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account         string `json:"account"`
			AudienceID      string `json:"audience_id"`
			TargetAccountID string `json:"target_account_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.AudienceID + "/adaccounts"
		vals := url.Values{}
		vals.Set("adaccounts", fmt.Sprintf(`["%s"]`, p.TargetAccountID))
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_sharing_status",
		Description: "Get the sharing status of a custom audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The custom audience ID."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.AudienceID + "/shared_account_info"
		vals := url.Values{}
		vals.Set("fields", "account_id,account_name,business_name,sharing_status")
		return client.Get(ctx, path, vals)
	})
}
