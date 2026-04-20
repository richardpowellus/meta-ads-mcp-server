package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/meta-ads-mcp-server/internal/paging"
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterInsights registers analytics/reporting insight tools.
func RegisterInsights(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── get_account_insights ───────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_account_insights",
		Description: "Get performance insights for the entire ad account. Supports field selection, date presets, time ranges, breakdowns, and sorting.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated metrics: impressions,clicks,spend,cpc,cpm,ctr,reach,frequency,actions,conversions,cost_per_action_type,purchase_roas."},
				"date_preset":   {Type: "string", Description: "Date preset: today, yesterday, last_7d, last_14d, last_30d, this_month, last_month, this_quarter, last_3d, maximum."},
				"time_range":    {Type: "string", Description: "JSON object: {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}. Mutually exclusive with date_preset."},
				"breakdowns":    {Type: "string", Description: "Comma-separated breakdowns: age, gender, country, region, dma, impression_device, placement, platform_position, device_platform, publisher_platform."},
				"level":         {Type: "string", Description: "Aggregation level: account, campaign, adset, ad (default: account)."},
				"filtering":     {Type: "string", Description: "JSON array of filter objects: [{\"field\":\"campaign.name\",\"operator\":\"CONTAIN\",\"value\":\"Brand\"}]."},
				"sort":          {Type: "string", Description: "Sort field and direction: spend_descending, impressions_ascending, etc."},
				"limit":         {Type: "string", Description: "Max rows to return (default: 25)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
			DatePreset  string `json:"date_preset"`
			TimeRange   string `json:"time_range"`
			Breakdowns  string `json:"breakdowns"`
			Level       string `json:"level"`
			Filtering   string `json:"filtering"`
			Sort        string `json:"sort"`
			Limit       string `json:"limit"`
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
		qp := insightParams(p.Fields, p.DatePreset, p.TimeRange, p.Breakdowns, p.Level, p.Filtering, p.Sort)
		if p.Limit != "" {
			qp.Set("limit", p.Limit)
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/insights", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── get_campaign_insights (insights module version) ────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_insights_report",
		Description: "Get insights across all campaigns (or a filtered subset) at the campaign level. Supports breakdowns and sorting for reporting.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated metrics."},
				"date_preset":   {Type: "string", Description: "Date preset."},
				"time_range":    {Type: "string", Description: "JSON {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}."},
				"breakdowns":    {Type: "string", Description: "Comma-separated breakdowns."},
				"filtering":     {Type: "string", Description: "JSON array of filter objects."},
				"sort":          {Type: "string", Description: "Sort field and direction."},
				"limit":         {Type: "string", Description: "Max rows to return."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
			DatePreset  string `json:"date_preset"`
			TimeRange   string `json:"time_range"`
			Breakdowns  string `json:"breakdowns"`
			Filtering   string `json:"filtering"`
			Sort        string `json:"sort"`
			Limit       string `json:"limit"`
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
		qp := insightParams(p.Fields, p.DatePreset, p.TimeRange, p.Breakdowns, "campaign", p.Filtering, p.Sort)
		if p.Limit != "" {
			qp.Set("limit", p.Limit)
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/insights", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── get_adset_insights_report ──────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_adset_insights_report",
		Description: "Get insights across all ad sets at the adset level. Supports breakdowns, filtering, and sorting for reporting.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated metrics."},
				"date_preset":   {Type: "string", Description: "Date preset."},
				"time_range":    {Type: "string", Description: "JSON {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}."},
				"breakdowns":    {Type: "string", Description: "Comma-separated breakdowns."},
				"filtering":     {Type: "string", Description: "JSON array of filter objects."},
				"sort":          {Type: "string", Description: "Sort field and direction."},
				"limit":         {Type: "string", Description: "Max rows to return."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
			DatePreset  string `json:"date_preset"`
			TimeRange   string `json:"time_range"`
			Breakdowns  string `json:"breakdowns"`
			Filtering   string `json:"filtering"`
			Sort        string `json:"sort"`
			Limit       string `json:"limit"`
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
		qp := insightParams(p.Fields, p.DatePreset, p.TimeRange, p.Breakdowns, "adset", p.Filtering, p.Sort)
		if p.Limit != "" {
			qp.Set("limit", p.Limit)
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/insights", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── get_ad_insights_report ─────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_insights_report",
		Description: "Get insights across all ads at the ad level. Supports breakdowns, filtering, and sorting for reporting.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated metrics."},
				"date_preset":   {Type: "string", Description: "Date preset."},
				"time_range":    {Type: "string", Description: "JSON {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}."},
				"breakdowns":    {Type: "string", Description: "Comma-separated breakdowns."},
				"filtering":     {Type: "string", Description: "JSON array of filter objects."},
				"sort":          {Type: "string", Description: "Sort field and direction."},
				"limit":         {Type: "string", Description: "Max rows to return."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
			DatePreset  string `json:"date_preset"`
			TimeRange   string `json:"time_range"`
			Breakdowns  string `json:"breakdowns"`
			Filtering   string `json:"filtering"`
			Sort        string `json:"sort"`
			Limit       string `json:"limit"`
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
		qp := insightParams(p.Fields, p.DatePreset, p.TimeRange, p.Breakdowns, "ad", p.Filtering, p.Sort)
		if p.Limit != "" {
			qp.Set("limit", p.Limit)
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/insights", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})
}
