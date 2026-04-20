// Package paging provides response-level pagination for MCP tool results.
//
// Many tools fetch all results from an API and return them in a single response.
// When the response is too large, the MCP client silently dumps it to a temp file
// and the caller never sees the data. This package solves that by:
//
//  1. Adding page/page_size parameters to tools (via Properties + ParseParams)
//  2. Slicing results to the requested page (via Emit / EmitAny)
//  3. Auto-paginating when no page_size is set but the response exceeds MaxResponseBytes
package paging

import (
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
)

const MaxResponseBytes = 200 << 10 // 200KB

// Params holds response-level pagination parameters.
type Params struct {
	Page     int // 1-based, defaults to 1
	PageSize int // 0 means "return all" (auto-paginate if too large)
}

// ParseParams extracts page/page_size from raw JSON tool arguments.
func ParseParams(raw json.RawMessage) Params {
	var p struct {
		Page     *int `json:"page"`
		PageSize *int `json:"page_size"`
	}
	_ = json.Unmarshal(raw, &p)
	params := Params{Page: 1}
	if p.Page != nil && *p.Page > 0 {
		params.Page = *p.Page
	}
	if p.PageSize != nil && *p.PageSize > 0 {
		params.PageSize = *p.PageSize
	}
	return params
}

// Properties returns the standard page/page_size property schemas
// to merge into a tool's InputSchema.Properties.
func Properties() map[string]mcp.PropertySchema {
	return map[string]mcp.PropertySchema{
		"page":      {Type: "integer", Description: "Page number for response pagination (default: 1). Use with page_size to page through large result sets."},
		"page_size": {Type: "integer", Description: "Max results per response page. When omitted, all results are returned unless the response would be too large, in which case results are auto-paginated."},
	}
}

// Emit returns a paginated response for a typed slice of results.
// If page_size is set, returns the requested page.
// If page_size is not set and the serialized response exceeds MaxResponseBytes,
// auto-paginates and returns the requested page (default 1) with a note.
func Emit[T any](results []T, p Params) map[string]any {
	totalCount := len(results)

	if p.PageSize > 0 {
		return pageSlice(results, totalCount, p.Page, p.PageSize, false)
	}

	resp := map[string]any{"count": totalCount, "results": results}
	data, err := json.Marshal(resp)
	if err != nil || len(data) <= MaxResponseBytes {
		return resp
	}

	autoPageSize := estimatePageSize(len(data), totalCount)
	return pageSlice(results, totalCount, p.Page, autoPageSize, true)
}

// EmitAny handles untyped results (typically []any from FetchPaginated calls).
// If results is not a recognized slice type, returns it unchanged.
func EmitAny(results any, p Params) any {
	switch arr := results.(type) {
	case []any:
		return Emit(arr, p)
	case []map[string]any:
		return Emit(arr, p)
	default:
		return results
	}
}

func pageSlice[T any](results []T, totalCount, page, pageSize int, auto bool) map[string]any {
	totalPages := (totalCount + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	if start > totalCount {
		start = totalCount
	}
	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	resp := map[string]any{
		"count":       end - start,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"results":     results[start:end],
	}
	if auto {
		resp["note"] = fmt.Sprintf(
			"Result set too large to return at once (%d results). Showing page %d of %d (page_size=%d). "+
				"Request more pages with the page parameter, or narrow your query.",
			totalCount, page, totalPages, pageSize,
		)
	}
	return resp
}

func estimatePageSize(totalBytes, totalCount int) int {
	if totalCount == 0 {
		return 1
	}
	avgPerResult := totalBytes / totalCount
	if avgPerResult == 0 {
		avgPerResult = 1
	}
	ps := (MaxResponseBytes * 80 / 100) / avgPerResult // 80% budget leaves room for metadata
	if ps < 1 {
		ps = 1
	}
	return ps
}
