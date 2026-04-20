package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/richardpowellus/meta-ads-mcp-server/internal/paging"
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterImages registers ad image management tools.
func RegisterImages(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "upload_image",
		Description: "Upload an image to the ad account. Provide either file_url (HTTP URL) or file_path (local file). Returns the image hash and URL.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"file_url":  {Type: "string", Description: "Public URL of the image to upload. Mutually exclusive with file_path."},
				"file_path": {Type: "string", Description: "Local file path of the image to upload. Mutually exclusive with file_url."},
				"name":      {Type: "string", Description: "Optional name for the image."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			FileURL  string `json:"file_url"`
			FilePath string `json:"file_path"`
			Name     string `json:"name"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		if p.FileURL == "" && p.FilePath == "" {
			return nil, fmt.Errorf("either file_url or file_path is required")
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/adimages"

		if p.FileURL != "" {
			vals := url.Values{}
			vals.Set("image_url", p.FileURL)
			if p.Name != "" {
				vals.Set("name", p.Name)
			}
			return client.Post(ctx, path, vals)
		}

		// Local file: read, base64-encode, POST as form field
		data, err := os.ReadFile(p.FilePath)
		if err != nil {
			return nil, fmt.Errorf("reading file %q: %w", p.FilePath, err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		vals := url.Values{}
		vals.Set("bytes", encoded)
		if p.Name != "" {
			vals.Set("name", p.Name)
		}
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_images",
		Description: "List ad images in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return (default: id,hash,name,url,status)."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Fields  string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + client.AdAccountID() + "/adimages"
		vals := url.Values{}
		fields := p.Fields
		if fields == "" {
			fields = "id,hash,name,url,status"
		}
		vals.Set("fields", fields)
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_image",
		Description: "Get details of a specific ad image by its hash.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"image_hash": {Type: "string", Description: "The image hash."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "image_hash"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ImageHash string `json:"image_hash"`
			Fields    string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/adimages"
		vals := url.Values{}
		vals.Set("hashes", fmt.Sprintf(`["%s"]`, p.ImageHash))
		if p.Fields != "" {
			vals.Set("fields", p.Fields)
		}
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_image",
		Description: "Delete an ad image by its hash.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"image_hash": {Type: "string", Description: "The image hash to delete."},
			},
			Required: []string{"account", "image_hash"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ImageHash string `json:"image_hash"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/adimages"
		vals := url.Values{}
		vals.Set("hash", p.ImageHash)
		return client.Post(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_image_crops",
		Description: "Get available image crops for an ad image.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"image_hash": {Type: "string", Description: "The image hash."},
			},
			Required: []string{"account", "image_hash"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ImageHash string `json:"image_hash"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/adimages"
		vals := url.Values{}
		vals.Set("hashes", fmt.Sprintf(`["%s"]`, p.ImageHash))
		vals.Set("fields", "id,hash,name,url,width,height")
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_image_by_url",
		Description: "Get image data by providing the image URL. Useful for validating images.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"image_url": {Type: "string", Description: "The URL of the image to check."},
			},
			Required: []string{"account", "image_url"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ImageURL string `json:"image_url"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + client.AdAccountID() + "/adimages"
		vals := url.Values{}
		vals.Set("image_url", p.ImageURL)
		vals.Set("fields", "id,hash,name,url,width,height,status")
		return client.Get(ctx, path, vals)
	})
}
