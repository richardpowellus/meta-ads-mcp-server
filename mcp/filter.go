package mcp

// ToolRegistrar is implemented by anything that accepts tool registrations.
// Both Server and FilteredServer satisfy this interface, allowing tool
// registration code to be shared between full-access and read-only binaries.
type ToolRegistrar interface {
	RegisterTool(tool Tool, handler ToolHandler)
}

// FilteredRegistrar wraps a ToolRegistrar and only forwards tools whose
// names are in the allowlist. Tools not in the list are silently dropped.
//
// This enables deny-by-default read-only servers: new tools added to the
// codebase are automatically excluded unless explicitly allowlisted.
type FilteredRegistrar struct {
	inner   ToolRegistrar
	allowed map[string]bool
}

// NewFilteredRegistrar creates a registrar that only allows the named tools.
func NewFilteredRegistrar(inner ToolRegistrar, allowedTools []string) *FilteredRegistrar {
	m := make(map[string]bool, len(allowedTools))
	for _, t := range allowedTools {
		m[t] = true
	}
	return &FilteredRegistrar{inner: inner, allowed: m}
}

// RegisterTool forwards the registration only if the tool name is in the allowlist.
func (f *FilteredRegistrar) RegisterTool(tool Tool, handler ToolHandler) {
	if f.allowed[tool.Name] {
		f.inner.RegisterTool(tool, handler)
	}
}
