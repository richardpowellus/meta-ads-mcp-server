package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterBudgetPlanning registers budget schedule, reach & frequency prediction, and minimum budget tools.
func RegisterBudgetPlanning(s mcp.ToolRegistrar, cfg *metaads.Config) {
	// --- Budget Schedules ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_budget_schedules",
		Description: "List budget schedules for a campaign. Budget schedules define how budget is distributed over time.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "Campaign ID to list budget schedules for."},
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
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}

		return client.Get(ctx, "/"+p.CampaignID+"/budget_schedules", q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_budget_schedule",
		Description: "Create a budget schedule for a campaign to define how budget is distributed over time periods.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"campaign_id":   {Type: "string", Description: "Campaign ID."},
				"budget_value":  {Type: "string", Description: "Budget value in account currency (in cents for USD)."},
				"time_start":    {Type: "string", Description: "Start time (ISO 8601)."},
				"time_end":      {Type: "string", Description: "End time (ISO 8601)."},
				"budget_value_type": {Type: "string", Description: "Budget type.", Enum: []string{"ABSOLUTE", "MULTIPLIER"}},
			},
			Required: []string{"account", "campaign_id", "budget_value", "time_start", "time_end"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account         string `json:"account"`
			CampaignID      string `json:"campaign_id"`
			BudgetValue     string `json:"budget_value"`
			TimeStart       string `json:"time_start"`
			TimeEnd         string `json:"time_end"`
			BudgetValueType string `json:"budget_value_type"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		q.Set("budget_value", p.BudgetValue)
		q.Set("time_start", p.TimeStart)
		q.Set("time_end", p.TimeEnd)
		if p.BudgetValueType != "" {
			q.Set("budget_value_type", p.BudgetValueType)
		}

		return client.Post(ctx, "/"+p.CampaignID+"/budget_schedules", q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_budget_schedule",
		Description: "Update an existing budget schedule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":            {Type: "string", Description: "Account name."},
				"budget_schedule_id": {Type: "string", Description: "Budget schedule ID."},
				"budget_value":       {Type: "string", Description: "Updated budget value."},
				"time_start":         {Type: "string", Description: "Updated start time (ISO 8601)."},
				"time_end":           {Type: "string", Description: "Updated end time (ISO 8601)."},
			},
			Required: []string{"account", "budget_schedule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			BudgetScheduleID string `json:"budget_schedule_id"`
			BudgetValue      string `json:"budget_value"`
			TimeStart        string `json:"time_start"`
			TimeEnd          string `json:"time_end"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.BudgetValue != "" {
			q.Set("budget_value", p.BudgetValue)
		}
		if p.TimeStart != "" {
			q.Set("time_start", p.TimeStart)
		}
		if p.TimeEnd != "" {
			q.Set("time_end", p.TimeEnd)
		}

		return client.Post(ctx, "/"+p.BudgetScheduleID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_budget_schedule",
		Description: "Delete a budget schedule by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":            {Type: "string", Description: "Account name."},
				"budget_schedule_id": {Type: "string", Description: "Budget schedule ID to delete."},
			},
			Required: []string{"account", "budget_schedule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			BudgetScheduleID string `json:"budget_schedule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		return client.Delete(ctx, "/"+p.BudgetScheduleID)
	})

	// --- Reach & Frequency Predictions ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_rf_predictions",
		Description: "List Reach & Frequency predictions for the ad account. These forecast ad delivery reach, frequency, and costs.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
				"limit":   {Type: "integer", Description: "Max results per page."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Fields  string `json:"fields"`
			Limit   int    `json:"limit"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}
		if p.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", p.Limit))
		}

		path := fmt.Sprintf("/%s/reachfrequencypredictions", client.AdAccountID())
		return client.Get(ctx, path, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_rf_prediction",
		Description: "Create a new Reach & Frequency prediction to forecast campaign delivery. Requires targeting spec, budget, and schedule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"prediction_body": {Type: "string", Description: "JSON object with prediction parameters: target_spec, start_time, stop_time, budget, frequency_cap, destination_id, objective, etc."},
			},
			Required: []string{"account", "prediction_body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string          `json:"account"`
			PredictionBody json.RawMessage `json:"prediction_body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf("/%s/reachfrequencypredictions", client.AdAccountID())
		return client.PostJSON(ctx, path, p.PredictionBody)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_rf_prediction",
		Description: "Get details of a specific Reach & Frequency prediction by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"prediction_id": {Type: "string", Description: "Reach & Frequency prediction ID."},
				"fields":        {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "prediction_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			PredictionID string `json:"prediction_id"`
			Fields       string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}

		return client.Get(ctx, "/"+p.PredictionID, q)
	})

	// --- Minimum Budgets ---

	s.RegisterTool(mcp.Tool{
		Name:        "get_minimum_budgets",
		Description: "Get minimum budget requirements for a campaign or ad set. Returns minimum daily and lifetime budgets based on targeting and bid strategy.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"campaign_id":   {Type: "string", Description: "Campaign ID to get minimums for (provide this or ad_set_id)."},
				"ad_set_id":     {Type: "string", Description: "Ad set ID to get minimums for (provide this or campaign_id)."},
				"fields":        {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			AdSetID    string `json:"ad_set_id"`
			Fields     string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		var parentID string
		switch {
		case p.AdSetID != "":
			parentID = p.AdSetID
		case p.CampaignID != "":
			parentID = p.CampaignID
		default:
			parentID = client.AdAccountID()
		}

		q := url.Values{}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}

		return client.Get(ctx, "/"+parentID+"/minimum_budgets", q)
	})
}
