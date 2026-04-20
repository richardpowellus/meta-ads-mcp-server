package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterAdLibrary registers Ad Library search and reporting tools.
func RegisterAdLibrary(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "search_ad_library",
		Description: "Search the Meta Ad Library for ads matching given criteria. Searches GET /ads_archive with search_terms and optional filters. Returns ad creative, spend, impressions, and funding entity.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"search_terms":     {Type: "string", Description: "Keywords to search for in ad text/creative."},
				"ad_type":          {Type: "string", Description: "Type of ads to search.", Enum: []string{"ALL", "POLITICAL_AND_ISSUE_ADS", "HOUSING_ADS", "CREDIT_ADS", "EMPLOYMENT_ADS"}},
				"ad_reached_countries": {Type: "string", Description: "Comma-separated ISO country codes (e.g. US,GB)."},
				"search_page_ids":  {Type: "string", Description: "Comma-separated page IDs to filter by."},
				"ad_active_status": {Type: "string", Description: "Filter by ad active status.", Enum: []string{"ALL", "ACTIVE", "INACTIVE"}},
				"publisher_platforms": {Type: "string", Description: "Comma-separated platforms: facebook,instagram,audience_network,messenger,whatsapp."},
				"media_type":       {Type: "string", Description: "Filter by media type.", Enum: []string{"ALL", "IMAGE", "MEME", "VIDEO", "NONE"}},
				"limit":            {Type: "integer", Description: "Max results to return (default 25, max 1000)."},
				"fields":           {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "search_terms", "ad_reached_countries"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account            string `json:"account"`
			SearchTerms        string `json:"search_terms"`
			AdType             string `json:"ad_type"`
			AdReachedCountries string `json:"ad_reached_countries"`
			SearchPageIDs      string `json:"search_page_ids"`
			AdActiveStatus     string `json:"ad_active_status"`
			PublisherPlatforms string `json:"publisher_platforms"`
			MediaType          string `json:"media_type"`
			Limit              int    `json:"limit"`
			Fields             string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		q.Set("search_terms", p.SearchTerms)
		q.Set("ad_reached_countries", fmt.Sprintf("[%q]", p.AdReachedCountries))

		if p.AdType != "" {
			q.Set("ad_type", p.AdType)
		}
		if p.SearchPageIDs != "" {
			q.Set("search_page_ids", p.SearchPageIDs)
		}
		if p.AdActiveStatus != "" {
			q.Set("ad_active_status", p.AdActiveStatus)
		}
		if p.PublisherPlatforms != "" {
			q.Set("publisher_platforms", p.PublisherPlatforms)
		}
		if p.MediaType != "" {
			q.Set("media_type", p.MediaType)
		}
		if p.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", p.Limit))
		}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}

		return client.Get(ctx, "/ads_archive", q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_library_report",
		Description: "Get an Ad Library report with aggregated spend and impression data for a page or country. Returns summary statistics from the Ad Library Report API.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":             {Type: "string", Description: "Account name."},
				"report_type":         {Type: "string", Description: "Report type.", Enum: []string{"page", "country", "region", "demographic"}},
				"ad_reached_countries": {Type: "string", Description: "Comma-separated ISO country codes."},
				"search_page_ids":     {Type: "string", Description: "Comma-separated page IDs for page-level report."},
				"time_range":          {Type: "string", Description: "JSON object with since/until dates, e.g. {\"since\":\"2024-01-01\",\"until\":\"2024-12-31\"}."},
				"fields":              {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account            string `json:"account"`
			ReportType         string `json:"report_type"`
			AdReachedCountries string `json:"ad_reached_countries"`
			SearchPageIDs      string `json:"search_page_ids"`
			TimeRange          string `json:"time_range"`
			Fields             string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.AdReachedCountries != "" {
			q.Set("ad_reached_countries", fmt.Sprintf("[%q]", p.AdReachedCountries))
		}
		if p.SearchPageIDs != "" {
			q.Set("search_page_ids", p.SearchPageIDs)
		}
		if p.TimeRange != "" {
			q.Set("time_range", p.TimeRange)
		}
		if p.Fields != "" {
			q.Set("fields", p.Fields)
		}

		return client.Get(ctx, "/ads_archive_report", q)
	})
}
