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

// RegisterConversions registers custom conversion and offline event tools.
func RegisterConversions(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_custom_conversions",
		Description: "List custom conversions for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
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
			v.Set("fields", "id,name,custom_event_type,rule,pixel,default_conversion_value,is_archived")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/customconversions", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_custom_conversion",
		Description: "Create a custom conversion.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":            {Type: "string", Description: "Account name."},
				"name":               {Type: "string", Description: "Custom conversion name."},
				"pixel_id":           {Type: "string", Description: "Pixel ID to associate with."},
				"custom_event_type":  {Type: "string", Description: "Event type (e.g. PURCHASE, LEAD, OTHER)."},
				"rule":               {Type: "string", Description: "JSON rule object for URL/event matching."},
				"default_conversion_value": {Type: "number", Description: "Default conversion value."},
			},
			Required: []string{"account", "name", "custom_event_type", "rule"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account                string  `json:"account"`
			Name                   string  `json:"name"`
			PixelID                string  `json:"pixel_id"`
			CustomEventType        string  `json:"custom_event_type"`
			Rule                   string  `json:"rule"`
			DefaultConversionValue float64 `json:"default_conversion_value"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("name", p.Name)
		v.Set("custom_event_type", p.CustomEventType)
		v.Set("rule", p.Rule)
		if p.PixelID != "" {
			v.Set("pixel_id", p.PixelID)
		}
		if p.DefaultConversionValue > 0 {
			v.Set("default_conversion_value", fmt.Sprintf("%.2f", p.DefaultConversionValue))
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/customconversions", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_custom_conversion",
		Description: "Get details of a specific custom conversion.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"conversion_id": {Type: "string", Description: "Custom conversion ID."},
				"fields":        {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "conversion_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ConversionID string `json:"conversion_id"`
			Fields       string `json:"fields"`
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
		return cl.Get(ctx, "/"+p.ConversionID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_custom_conversion",
		Description: "Delete a custom conversion.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"conversion_id": {Type: "string", Description: "Custom conversion ID to delete."},
			},
			Required: []string{"account", "conversion_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ConversionID string `json:"conversion_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.ConversionID)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_offline_event_sets",
		Description: "List offline event sets for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
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
			v.Set("fields", "id,name,description,upload_tag,event_stats,is_auto_assigned")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/offline_conversion_data_sets", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_offline_event_set",
		Description: "Create an offline event set for tracking offline conversions.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"name":        {Type: "string", Description: "Event set name."},
				"description": {Type: "string", Description: "Event set description."},
				"upload_tag":  {Type: "string", Description: "Upload tag for grouping events."},
			},
			Required: []string{"account", "name"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			Name        string `json:"name"`
			Description string `json:"description"`
			UploadTag   string `json:"upload_tag"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("name", p.Name)
		if p.Description != "" {
			v.Set("description", p.Description)
		}
		if p.UploadTag != "" {
			v.Set("upload_tag", p.UploadTag)
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/offline_conversion_data_sets", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "upload_offline_events",
		Description: "Upload offline conversion events to an event set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"event_set_id": {Type: "string", Description: "Offline event set ID."},
				"upload_tag":   {Type: "string", Description: "Upload tag for this batch."},
				"data":         {Type: "string", Description: "JSON array of offline event objects."},
			},
			Required: []string{"account", "event_set_id", "data"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			EventSetID string `json:"event_set_id"`
			UploadTag  string `json:"upload_tag"`
			Data       string `json:"data"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("data", p.Data)
		if p.UploadTag != "" {
			v.Set("upload_tag", p.UploadTag)
		}
		return cl.Post(ctx, "/"+p.EventSetID+"/events", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_offline_event_set_stats",
		Description: "Get statistics for an offline event set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"event_set_id": {Type: "string", Description: "Offline event set ID."},
				"fields":       {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "event_set_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			EventSetID string `json:"event_set_id"`
			Fields     string `json:"fields"`
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
			v.Set("fields", "id,name,event_stats,duplicate_entries,matched_entries,valid_entries")
		}
		return cl.Get(ctx, "/"+p.EventSetID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "assign_offline_event_set_accounts",
		Description: "Assign ad accounts to an offline event set for attribution.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"event_set_id":   {Type: "string", Description: "Offline event set ID."},
				"ad_account_ids": {Type: "string", Description: "Comma-separated ad account IDs to assign."},
			},
			Required: []string{"account", "event_set_id", "ad_account_ids"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			EventSetID   string `json:"event_set_id"`
			AdAccountIDs string `json:"ad_account_ids"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("ad_accounts", p.AdAccountIDs)
		return cl.Post(ctx, "/"+p.EventSetID+"/adaccounts", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "send_conversion_events",
		Description: "Send server-side conversion events via the Conversions API (CAPI). POST /{pixel_id}/events.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"pixel_id": {Type: "string", Description: "Pixel ID for the Conversions API."},
				"data":     {Type: "string", Description: "JSON array of server event objects (event_name, event_time, user_data, etc.)."},
				"test_event_code": {Type: "string", Description: "Test event code for validation (omit for production)."},
			},
			Required: []string{"account", "pixel_id", "data"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			PixelID       string `json:"pixel_id"`
			Data          string `json:"data"`
			TestEventCode string `json:"test_event_code"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("data", p.Data)
		if p.TestEventCode != "" {
			v.Set("test_event_code", p.TestEventCode)
		}
		return cl.Post(ctx, "/"+p.PixelID+"/events", v)
	})
}
