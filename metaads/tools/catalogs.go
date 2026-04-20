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

// RegisterCatalogs registers product catalog management tools.
func RegisterCatalogs(s mcp.ToolRegistrar, cfg *metaads.Config) {
	// --- Catalogs ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_catalogs",
		Description: "List product catalogs owned by the business.",
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
			v.Set("fields", "id,name,product_count,vertical")
		}
		items, err := cl.FetchAll(ctx, fmt.Sprintf("/%s/owned_product_catalogs", cl.BusinessID()), "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_catalog",
		Description: "Get details of a specific product catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "catalog_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Fields    string `json:"fields"`
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
		return cl.Get(ctx, "/"+p.CatalogID, v)
	})

	// --- Product Sets ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_product_sets",
		Description: "List product sets in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
			}, paging.Properties()),
			Required: []string{"account", "catalog_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Fields    string `json:"fields"`
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
			v.Set("fields", "id,name,product_count,filter")
		}
		items, err := cl.FetchAll(ctx, "/"+p.CatalogID+"/product_sets", "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_product_set",
		Description: "Create a product set in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"name":       {Type: "string", Description: "Product set name."},
				"filter":     {Type: "string", Description: "JSON filter rules for the product set."},
			},
			Required: []string{"account", "catalog_id", "name"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Name      string `json:"name"`
			Filter    string `json:"filter"`
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
		if p.Filter != "" {
			v.Set("filter", p.Filter)
		}
		return cl.Post(ctx, "/"+p.CatalogID+"/product_sets", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_product_set",
		Description: "Get details of a specific product set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"product_set_id": {Type: "string", Description: "Product set ID."},
				"fields":         {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "product_set_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ProductSetID string `json:"product_set_id"`
			Fields       string `json:"fields"`
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
		return cl.Get(ctx, "/"+p.ProductSetID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_product_set",
		Description: "Delete a product set.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"product_set_id": {Type: "string", Description: "Product set ID to delete."},
			},
			Required: []string{"account", "product_set_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ProductSetID string `json:"product_set_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.ProductSetID)
	})

	// --- Products ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_products",
		Description: "List products in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
				"filter":     {Type: "string", Description: "JSON filter for products."},
			}, paging.Properties()),
			Required: []string{"account", "catalog_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Fields    string `json:"fields"`
			Filter    string `json:"filter"`
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
			v.Set("fields", "id,retailer_id,name,price,availability,image_url,url")
		}
		if p.Filter != "" {
			v.Set("filter", p.Filter)
		}
		items, err := cl.FetchAll(ctx, "/"+p.CatalogID+"/products", "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_product",
		Description: "Get details of a specific product in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"product_id": {Type: "string", Description: "Product ID."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ProductID string `json:"product_id"`
			Fields    string `json:"fields"`
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
		return cl.Get(ctx, "/"+p.ProductID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_product",
		Description: "Create a product in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":      {Type: "string", Description: "Account name."},
				"catalog_id":   {Type: "string", Description: "Product catalog ID."},
				"retailer_id":  {Type: "string", Description: "Unique retailer product ID."},
				"name":         {Type: "string", Description: "Product name."},
				"price":        {Type: "string", Description: "Price with currency (e.g. '9.99 USD')."},
				"availability": {Type: "string", Description: "Availability (in stock, out of stock)."},
				"url":          {Type: "string", Description: "Product page URL."},
				"image_url":    {Type: "string", Description: "Product image URL."},
				"description":  {Type: "string", Description: "Product description."},
				"category":     {Type: "string", Description: "Product category."},
				"brand":        {Type: "string", Description: "Brand name."},
			},
			Required: []string{"account", "catalog_id", "retailer_id", "name", "price", "availability", "url", "image_url"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			CatalogID    string `json:"catalog_id"`
			RetailerID   string `json:"retailer_id"`
			Name         string `json:"name"`
			Price        string `json:"price"`
			Availability string `json:"availability"`
			URL          string `json:"url"`
			ImageURL     string `json:"image_url"`
			Description  string `json:"description"`
			Category     string `json:"category"`
			Brand        string `json:"brand"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("retailer_id", p.RetailerID)
		v.Set("name", p.Name)
		v.Set("price", p.Price)
		v.Set("availability", p.Availability)
		v.Set("url", p.URL)
		v.Set("image_url", p.ImageURL)
		if p.Description != "" {
			v.Set("description", p.Description)
		}
		if p.Category != "" {
			v.Set("category", p.Category)
		}
		if p.Brand != "" {
			v.Set("brand", p.Brand)
		}
		return cl.Post(ctx, "/"+p.CatalogID+"/products", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_product",
		Description: "Delete a product from a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"product_id": {Type: "string", Description: "Product ID to delete."},
			},
			Required: []string{"account", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ProductID string `json:"product_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.ProductID)
	})

	// --- Product Feeds ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_product_feeds",
		Description: "List product feeds in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"fields":     {Type: "string", Description: "Comma-separated fields to return."},
			}, paging.Properties()),
			Required: []string{"account", "catalog_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Fields    string `json:"fields"`
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
			v.Set("fields", "id,name,product_count,schedule,update_schedule,latest_upload")
		}
		items, err := cl.FetchAll(ctx, "/"+p.CatalogID+"/product_feeds", "data", v)
		if err != nil {
			return nil, err
		}
		return paging.Emit(items, paging.ParseParams(params)), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_product_feed",
		Description: "Create a product feed for a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"name":       {Type: "string", Description: "Feed name."},
				"schedule":   {Type: "string", Description: "JSON schedule object (url, interval, hour)."},
			},
			Required: []string{"account", "catalog_id", "name"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Name      string `json:"name"`
			Schedule  string `json:"schedule"`
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
		if p.Schedule != "" {
			v.Set("schedule", p.Schedule)
		}
		return cl.Post(ctx, "/"+p.CatalogID+"/product_feeds", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_product_feed",
		Description: "Get details of a specific product feed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"feed_id": {Type: "string", Description: "Product feed ID."},
				"fields":  {Type: "string", Description: "Comma-separated fields to return."},
			},
			Required: []string{"account", "feed_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			FeedID  string `json:"feed_id"`
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
		return cl.Get(ctx, "/"+p.FeedID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_product_feed",
		Description: "Update a product feed (name, schedule).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"feed_id":  {Type: "string", Description: "Product feed ID."},
				"name":     {Type: "string", Description: "New feed name."},
				"schedule": {Type: "string", Description: "JSON schedule object."},
			},
			Required: []string{"account", "feed_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			FeedID   string `json:"feed_id"`
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		if p.Name != "" {
			v.Set("name", p.Name)
		}
		if p.Schedule != "" {
			v.Set("schedule", p.Schedule)
		}
		return cl.Post(ctx, "/"+p.FeedID, v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_product_feed",
		Description: "Delete a product feed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"feed_id": {Type: "string", Description: "Product feed ID to delete."},
			},
			Required: []string{"account", "feed_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			FeedID  string `json:"feed_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return cl.Delete(ctx, "/"+p.FeedID)
	})

	// --- Batch Operations ---

	s.RegisterTool(mcp.Tool{
		Name:        "batch_products",
		Description: "Batch create, update, or delete products in a catalog.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
				"requests":   {Type: "string", Description: "JSON array of batch request objects (method, retailer_id, data)."},
			},
			Required: []string{"account", "catalog_id", "requests"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			CatalogID string `json:"catalog_id"`
			Requests  string `json:"requests"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("requests", p.Requests)
		return cl.Post(ctx, "/"+p.CatalogID+"/items_batch", v)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_batch_status",
		Description: "Get the status of a batch product operation.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"handle":   {Type: "string", Description: "Batch handle returned from batch_products."},
				"catalog_id": {Type: "string", Description: "Product catalog ID."},
			},
			Required: []string{"account", "catalog_id", "handle"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			Handle    string `json:"handle"`
			CatalogID string `json:"catalog_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		cl, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		v := url.Values{}
		v.Set("handle", p.Handle)
		return cl.Get(ctx, "/"+p.CatalogID+"/check_batch_request_status", v)
	})
}
