package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterCanvas registers Instant Experience (Canvas) tools.
func RegisterCanvas(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_canvases",
		Description: "List Instant Experiences (Canvases) for the ad account. Returns canvas ID, name, status, and URL.",
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

		path := fmt.Sprintf("/%s/canvases", client.AdAccountID())
		return client.Get(ctx, path, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_canvas",
		Description: "Create a new Instant Experience (Canvas) for the ad account. Provide a body_elements JSON array defining the canvas layout.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"name":          {Type: "string", Description: "Canvas name."},
				"body_elements": {Type: "string", Description: "JSON array of canvas body elements (buttons, photos, videos, carousels, etc.)."},
				"is_published":  {Type: "string", Description: "Set to 'true' to publish immediately."},
			},
			Required: []string{"account", "name", "body_elements"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			Name         string `json:"name"`
			BodyElements string `json:"body_elements"`
			IsPublished  string `json:"is_published"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		q.Set("name", p.Name)
		q.Set("body_element_ids", p.BodyElements)
		if p.IsPublished != "" {
			q.Set("is_published", p.IsPublished)
		}

		path := fmt.Sprintf("/%s/canvases", client.AdAccountID())
		return client.Post(ctx, path, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_canvas",
		Description: "Get details of a specific Instant Experience (Canvas) by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"canvas_id": {Type: "string", Description: "The Canvas/Instant Experience ID."},
				"fields":    {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "canvas_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			CanvasID string `json:"canvas_id"`
			Fields   string `json:"fields"`
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

		return client.Get(ctx, "/"+p.CanvasID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_canvas",
		Description: "Update an existing Instant Experience (Canvas).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"canvas_id":     {Type: "string", Description: "The Canvas/Instant Experience ID."},
				"name":          {Type: "string", Description: "Updated canvas name."},
				"body_elements": {Type: "string", Description: "Updated JSON array of canvas body elements."},
				"is_published":  {Type: "string", Description: "Set to 'true' to publish."},
			},
			Required: []string{"account", "canvas_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			CanvasID     string `json:"canvas_id"`
			Name         string `json:"name"`
			BodyElements string `json:"body_elements"`
			IsPublished  string `json:"is_published"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		if p.Name != "" {
			q.Set("name", p.Name)
		}
		if p.BodyElements != "" {
			q.Set("body_element_ids", p.BodyElements)
		}
		if p.IsPublished != "" {
			q.Set("is_published", p.IsPublished)
		}

		return client.Post(ctx, "/"+p.CanvasID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_canvas",
		Description: "Delete an Instant Experience (Canvas) by ID.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"canvas_id": {Type: "string", Description: "The Canvas/Instant Experience ID."},
			},
			Required: []string{"account", "canvas_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			CanvasID string `json:"canvas_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		return client.Delete(ctx, "/"+p.CanvasID)
	})
}
