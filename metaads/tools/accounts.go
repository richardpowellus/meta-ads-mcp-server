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

// RegisterAccounts registers ad account management tools.
func RegisterAccounts(s mcp.ToolRegistrar, cfg *metaads.Config) {

	// ── list_accounts ──────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_accounts",
		Description: "List all configured Meta Ads accounts. Returns account names (never tokens).",
		InputSchema: mcp.InputSchema{Type: "object"},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		return cfg.Accounts(ctx), nil
	})

	// ── get_ad_account ─────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_account",
		Description: "Get details of a Meta ad account including name, currency, timezone, status, and spend cap.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"fields":        {Type: "string", Description: "Comma-separated fields to return (e.g. name,currency,timezone_name,account_status,spend_cap,amount_spent,balance)."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Fields      string `json:"fields"`
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
		qp := url.Values{}
		if p.Fields != "" {
			qp.Set("fields", p.Fields)
		} else {
			qp.Set("fields", "id,name,account_id,currency,timezone_name,account_status,spend_cap,amount_spent,balance,business,owner")
		}
		return client.Get(ctx, fmt.Sprintf("/%s", acctID), qp)
	})

	// ── update_ad_account ──────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "update_ad_account",
		Description: "Update ad account settings such as name or spend cap.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"name":          {Type: "string", Description: "New account name."},
				"spend_cap":     {Type: "string", Description: "Spend cap in ACCOUNT CURRENCY (dollars for USD accounts), NOT cents. Meta's API multiplies this value by 100 to store cents. Example: '1000' for $1,000."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Name        string `json:"name"`
			SpendCap    string `json:"spend_cap"`
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
		if p.Name != "" {
			fp.Set("name", p.Name)
		}
		if p.SpendCap != "" {
			fp.Set("spend_cap", p.SpendCap)
		}
		return client.Post(ctx, fmt.Sprintf("/%s", acctID), fp)
	})

	// ── list_ad_account_users ──────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_ad_account_users",
		Description: "List users who have access to the ad account with their roles and permissions.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/users", acctID), "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── get_user_ad_accounts ───────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_user_ad_accounts",
		Description: "List ad accounts accessible by a user (defaults to 'me').",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"user_id": {Type: "string", Description: "User ID or 'me' (default: me)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			UserID  string `json:"user_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		uid := p.UserID
		if uid == "" {
			uid = "me"
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/adaccounts", uid), "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── set_account_spend_cap ──────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "set_account_spend_cap",
		Description: "Set the account-level spend cap. Value is in ACCOUNT CURRENCY (dollars for USD accounts) — NOT cents, despite what earlier versions of this description claimed. Meta's Marketing API multiplies this value by 100 when storing. Example: pass '1000' to set a $1,000 cap (stored as 100000 cents). Set to '0' to remove the cap. Note: GET responses (get_ad_account) return spend_cap in cents — read cents, write dollars.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"spend_cap":     {Type: "string", Description: "Spend cap amount in ACCOUNT CURRENCY (dollars for USD accounts), NOT cents. Meta's API multiplies this value by 100 to store cents. Example: '1000' for $1,000 cap. '0' removes the cap."},
			},
			Required: []string{"account", "spend_cap"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			SpendCap    string `json:"spend_cap"`
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
		fp.Set("spend_cap", p.SpendCap)
		return client.Post(ctx, fmt.Sprintf("/%s", acctID), fp)
	})

	// ── list_ad_labels ─────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "list_ad_labels",
		Description: "List all ad labels in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/%s/adlabels", acctID), "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// ── create_ad_label ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "create_ad_label",
		Description: "Create a new ad label in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"name":          {Type: "string", Description: "Label name."},
			},
			Required: []string{"account", "name"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Name        string `json:"name"`
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
		return client.Post(ctx, fmt.Sprintf("/%s/adlabels", acctID), fp)
	})

	// ── get_ad_label ───────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "get_ad_label",
		Description: "Get details of a specific ad label by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"label_id": {Type: "string", Description: "Ad label ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "label_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			LabelID string `json:"label_id"`
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
		}
		return client.Get(ctx, fmt.Sprintf("/%s", p.LabelID), qp)
	})

	// ── delete_ad_label ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "delete_ad_label",
		Description: "Delete an ad label by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"label_id": {Type: "string", Description: "Ad label ID to delete."},
			},
			Required: []string{"account", "label_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			LabelID string `json:"label_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/%s", p.LabelID))
	})

	// ── assign_ad_label ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "assign_ad_label",
		Description: "Assign ad labels to campaigns, adsets, or ads. Provide the object ID and label IDs.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"object_id": {Type: "string", Description: "ID of the campaign, adset, or ad to label."},
				"label_ids": {Type: "string", Description: "JSON array of label IDs, e.g. [\"123\",\"456\"]."},
			},
			Required: []string{"account", "object_id", "label_ids"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ObjectID string `json:"object_id"`
			LabelIDs string `json:"label_ids"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		fp.Set("adlabels", p.LabelIDs)
		return client.Post(ctx, fmt.Sprintf("/%s/adlabels", p.ObjectID), fp)
	})

	// ── remove_ad_label ────────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "remove_ad_label",
		Description: "Remove ad labels from a campaign, adset, or ad.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"object_id": {Type: "string", Description: "ID of the campaign, adset, or ad."},
				"label_ids": {Type: "string", Description: "JSON array of label IDs to remove, e.g. [\"123\"]."},
			},
			Required: []string{"account", "object_id", "label_ids"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ObjectID string `json:"object_id"`
			LabelIDs string `json:"label_ids"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		fp := url.Values{}
		fp.Set("adlabels", p.LabelIDs)
		return client.Post(ctx, fmt.Sprintf("/%s/adlabels", p.ObjectID), fp)
	})

	// ── generate_ad_preview ────────────────────────────────────────────
	s.RegisterTool(mcp.Tool{
		Name:        "generate_ad_preview",
		Description: "Generate an ad preview for a given creative spec or existing ad. Returns HTML for rendering.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"ad_account_id": {Type: "string", Description: "Ad account ID (overrides default)."},
				"creative":      {Type: "string", Description: "JSON creative spec (object_story_spec, etc.)."},
				"ad_format":     {Type: "string", Description: "Preview format: DESKTOP_FEED_STANDARD, MOBILE_FEED_STANDARD, RIGHT_COLUMN_STANDARD, etc."},
			},
			Required: []string{"account", "ad_format"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			AdAccountID string `json:"ad_account_id"`
			Creative    string `json:"creative"`
			AdFormat    string `json:"ad_format"`
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
		qp := url.Values{}
		qp.Set("ad_format", p.AdFormat)
		if p.Creative != "" {
			qp.Set("creative", p.Creative)
		}
		return client.Get(ctx, fmt.Sprintf("/%s/generatepreviews", acctID), qp)
	})
}
