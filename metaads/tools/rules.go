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

// RegisterRules registers automated ad rule tools.
func RegisterRules(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_ad_rules",
		Description: "List automated ad rules for the ad account via GET /{ad_account_id}/adrules_library.",
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
			v.Set("fields", "id,name,status,evaluation_spec,execution_spec,schedule_spec")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/adrules_library", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_ad_rule",
		Description: "Create an automated ad rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"name":           {Type: "string", Description: "Rule name."},
				"evaluation_spec": {Type: "string", Description: "JSON evaluation spec (trigger conditions)."},
				"execution_spec": {Type: "string", Description: "JSON execution spec (actions to take)."},
				"schedule_spec":  {Type: "string", Description: "JSON schedule spec (when to evaluate)."},
			},
			Required: []string{"account", "name", "evaluation_spec", "execution_spec"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			Name           string `json:"name"`
			EvaluationSpec string `json:"evaluation_spec"`
			ExecutionSpec  string `json:"execution_spec"`
			ScheduleSpec   string `json:"schedule_spec"`
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
		v.Set("evaluation_spec", p.EvaluationSpec)
		v.Set("execution_spec", p.ExecutionSpec)
		if p.ScheduleSpec != "" {
			v.Set("schedule_spec", p.ScheduleSpec)
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/adrules_library", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_rule",
		Description: "Get details of a specific ad rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"rule_id": {Type: "string", Description: "Ad rule ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			RuleID  string `json:"rule_id"`
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
		}
		return cl.Get(ctx, "/"+p.RuleID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_ad_rule",
		Description: "Update an existing ad rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"rule_id":        {Type: "string", Description: "Ad rule ID."},
				"name":           {Type: "string", Description: "New rule name."},
				"status":         {Type: "string", Description: "Rule status (ENABLED, DISABLED)."},
				"evaluation_spec": {Type: "string", Description: "JSON evaluation spec."},
				"execution_spec": {Type: "string", Description: "JSON execution spec."},
				"schedule_spec":  {Type: "string", Description: "JSON schedule spec."},
			},
			Required: []string{"account", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			RuleID         string `json:"rule_id"`
			Name           string `json:"name"`
			Status         string `json:"status"`
			EvaluationSpec string `json:"evaluation_spec"`
			ExecutionSpec  string `json:"execution_spec"`
			ScheduleSpec   string `json:"schedule_spec"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		if p.Name != "" {
			v.Set("name", p.Name)
		}
		if p.Status != "" {
			v.Set("status", p.Status)
		}
		if p.EvaluationSpec != "" {
			v.Set("evaluation_spec", p.EvaluationSpec)
		}
		if p.ExecutionSpec != "" {
			v.Set("execution_spec", p.ExecutionSpec)
		}
		if p.ScheduleSpec != "" {
			v.Set("schedule_spec", p.ScheduleSpec)
		}
		return cl.Post(ctx, "/"+p.RuleID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_ad_rule",
		Description: "Delete an ad rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"rule_id": {Type: "string", Description: "Ad rule ID to delete."},
			},
			Required: []string{"account", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			RuleID  string `json:"rule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.RuleID)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_rule_execution_history",
		Description: "Get execution history for an ad rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"rule_id": {Type: "string", Description: "Ad rule ID."},
			}, paging.Properties()),
			Required: []string{"account", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			RuleID  string `json:"rule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := cl.FetchAll(ctx, "/"+p.RuleID+"/history", "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})
}
