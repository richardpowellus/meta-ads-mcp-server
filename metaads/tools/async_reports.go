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

// RegisterAsyncReports registers async report tools.
func RegisterAsyncReports(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "create_async_report",
		Description: "Create an async insights report for the ad account. Returns a report_run_id to poll.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"level":        {Type: "string", Description: "Reporting level: account, campaign, adset, or ad."},
				"fields":       {Type: "string", Description: "Comma-separated metrics (e.g. impressions,spend,clicks,cpc,ctr)."},
				"time_range":   {Type: "string", Description: "JSON time range object with since/until (YYYY-MM-DD)."},
				"filtering":    {Type: "string", Description: "JSON array of filter objects."},
				"breakdowns":   {Type: "string", Description: "Comma-separated breakdowns (e.g. age,gender,country)."},
				"time_increment": {Type: "string", Description: "Time granularity: 1 (daily), 7 (weekly), monthly, all_days."},
			},
			Required: []string{"account", "fields"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			Level         string `json:"level"`
			Fields        string `json:"fields"`
			TimeRange     string `json:"time_range"`
			Filtering     string `json:"filtering"`
			Breakdowns    string `json:"breakdowns"`
			TimeIncrement string `json:"time_increment"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", p.Fields)
		if p.Level != "" {
			v.Set("level", p.Level)
		}
		if p.TimeRange != "" {
			v.Set("time_range", p.TimeRange)
		}
		if p.Filtering != "" {
			v.Set("filtering", p.Filtering)
		}
		if p.Breakdowns != "" {
			v.Set("breakdowns", p.Breakdowns)
		}
		if p.TimeIncrement != "" {
			v.Set("time_increment", p.TimeIncrement)
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/insights", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_async_report_status",
		Description: "Check the status of an async insights report run.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"report_run_id": {Type: "string", Description: "Report run ID from create_async_report."},
			},
			Required: []string{"account", "report_run_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			ReportRunID string `json:"report_run_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "id,async_status,async_percent_completion,date_start,date_stop")
		return cl.Get(ctx, "/"+p.ReportRunID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_async_report_results",
		Description: "Get results of a completed async insights report.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"report_run_id": {Type: "string", Description: "Report run ID."},
			}, paging.Properties()),
			Required: []string{"account", "report_run_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			ReportRunID string `json:"report_run_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := cl.FetchAll(ctx, "/"+p.ReportRunID+"/insights", "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_report_runs",
		Description: "List recent async report runs for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "id,async_status,async_percent_completion,date_start,date_stop,time_completed")
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/insights", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})
}
