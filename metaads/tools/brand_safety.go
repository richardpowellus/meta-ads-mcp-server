package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterBrandSafety registers publisher block list management tools for brand safety.
func RegisterBrandSafety(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_block_lists",
		Description: "List publisher block lists for the ad account. Block lists prevent ads from appearing on specified publisher URLs or apps.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Fields  string `json:"fields"`
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

		path := fmt.Sprintf("/%s/publisher_block_lists", client.AdAccountID())
		return client.Get(ctx, path, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_block_list",
		Description: "Create a new publisher block list for the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"name":    {Type: "string", Description: "Block list name."},
			},
			Required: []string{"account", "name"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Name    string `json:"name"`
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

		path := fmt.Sprintf("/%s/publisher_block_lists", client.AdAccountID())
		return client.Post(ctx, path, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_block_list",
		Description: "Get details and entries of a specific publisher block list.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"block_list_id": {Type: "string", Description: "Publisher block list ID."},
				"fields":        {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "block_list_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			BlockListID string `json:"block_list_id"`
			Fields      string `json:"fields"`
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

		return client.Get(ctx, "/"+p.BlockListID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "add_to_block_list",
		Description: "Add publisher URLs or app store IDs to a block list. Provide a newline-separated list of URLs or app IDs.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"block_list_id": {Type: "string", Description: "Publisher block list ID."},
				"publisher_urls": {Type: "string", Description: "Newline-separated list of publisher URLs or app store IDs to block."},
			},
			Required: []string{"account", "block_list_id", "publisher_urls"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			BlockListID   string `json:"block_list_id"`
			PublisherURLs string `json:"publisher_urls"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		q.Set("publisher_urls", p.PublisherURLs)

		return client.Post(ctx, "/"+p.BlockListID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "remove_from_block_list",
		Description: "Remove publisher URLs or app store IDs from a block list.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"block_list_id": {Type: "string", Description: "Publisher block list ID."},
				"publisher_urls": {Type: "string", Description: "Newline-separated list of publisher URLs or app store IDs to remove from the block list."},
			},
			Required: []string{"account", "block_list_id", "publisher_urls"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account       string `json:"account"`
			BlockListID   string `json:"block_list_id"`
			PublisherURLs string `json:"publisher_urls"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		q := url.Values{}
		q.Set("publisher_urls", p.PublisherURLs)

		// Use DELETE with query params by constructing the request via Post with a _method override,
		// since the Graph API supports removing entries via POST with a delete flag.
		q.Set("delete_urls", "true")
		return client.Post(ctx, "/"+p.BlockListID, q)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_block_list",
		Description: "Delete a publisher block list entirely.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":       {Type: "string", Description: "Account name."},
				"block_list_id": {Type: "string", Description: "Publisher block list ID to delete."},
			},
			Required: []string{"account", "block_list_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account     string `json:"account"`
			BlockListID string `json:"block_list_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}

		return client.Delete(ctx, "/"+p.BlockListID)
	})
}
