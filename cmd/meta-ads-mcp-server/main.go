package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/richardpowellus/meta-ads-mcp-server/internal/tz"
	"github.com/richardpowellus/meta-ads-mcp-server/internal/version"
	"github.com/richardpowellus/meta-ads-mcp-server/internal/watchdog"
	"github.com/richardpowellus/meta-ads-mcp-server/mcp"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads"
	"github.com/richardpowellus/meta-ads-mcp-server/metaads/tools"
)

var instructions = `Meta Ads MCP Server — Facebook/Instagram advertising management.

Unofficial — not affiliated with or endorsed by Meta Platforms, Inc.

This server provides full read+write access to the Meta Marketing API v25.0
with 220+ tools covering every major endpoint.

Key capabilities:
- Campaigns: Create, manage, and optimize ad campaigns with bid strategies
- Ad Sets: Targeting, budgets, scheduling, promoted objects
- Ads: Creative management, status control, A/B testing
- Creative Assets: Image/video upload (multipart + resumable), ad creatives
- Audiences: Custom audiences, lookalikes, targeting search
- Insights: Campaign analytics with breakdowns, date ranges, async reports
- E-commerce: Product catalogs, feeds, dynamic ads
- Lead Gen: Lead forms and lead data retrieval
- Business: Account management, pages, Instagram accounts
- Billing: Transaction history and spend tracking

Configuration:
Set META_ADS_ACCESS_TOKEN environment variable (required).
Optional: META_APP_SECRET, META_APP_ID, META_AD_ACCOUNT_ID, META_BUSINESS_ID.
` + tz.Suffix()

func main() {
	log.SetOutput(os.Stderr)

	ctx := context.Background()
	ctx, cancel := watchdog.Start(ctx)
	defer cancel()

	cfg := metaads.NewEnvConfig()

	server := mcp.New("Meta Ads", version.Version, instructions)

	tools.RegisterAll(server, cfg)

	if err := server.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
