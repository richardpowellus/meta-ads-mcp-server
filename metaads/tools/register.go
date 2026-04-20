package tools

import (
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
)

// RegisterAll registers all Meta Ads tools with the given server.
func RegisterAll(s mcp.ToolRegistrar, cfg *metaads.Config) {
	RegisterAccounts(s, cfg)
	RegisterCampaigns(s, cfg)
	RegisterAdsets(s, cfg)
	RegisterAds(s, cfg)
	RegisterCreatives(s, cfg)
	RegisterInsights(s, cfg)
	RegisterImages(s, cfg)
	RegisterVideos(s, cfg)
	RegisterBilling(s, cfg)
	RegisterTargeting(s, cfg)
	RegisterAudiences(s, cfg)
	RegisterAdLibrary(s, cfg)
	RegisterCanvas(s, cfg)
	RegisterBudgetPlanning(s, cfg)
	RegisterBrandSafety(s, cfg)
	RegisterLeads(s, cfg)
	RegisterCatalogs(s, cfg)
	RegisterBusiness(s, cfg)
	RegisterAuth(s, cfg)
	RegisterPreviews(s, cfg)
	RegisterPixels(s, cfg)
	RegisterAsyncReports(s, cfg)
	RegisterRules(s, cfg)
	RegisterDiagnostics(s, cfg)
	RegisterConversions(s, cfg)
	RegisterActivities(s, cfg)
	RegisterExperiments(s, cfg)
}
