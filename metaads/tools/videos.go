package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path/filepath"

	"github.com/richardpowellus/meta-ads-mcp-server/internal/paging"
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterVideos registers ad video management tools.
func RegisterVideos(s mcp.ToolRegistrar, cfg *metaads.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "upload_video",
		Description: "Upload a video to the ad account. Provide either file_url (HTTP URL) or file_path (local file). For URLs, uses standard form POST. For local files, uploads via the graph-video endpoint.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"file_url":  {Type: "string", Description: "Public URL of the video to upload. Mutually exclusive with file_path."},
				"file_path": {Type: "string", Description: "Local file path of the video to upload. Mutually exclusive with file_url."},
				"title":     {Type: "string", Description: "Title for the video."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			FileURL  string `json:"file_url"`
			FilePath string `json:"file_path"`
			Title    string `json:"title"`
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
		path := "/" + client.AdAccountID() + "/advideos"

		if p.FileURL != "" {
			vals := url.Values{}
			vals.Set("file_url", p.FileURL)
			if p.Title != "" {
				vals.Set("title", p.Title)
			}
			return client.Post(ctx, path, vals)
		}

		// Local file: read and upload via video graph endpoint
		data, err := os.ReadFile(p.FilePath)
		if err != nil {
			return nil, fmt.Errorf("reading file %q: %w", p.FilePath, err)
		}
		ext := filepath.Ext(p.FilePath)
		ct := mime.TypeByExtension(ext)
		if ct == "" {
			ct = "video/mp4"
		}
		reader := bytes.NewReader(data)
		return client.PostVideo(ctx, path, reader, ct)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_videos",
		Description: "List ad videos in the ad account.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"fields":  {Type: "string", Description: "Comma-separated fields (default: id,title,length,status,picture,created_time)."},
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
		path := "/" + client.AdAccountID() + "/advideos"
		vals := url.Values{}
		fields := p.Fields
		if fields == "" {
			fields = "id,title,length,status,picture,created_time"
		}
		vals.Set("fields", fields)
		items, err := client.FetchAll(ctx, path, "data", vals)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_video",
		Description: "Get details of a specific ad video.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"video_id": {Type: "string", Description: "The video ID."},
				"fields":   {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "video_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			VideoID string `json:"video_id"`
			Fields  string `json:"fields"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		path := "/" + p.VideoID
		vals := url.Values{}
		if p.Fields != "" {
			vals.Set("fields", p.Fields)
		}
		return client.Get(ctx, path, vals)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_video_thumbnails",
		Description: "Get thumbnail images for a specific ad video.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"video_id": {Type: "string", Description: "The video ID."},
			}, paging.Properties()),
			Required: []string{"account", "video_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			VideoID string `json:"video_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		path := "/" + p.VideoID + "/thumbnails"
		items, err := client.FetchAll(ctx, path, "data", nil)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_video",
		Description: "Delete an ad video.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"video_id": {Type: "string", Description: "The video ID to delete."},
			},
			Required: []string{"account", "video_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			VideoID string `json:"video_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, "/"+p.VideoID)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_video_upload_status",
		Description: "Check the upload/processing status of an ad video.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"video_id": {Type: "string", Description: "The video ID."},
			},
			Required: []string{"account", "video_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			VideoID string `json:"video_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("parsing params: %w", err)
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		vals := url.Values{}
		vals.Set("fields", "status,length,title,picture")
		return client.Get(ctx, "/"+p.VideoID, vals)
	})
}
