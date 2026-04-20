package version

// Version is the git commit short hash, injected at build time via:
//
//	go build -ldflags "-X github.com/richardpowellus/meta-ads-mcp-server/internal/version.Version=$(git rev-parse --short HEAD)"
//
// Falls back to "dev" if not set.
var Version = "dev"
