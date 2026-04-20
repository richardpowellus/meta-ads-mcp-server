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

// RegisterTargeting registers audience targeting and research tools.
func RegisterTargeting(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "search_interests",
		Description: "Search for interest-based targeting options by keyword.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"query":   {Type: "string", Description: "Search query for interests."},
				"limit":   {Type: "string", Description: "Max results (default 50)."},
			}, paging.Properties()),
			Required: []string{"account", "query"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Query   string `json:"query"`
			Limit   string `json:"limit"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		vals := url.Values{}
		vals.Set("type", "adinterest")
		vals.Set("q", p.Query)
		if p.Limit != "" {
			vals.Set("limit", p.Limit)
		}
		items, err := client.FetchAll(ctx, "/search", "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "search_behaviors",
		Description: "Search for behavior-based targeting categories.",
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
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		vals := url.Values{}
		vals.Set("type", "adTargetingCategory")
		vals.Set("class", "behaviors")
		items, err := client.FetchAll(ctx, "/search", "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "search_demographics",
		Description: "Search for demographic targeting categories.",
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
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		vals := url.Values{}
		vals.Set("type", "adTargetingCategory")
		vals.Set("class", "demographics")
		items, err := client.FetchAll(ctx, "/search", "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "search_locations",
		Description: "Search for geographic targeting locations by keyword.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"query":          {Type: "string", Description: "Search query for locations (e.g. city/country name)."},
				"location_types": {Type: "string", Description: "Comma-separated types: country, region, city, zip, geo_market, electoral_district (default: all)."},
			}, paging.Properties()),
			Required: []string{"account", "query"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			Query         string `json:"query"`
			LocationTypes string `json:"location_types"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		vals := url.Values{}
		vals.Set("type", "adgeolocation")
		vals.Set("q", p.Query)
		if p.LocationTypes != "" {
			vals.Set("location_types", fmt.Sprintf(`["%s"]`, p.LocationTypes))
		}
		items, err := client.FetchAll(ctx, "/search", "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_targeting_categories",
		Description: "Browse available targeting categories for the ad account.",
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
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/targetingbrowse"
		items, err := client.FetchAll(ctx, path, "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_reach_estimate",
		Description: "Get a reach estimate for the given targeting specification.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"targeting_spec": {Type: "string", Description: "JSON targeting spec (e.g. {\"geo_locations\":{\"countries\":[\"US\"]},\"interests\":[{\"id\":\"6003139266461\"}]})."},
			},
			Required: []string{"account", "targeting_spec"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			TargetingSpec string `json:"targeting_spec"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/reachestimate"
		vals := url.Values{}
		vals.Set("targeting_spec", p.TargetingSpec)
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_delivery_estimate",
		Description: "Get delivery estimates for a targeting specification including daily outcomes curve.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"targeting_spec": {Type: "string", Description: "JSON targeting spec."},
				"optimization_goal": {Type: "string", Description: "Optimization goal (e.g. LINK_CLICKS, IMPRESSIONS, REACH)."},
			},
			Required: []string{"account", "targeting_spec", "optimization_goal"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account            string `json:"account"`
			TargetingSpec      string `json:"targeting_spec"`
			OptimizationGoal   string `json:"optimization_goal"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/delivery_estimate"
		vals := url.Values{}
		vals.Set("targeting_spec", p.TargetingSpec)
		vals.Set("optimization_goal", p.OptimizationGoal)
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_targeting_browse",
		Description: "Browse detailed targeting options organized by category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"limit_type":   {Type: "string", Description: "Filter by type (e.g. interests, behaviors, demographics)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			LimitType string `json:"limit_type"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/targetingbrowse"
		vals := url.Values{}
		if p.LimitType != "" {
			vals.Set("limit_type", p.LimitType)
		}
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_broad_targeting_categories",
		Description: "Get broad targeting categories for the ad account.",
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
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/broadtargetingcategories"
		items, err := client.FetchAll(ctx, path, "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "search_targeting_options",
		Description: "General-purpose targeting search. Specify the targeting class to search within.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"query":   {Type: "string", Description: "Search query."},
				"class":   {Type: "string", Description: "Targeting class: interests, behaviors, demographics, life_events, industries, income, family_statuses, user_device, user_os."},
			}, paging.Properties()),
			Required: []string{"account", "query", "class"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Query   string `json:"query"`
			Class   string `json:"class"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/targetingsearch"
		vals := url.Values{}
		vals.Set("q", p.Query)
		vals.Set("targeting_list", fmt.Sprintf(`[{"type":"%s"}]`, p.Class))
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})
}
