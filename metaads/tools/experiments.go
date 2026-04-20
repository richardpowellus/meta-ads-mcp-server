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

// RegisterExperiments registers A/B testing and brand lift study tools.
func RegisterExperiments(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_ad_studies",
		Description: "List ad studies (experiments) for the business.",
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
			v.Set("fields", "id,name,description,type,start_time,end_time,results")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/ad_studies", cl.BusinessID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_ad_study",
		Description: "Create a new ad study (A/B test or holdout experiment).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"name":        {Type: "string", Description: "Study name."},
				"description": {Type: "string", Description: "Study description."},
				"type":        {Type: "string", Description: "Study type (e.g. SPLIT_TEST, HOLDOUT)."},
				"start_time":  {Type: "string", Description: "Start time (ISO 8601)."},
				"end_time":    {Type: "string", Description: "End time (ISO 8601)."},
				"cells":       {Type: "string", Description: "JSON array of study cell objects."},
				"objectives":  {Type: "string", Description: "JSON array of study objective objects."},
			},
			Required: []string{"account", "name", "type"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Type        string `json:"type"`
			StartTime   string `json:"start_time"`
			EndTime     string `json:"end_time"`
			Cells       string `json:"cells"`
			Objectives  string `json:"objectives"`
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
		v.Set("type", p.Type)
		if p.Description != "" {
			v.Set("description", p.Description)
		}
		if p.StartTime != "" {
			v.Set("start_time", p.StartTime)
		}
		if p.EndTime != "" {
			v.Set("end_time", p.EndTime)
		}
		if p.Cells != "" {
			v.Set("cells", p.Cells)
		}
		if p.Objectives != "" {
			v.Set("objectives", p.Objectives)
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/ad_studies", cl.BusinessID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_study",
		Description: "Get details of a specific ad study.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"study_id": {Type: "string", Description: "Ad study ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "study_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StudyID string `json:"study_id"`
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
		return cl.Get(ctx, "/"+p.StudyID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_study_results",
		Description: "Get the results of an ad study.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"study_id": {Type: "string", Description: "Ad study ID."},
			},
			Required: []string{"account", "study_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StudyID string `json:"study_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("fields", "id,name,type,results,confidence_level,winner_cell")
		return cl.Get(ctx, "/"+p.StudyID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_brand_lift_results",
		Description: "Get brand lift study results including survey data.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"study_id": {Type: "string", Description: "Brand lift study ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "study_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StudyID string `json:"study_id"`
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
			v.Set("fields", "id,name,type,results,questions,start_time,end_time")
		}
		return cl.Get(ctx, "/"+p.StudyID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_ad_study",
		Description: "Delete an ad study.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"study_id": {Type: "string", Description: "Ad study ID to delete."},
			},
			Required: []string{"account", "study_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StudyID string `json:"study_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.StudyID)
	})
}
