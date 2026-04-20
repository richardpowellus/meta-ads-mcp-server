package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

// Server is an MCP server that communicates over stdio using JSON-RPC 2.0.
type Server struct {
	name         string
	version      string
	instructions string
	tools        map[string]registeredTool
	mu           sync.RWMutex
}

type registeredTool struct {
	Tool    Tool
	Handler ToolHandler
}

// ToolHandler processes a tool call and returns a result or error.
type ToolHandler func(ctx context.Context, params json.RawMessage) (any, error)

// New creates a new MCP server.
func New(name, version, instructions string) *Server {
	return &Server{
		name:         name,
		version:      version,
		instructions: instructions,
		tools:        make(map[string]registeredTool),
	}
}

// RegisterTool adds a tool to the server. Panics if a tool with the same name is already registered.
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tools[tool.Name]; exists {
		panic(fmt.Sprintf("duplicate tool registration: %q", tool.Name))
	}
	s.tools[tool.Name] = registeredTool{Tool: tool, Handler: handler}
}

// Run starts the server, reading JSON-RPC requests from stdin and writing responses to stdout.
// It blocks until ctx is cancelled or stdin is closed.
func (s *Server) Run(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("reading stdin: %w", err)
		}

		if len(line) == 0 || (len(line) == 1 && line[0] == '\n') {
			continue
		}

		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			resp := jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error:   &jsonRPCError{Code: -32700, Message: "Parse error"},
			}
			writeResponse(writer, resp)
			continue
		}

		resp := s.handleRequest(ctx, &req)
		if resp != nil {
			writeResponse(writer, *resp)
		}
	}
}

func (s *Server) handleRequest(ctx context.Context, req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		return nil // notification, no response
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	case "ping":
		return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: json.RawMessage(`{}`)}
	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonRPCError{Code: -32601, Message: "Method not found"},
		}
	}
}

func (s *Server) handleInitialize(req *jsonRPCRequest) *jsonRPCResponse {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{ListChanged: false},
		},
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
		Instructions: s.instructions,
	}

	data, _ := json.Marshal(result)
	return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: data}
}

func (s *Server) handleToolsList(req *jsonRPCRequest) *jsonRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]Tool, 0, len(s.tools))
	for _, rt := range s.tools {
		tools = append(tools, rt.Tool)
	}
	sort.Slice(tools, func(i, j int) bool { return tools[i].Name < tools[j].Name })

	result := ToolsListResult{Tools: tools}
	data, _ := json.Marshal(result)
	return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: data}
}

func (s *Server) handleToolsCall(ctx context.Context, req *jsonRPCRequest) *jsonRPCResponse {
	var call ToolCallRequest
	if err := json.Unmarshal(req.Params, &call); err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonRPCError{Code: -32602, Message: "Invalid params"},
		}
	}

	s.mu.RLock()
	rt, ok := s.tools[call.Name]
	s.mu.RUnlock()

	if !ok {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonRPCError{Code: -32602, Message: fmt.Sprintf("Unknown tool: %s", call.Name)},
		}
	}

	result, err := rt.Handler(ctx, call.Arguments)
	if err != nil {
		callResult := ToolCallResult{
			Content: []ContentBlock{{Type: "text", Text: fmt.Sprintf("Error: %s", err.Error())}},
			IsError: true,
		}
		data, _ := json.Marshal(callResult)
		return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: data}
	}

	text, err := json.Marshal(result)
	if err != nil {
		text = []byte(fmt.Sprintf("%v", result))
	}

	callResult := ToolCallResult{
		Content: []ContentBlock{{Type: "text", Text: string(text)}},
	}
	data, _ := json.Marshal(callResult)
	return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: data}
}

func writeResponse(w io.Writer, resp jsonRPCResponse) {
	data, _ := json.Marshal(resp)
	data = append(data, '\n')
	_, _ = w.Write(data)
}
