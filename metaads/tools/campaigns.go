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

// RegisterCampaigns registers campaign management tools.
func RegisterCampaigns(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── list_campaigns ─────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_campaigns",
		Description: "List all campaigns in the ad account. Supports filtering by status and date.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated fields (e.g. id,name,status,objective,daily_budget,lifetime_budget,start_time,stop_time)."},
				"status":        {Type: "string", Description: "Filter by effective_status: ACTIVE, PAUSED, ARCHIVED, etc. Comma-separated for multiple."},
				"date_preset":   {Type: "string", Description: "Date preset for filtering: today, yesterday, last_7d, last_30d, etc."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
			Status      string `json:"status"`
			DatePreset  string `json:"date_preset"`
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
			qp.Set("fields", "id,name,objective,status,effective_status,daily_budget,lifetime_budget,budget_remaining,start_time,stop_time,created_time,updated_time")
		}
		if p.Status != "" {
			qp.Set("effective_status", fmt.Sprintf("[%q]", p.Status))
		}
		if p.DatePreset != "" {
			qp.Set("date_preset", p.DatePreset)
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/campaigns", acctID), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── create_campaign ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "create_campaign",
		Description: "Create a new campaign. Requires name, objective, and status at minimum.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":                       {Type: "string", Description: "Account name."},
				"ad_account_id":                  {Type: "string", Description: "Ad account ID (overrides default)."},
				"name":                           {Type: "string", Description: "Campaign name."},
				"objective":                      {Type: "string", Description: "Campaign objective: OUTCOME_AWARENESS, OUTCOME_ENGAGEMENT, OUTCOME_LEADS, OUTCOME_SALES, OUTCOME_TRAFFIC, OUTCOME_APP_PROMOTION."},
				"status":                         {Type: "string", Description: "Initial status: ACTIVE or PAUSED."},
				"special_ad_categories":          {Type: "string", Description: "JSON array of special ad categories: NONE, EMPLOYMENT, HOUSING, CREDIT, ISSUES_ELECTIONS_POLITICS."},
				"daily_budget":                   {Type: "string", Description: "Daily budget in cents (account currency)."},
				"lifetime_budget":                {Type: "string", Description: "Lifetime budget in cents (account currency)."},
				"spend_cap":                      {Type: "string", Description: "Campaign spend cap in cents."},
				"bid_strategy":                   {Type: "string", Description: "Bid strategy: LOWEST_COST_WITHOUT_CAP, LOWEST_COST_WITH_BID_CAP, COST_CAP."},
				"buying_type":                    {Type: "string", Description: "AUCTION (default) or RESERVED."},
				"start_time":                     {Type: "string", Description: "Start time in ISO 8601 format."},
				"stop_time":                      {Type: "string", Description: "Stop time in ISO 8601 format."},
				"is_adset_budget_sharing_enabled": {Type: "string", Description: "Whether Advantage Campaign Budget is enabled: true or false."},
			},
			Required: []string{"account", "name", "objective", "status", "special_ad_categories"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account                    string `json:"account"`
			AdAccountID                string `json:"ad_account_id"`
			Name                       string `json:"name"`
			Objective                  string `json:"objective"`
			Status                     string `json:"status"`
			SpecialAdCategories        string `json:"special_ad_categories"`
			DailyBudget                string `json:"daily_budget"`
			LifetimeBudget             string `json:"lifetime_budget"`
			SpendCap                   string `json:"spend_cap"`
			BidStrategy                string `json:"bid_strategy"`
			BuyingType                 string `json:"buying_type"`
			StartTime                  string `json:"start_time"`
			StopTime                   string `json:"stop_time"`
			IsAdsetBudgetSharingEnabled string `json:"is_adset_budget_sharing_enabled"`
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
		fp.Set("objective", p.Objective)
		fp.Set("status", p.Status)
		fp.Set("special_ad_categories", p.SpecialAdCategories)
		setOptional(fp, "daily_budget", p.DailyBudget)
		setOptional(fp, "lifetime_budget", p.LifetimeBudget)
		setOptional(fp, "spend_cap", p.SpendCap)
		setOptional(fp, "bid_strategy", p.BidStrategy)
		setOptional(fp, "buying_type", p.BuyingType)
		setOptional(fp, "start_time", p.StartTime)
		setOptional(fp, "stop_time", p.StopTime)
		setOptional(fp, "is_adset_budget_sharing_enabled", p.IsAdsetBudgetSharingEnabled)
		return client.Post(ctx, fmt.Sprintf("/%s/campaigns", acctID), fp)
	})

	// ── get_campaign ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign",
		Description: "Get details of a specific campaign by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID."},
				"fields":      {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
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
			qp.Set("fields", "id,name,objective,status,effective_status,daily_budget,lifetime_budget,budget_remaining,bid_strategy,start_time,stop_time,created_time,updated_time,special_ad_categories")
		}
		return client.Get(ctx, fmt.Sprintf("/%s", p.CampaignID), qp)
	})

	// ── update_campaign ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "update_campaign",
		Description: "Update an existing campaign. Only provided fields are changed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"campaign_id":    {Type: "string", Description: "Campaign ID to update."},
				"name":           {Type: "string", Description: "New campaign name."},
				"status":         {Type: "string", Description: "ACTIVE, PAUSED, or ARCHIVED."},
				"daily_budget":   {Type: "string", Description: "Daily budget in cents."},
				"lifetime_budget": {Type: "string", Description: "Lifetime budget in cents."},
				"spend_cap":      {Type: "string", Description: "Campaign spend cap in cents."},
				"bid_strategy":   {Type: "string", Description: "Bid strategy."},
				"start_time":     {Type: "string", Description: "Start time in ISO 8601."},
				"stop_time":      {Type: "string", Description: "Stop time in ISO 8601."},
			},
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			CampaignID     string `json:"campaign_id"`
			Name           string `json:"name"`
			Status         string `json:"status"`
			DailyBudget    string `json:"daily_budget"`
			LifetimeBudget string `json:"lifetime_budget"`
			SpendCap       string `json:"spend_cap"`
			BidStrategy    string `json:"bid_strategy"`
			StartTime      string `json:"start_time"`
			StopTime       string `json:"stop_time"`
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
		setOptional(fp, "daily_budget", p.DailyBudget)
		setOptional(fp, "lifetime_budget", p.LifetimeBudget)
		setOptional(fp, "spend_cap", p.SpendCap)
		setOptional(fp, "bid_strategy", p.BidStrategy)
		setOptional(fp, "start_time", p.StartTime)
		setOptional(fp, "stop_time", p.StopTime)
		return client.Post(ctx, fmt.Sprintf("/%s", p.CampaignID), fp)
	})

	// ── delete_campaign ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "delete_campaign",
		Description: "Delete a campaign by ID. Sets status to DELETED.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID to delete."},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/%s", p.CampaignID))
	})

	// ── get_campaign_insights ──────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_insights",
		Description: "Get performance insights/metrics for a specific campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID."},
				"fields":      {Type: "string", Description: "Comma-separated metrics: impressions,clicks,spend,cpc,cpm,ctr,reach,frequency,actions,conversions,cost_per_action_type."},
				"date_preset": {Type: "string", Description: "Date preset: today, yesterday, last_7d, last_30d, this_month, last_month, maximum."},
				"time_range":  {Type: "string", Description: "JSON object with since/until dates: {\"since\":\"2024-01-01\",\"until\":\"2024-01-31\"}."},
				"breakdowns":  {Type: "string", Description: "Comma-separated breakdowns: age, gender, country, placement, device_platform."},
			},
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
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
		return client.Get(ctx, fmt.Sprintf("/%s/insights", p.CampaignID), qp)
	})

	// ── duplicate_campaign ─────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "duplicate_campaign",
		Description: "Duplicate an existing campaign including its adsets and ads.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID to duplicate."},
				"deep_copy":   {Type: "string", Description: "Whether to deep-copy adsets and ads: true or false (default: true)."},
				"status_option": {Type: "string", Description: "Status for the new campaign: PAUSED or INHERITED (default: PAUSED)."},
			},
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			CampaignID   string `json:"campaign_id"`
			DeepCopy     string `json:"deep_copy"`
			StatusOption string `json:"status_option"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		if p.DeepCopy != "" {
			fp.Set("deep_copy", p.DeepCopy)
		} else {
			fp.Set("deep_copy", "true")
		}
		if p.StatusOption != "" {
			fp.Set("status_option", p.StatusOption)
		} else {
			fp.Set("status_option", "PAUSED")
		}
		return client.Post(ctx, fmt.Sprintf("/%s/copies", p.CampaignID), fp)
	})

	// ── copy_campaign ──────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "copy_campaign",
		Description: "Copy a campaign to a different ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":              {Type: "string", Description: "Account name."},
				"campaign_id":          {Type: "string", Description: "Source campaign ID."},
				"target_ad_account_id": {Type: "string", Description: "Target ad account ID (with act_ prefix)."},
				"deep_copy":            {Type: "string", Description: "Whether to deep-copy adsets and ads: true or false (default: true)."},
				"status_option":        {Type: "string", Description: "Status for the copy: PAUSED or INHERITED (default: PAUSED)."},
			},
			Required: []string{"account", "campaign_id", "target_ad_account_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account           string `json:"account"`
			CampaignID        string `json:"campaign_id"`
			TargetAdAccountID string `json:"target_ad_account_id"`
			DeepCopy          string `json:"deep_copy"`
			StatusOption      string `json:"status_option"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		if p.DeepCopy != "" {
			fp.Set("deep_copy", p.DeepCopy)
		} else {
			fp.Set("deep_copy", "true")
		}
		if p.StatusOption != "" {
			fp.Set("status_option", p.StatusOption)
		} else {
			fp.Set("status_option", "PAUSED")
		}
		fp.Set("target_ad_account_id", p.TargetAdAccountID)
		return client.Post(ctx, fmt.Sprintf("/%s/copies", p.CampaignID), fp)
	})
}
