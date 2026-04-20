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

// RegisterLeads registers lead generation tools.
func RegisterLeads(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_lead_forms",
		Description: "List lead generation forms for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return (e.g. id,name,status,leads_count)."},
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
			v.Set("fields", "id,name,status,leads_count,created_time")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/leadgen_forms", cl.AdAccountID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_lead_form",
		Description: "Create a new lead generation form.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"name":             {Type: "string", Description: "Form name."},
				"questions":        {Type: "string", Description: "JSON array of question objects."},
				"privacy_policy":   {Type: "string", Description: "JSON privacy policy object with url and link_text."},
				"follow_up_action": {Type: "string", Description: "Follow-up action URL or message."},
			},
			Required: []string{"account", "name", "questions", "privacy_policy"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			Name           string `json:"name"`
			Questions      string `json:"questions"`
			PrivacyPolicy  string `json:"privacy_policy"`
			FollowUpAction string `json:"follow_up_action"`
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
		v.Set("questions", p.Questions)
		v.Set("privacy_policy", p.PrivacyPolicy)
		if p.FollowUpAction != "" {
			v.Set("follow_up_action_url", p.FollowUpAction)
		}
		return cl.Post(ctx, fmt.Sprintf("/%s/leadgen_forms", cl.AdAccountID()), v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_lead_form",
		Description: "Get details of a specific lead generation form.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"form_id": {Type: "string", Description: "Lead form ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "form_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			FormID  string `json:"form_id"`
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
		return cl.Get(ctx, "/"+p.FormID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_leads",
		Description: "Get leads submitted to a lead generation form.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"form_id": {Type: "string", Description: "Lead form ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields (e.g. id,created_time,field_data)."},
			}, paging.Properties()),
			Required: []string{"account", "form_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			FormID  string `json:"form_id"`
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
			v.Set("fields", "id,created_time,field_data")
		}
		items, err := cl.FetchAll(ctx, "/"+p.FormID+"/leads", "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_lead",
		Description: "Get a single lead by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"lead_id": {Type: "string", Description: "Lead ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "lead_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			LeadID  string `json:"lead_id"`
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
		return cl.Get(ctx, "/"+p.LeadID, v)
	})
}
