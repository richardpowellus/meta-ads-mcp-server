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

// RegisterAdsets registers ad set management tools.
func RegisterAdsets(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── list_adsets ────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_adsets",
		Description: "List all ad sets in the ad account, optionally filtered by campaign or status.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"campaign_id":   {Type: "string", Description: "Filter by campaign ID. If set, lists adsets under this campaign instead of the account."},
				"fields":        {Type: "string", Description: "Comma-separated fields."},
				"status":        {Type: "string", Description: "Filter by effective_status: ACTIVE, PAUSED, etc."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
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
			qp.Set("fields", "id,name,campaign_id,status,effective_status,daily_budget,lifetime_budget,budget_remaining,optimization_goal,billing_event,bid_amount,start_time,end_time,targeting,promoted_object,created_time,updated_time")
		}
		if p.Status != "" {
			qp.Set("effective_status", fmt.Sprintf("[%q]", p.Status))
		}
		parent := p.CampaignID
		if parent == "" {
			parent = p.AdAccountID
			if parent == "" {
				parent = client.AdAccountID()
			}
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/adsets", parent), "data", qp)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── create_adset ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "create_adset",
		Description: "Create a new ad set. Requires campaign_id, name, optimization_goal, billing_event, and status.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":           {Type: "string", Description: "Account name."},
				"ad_account_id":     {Type: "string", Description: "Ad account ID (overrides default)."},
				"campaign_id":       {Type: "string", Description: "Parent campaign ID."},
				"name":              {Type: "string", Description: "Ad set name."},
				"optimization_goal": {Type: "string", Description: "Optimization goal: LINK_CLICKS, IMPRESSIONS, REACH, CONVERSIONS, LANDING_PAGE_VIEWS, LEAD_GENERATION, etc."},
				"billing_event":     {Type: "string", Description: "Billing event: IMPRESSIONS, LINK_CLICKS, etc."},
				"status":            {Type: "string", Description: "Initial status: ACTIVE or PAUSED."},
				"daily_budget":      {Type: "string", Description: "Daily budget in cents (account currency)."},
				"lifetime_budget":   {Type: "string", Description: "Lifetime budget in cents."},
				"bid_amount":        {Type: "string", Description: "Bid amount in cents (for manual bidding)."},
				"bid_strategy":      {Type: "string", Description: "Bid strategy: LOWEST_COST_WITHOUT_CAP, LOWEST_COST_WITH_BID_CAP, COST_CAP."},
				"start_time":        {Type: "string", Description: "Start time in ISO 8601."},
				"end_time":          {Type: "string", Description: "End time in ISO 8601 (required for lifetime_budget)."},
				"targeting":         {Type: "string", Description: "JSON targeting spec: {\"geo_locations\":{\"countries\":[\"US\"]},\"age_min\":18,\"age_max\":65}."},
				"promoted_object":   {Type: "string", Description: "JSON promoted object: {\"pixel_id\":\"123\",\"custom_event_type\":\"PURCHASE\"} or {\"page_id\":\"456\"}."},
				"destination_type":  {Type: "string", Description: "Destination type: WEBSITE, APP, etc."},
			},
			Required: []string{"account", "campaign_id", "name", "optimization_goal", "billing_event", "status"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			AdAccountID      string `json:"ad_account_id"`
			CampaignID       string `json:"campaign_id"`
			Name             string `json:"name"`
			OptimizationGoal string `json:"optimization_goal"`
			BillingEvent     string `json:"billing_event"`
			Status           string `json:"status"`
			DailyBudget      string `json:"daily_budget"`
			LifetimeBudget   string `json:"lifetime_budget"`
			BidAmount        string `json:"bid_amount"`
			BidStrategy      string `json:"bid_strategy"`
			StartTime        string `json:"start_time"`
			EndTime          string `json:"end_time"`
			Targeting        string `json:"targeting"`
			PromotedObject   string `json:"promoted_object"`
			DestinationType  string `json:"destination_type"`
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
		fp.Set("campaign_id", p.CampaignID)
		fp.Set("name", p.Name)
		fp.Set("optimization_goal", p.OptimizationGoal)
		fp.Set("billing_event", p.BillingEvent)
		fp.Set("status", p.Status)
		setOptional(fp, "daily_budget", p.DailyBudget)
		setOptional(fp, "lifetime_budget", p.LifetimeBudget)
		setOptional(fp, "bid_amount", p.BidAmount)
		setOptional(fp, "bid_strategy", p.BidStrategy)
		setOptional(fp, "start_time", p.StartTime)
		setOptional(fp, "end_time", p.EndTime)
		setOptional(fp, "targeting", p.Targeting)
		setOptional(fp, "promoted_object", p.PromotedObject)
		setOptional(fp, "destination_type", p.DestinationType)
		return client.Post(ctx, fmt.Sprintf("/%s/adsets", acctID), fp)
	})

	// ── get_adset ──────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_adset",
		Description: "Get details of a specific ad set by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"adset_id": {Type: "string", Description: "Ad set ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "adset_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			AdsetID string `json:"adset_id"`
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
			qp.Set("fields", "id,name,campaign_id,status,effective_status,daily_budget,lifetime_budget,budget_remaining,optimization_goal,billing_event,bid_amount,bid_strategy,start_time,end_time,targeting,promoted_object,created_time,updated_time")
		}
		return client.Get(ctx, fmt.Sprintf("/%s", p.AdsetID), qp)
	})

	// ── update_adset ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "update_adset",
		Description: "Update an existing ad set. Only provided fields are changed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":           {Type: "string", Description: "Account name."},
				"adset_id":          {Type: "string", Description: "Ad set ID to update."},
				"name":              {Type: "string", Description: "New ad set name."},
				"status":            {Type: "string", Description: "ACTIVE, PAUSED, or ARCHIVED."},
				"daily_budget":      {Type: "string", Description: "Daily budget in cents."},
				"lifetime_budget":   {Type: "string", Description: "Lifetime budget in cents."},
				"bid_amount":        {Type: "string", Description: "Bid amount in cents."},
				"bid_strategy":      {Type: "string", Description: "Bid strategy."},
				"optimization_goal": {Type: "string", Description: "Optimization goal."},
				"billing_event":     {Type: "string", Description: "Billing event."},
				"start_time":        {Type: "string", Description: "Start time in ISO 8601."},
				"end_time":          {Type: "string", Description: "End time in ISO 8601."},
				"targeting":         {Type: "string", Description: "JSON targeting spec."},
				"promoted_object":   {Type: "string", Description: "JSON promoted object."},
			},
			Required: []string{"account", "adset_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			AdsetID          string `json:"adset_id"`
			Name             string `json:"name"`
			Status           string `json:"status"`
			DailyBudget      string `json:"daily_budget"`
			LifetimeBudget   string `json:"lifetime_budget"`
			BidAmount        string `json:"bid_amount"`
			BidStrategy      string `json:"bid_strategy"`
			OptimizationGoal string `json:"optimization_goal"`
			BillingEvent     string `json:"billing_event"`
			StartTime        string `json:"start_time"`
			EndTime          string `json:"end_time"`
			Targeting        string `json:"targeting"`
			PromotedObject   string `json:"promoted_object"`
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
		setOptional(fp, "bid_amount", p.BidAmount)
		setOptional(fp, "bid_strategy", p.BidStrategy)
		setOptional(fp, "optimization_goal", p.OptimizationGoal)
		setOptional(fp, "billing_event", p.BillingEvent)
		setOptional(fp, "start_time", p.StartTime)
		setOptional(fp, "end_time", p.EndTime)
		setOptional(fp, "targeting", p.Targeting)
		setOptional(fp, "promoted_object", p.PromotedObject)
		return client.Post(ctx, fmt.Sprintf("/%s", p.AdsetID), fp)
	})

	// ── delete_adset ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "delete_adset",
		Description: "Delete an ad set by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"adset_id": {Type: "string", Description: "Ad set ID to delete."},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/%s", p.AdsetID))
	})

	// ── get_adset_insights ─────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_adset_insights",
		Description: "Get performance insights/metrics for a specific ad set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"adset_id":    {Type: "string", Description: "Ad set ID."},
				"fields":      {Type: "string", Description: "Comma-separated metrics."},
				"date_preset": {Type: "string", Description: "Date preset: today, yesterday, last_7d, last_30d, etc."},
				"time_range":  {Type: "string", Description: "JSON {\"since\":\"YYYY-MM-DD\",\"until\":\"YYYY-MM-DD\"}."},
				"breakdowns":  {Type: "string", Description: "Comma-separated breakdowns."},
			},
			Required: []string{"account", "adset_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AdsetID    string `json:"adset_id"`
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
		return client.Get(ctx, fmt.Sprintf("/%s/insights", p.AdsetID), qp)
	})

	// ── copy_adset ─────────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "copy_adset",
		Description: "Duplicate an ad set, optionally to a different campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"adset_id":      {Type: "string", Description: "Ad set ID to copy."},
				"campaign_id":   {Type: "string", Description: "Target campaign ID (optional, defaults to same campaign)."},
				"deep_copy":     {Type: "string", Description: "Whether to deep-copy ads: true or false (default: true)."},
				"status_option": {Type: "string", Description: "Status for the copy: PAUSED or INHERITED (default: PAUSED)."},
				"rename_options": {Type: "string", Description: "JSON rename options: {\"rename_strategy\":\"DEEP_RENAME\",\"rename_prefix\":\"Copy of\",\"rename_suffix\":\"\"}."},
			},
			Required: []string{"account", "adset_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			AdsetID       string `json:"adset_id"`
			CampaignID    string `json:"campaign_id"`
			DeepCopy      string `json:"deep_copy"`
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
		setOptional(fp, "campaign_id", p.CampaignID)
		setOptional(fp, "rename_options", p.RenameOptions)
		return client.Post(ctx, fmt.Sprintf("/%s/copies", p.AdsetID), fp)
	})

	// ── get_adset_targeting_sentence ────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_adset_targeting_sentence",
		Description: "Get a human-readable description of the ad set's targeting configuration.",
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/%s/targetingsentencelines", p.AdsetID), nil)
	})
}
