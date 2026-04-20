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

// RegisterBilling registers billing and spend tools.
func RegisterBilling(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "get_billing_transactions",
		Description: "Get billing-related activity log entries for the ad account. Filters to billing events such as funding, billing, invoice, and payment events.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"since":   {Type: "string", Description: "Start timestamp (ISO 8601 or Unix epoch)."},
				"until":   {Type: "string", Description: "End timestamp (ISO 8601 or Unix epoch)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Since   string `json:"since"`
			Until   string `json:"until"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/activities"
		vals := url.Values{}
		vals.Set("fields", "event_type,event_time,extra_data,object_id,object_name")
		if p.Since != "" {
			vals.Set("since", p.Since)
		}
		if p.Until != "" {
			vals.Set("until", p.Until)
		}
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}

		// Filter to billing-related events
		billingKeywords := []string{"billing", "fund", "invoice", "payment", "charge", "credit", "spend"}
		var filtered []json.RawMessage
		for _, item := range items {
			var activity struct {
				EventType string `json:"event_type"`
			}
			_ = json.Unmarshal(item, &activity)
			for _, kw := range billingKeywords {
				if containsCI(activity.EventType, kw) {
					filtered = append(filtered, item)
					break
				}
			}
		}
		if filtered == nil {
			filtered = items // fall back to all activities if no billing-specific ones found
		}
		return paging.Emit(filtered, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_account_spend",
		Description: "Get account-level spend insights (spend, impressions, clicks) for a date range.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"date_start": {Type: "string", Description: "Start date (YYYY-MM-DD). Defaults to last 7 days if omitted."},
				"date_end":   {Type: "string", Description: "End date (YYYY-MM-DD). Defaults to today if omitted."},
				"fields":     {Type: "string", Description: "Comma-separated fields (default: spend,impressions,clicks,cpc,cpm,ctr,reach,frequency)."},
				"level":      {Type: "string", Description: "Aggregation level: account (default), campaign, adset, ad.", Enum: []string{"account", "campaign", "adset", "ad"}},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			DateStart string `json:"date_start"`
			DateEnd   string `json:"date_end"`
			Fields    string `json:"fields"`
			Level     string `json:"level"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/insights"
		vals := url.Values{}
		fields := p.Fields
		if fields == "" {
			fields = "spend,impressions,clicks,cpc,cpm,ctr,reach,frequency"
		}
		vals.Set("fields", fields)
		if p.DateStart != "" {
			vals.Set("time_range", fmt.Sprintf(`{"since":"%s","until":"%s"}`, p.DateStart, p.DateEnd))
		}
		if p.Level != "" {
			vals.Set("level", p.Level)
		}
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_spend_by_day",
		Description: "Get daily spend breakdown for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"date_start": {Type: "string", Description: "Start date (YYYY-MM-DD)."},
				"date_end":   {Type: "string", Description: "End date (YYYY-MM-DD)."},
			}, paging.Properties()),
			Required: []string{"account", "date_start", "date_end"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			DateStart string `json:"date_start"`
			DateEnd   string `json:"date_end"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/insights"
		vals := url.Values{}
		vals.Set("fields", "spend,impressions,clicks,date_start,date_stop")
		vals.Set("time_range", fmt.Sprintf(`{"since":"%s","until":"%s"}`, p.DateStart, p.DateEnd))
		vals.Set("time_increment", "1")
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})
}

// containsCI checks if s contains substr (case-insensitive).
func containsCI(s, substr string) bool {
	sl := len(substr)
	if sl == 0 {
		return true
	}
	if len(s) < sl {
		return false
	}
	for i := 0; i <= len(s)-sl; i++ {
		match := true
		for j := 0; j < sl; j++ {
			sc, tc := s[i+j], substr[j]
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if tc >= 'A' && tc <= 'Z' {
				tc += 32
			}
			if sc != tc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
