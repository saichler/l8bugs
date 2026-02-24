package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/saichler/l8types/go/ifs"
	"os"
)

const protocolVersion = "2024-11-05"

type ToolHandler func(args map[string]interface{}) (*CallToolResult, error)

type Server struct {
	vnic     ifs.IVNic
	tools    map[string]ToolHandler
	toolDefs []ToolDef
}

func NewServer(vnic ifs.IVNic) *Server {
	s := &Server{
		vnic:  vnic,
		tools: make(map[string]ToolHandler),
	}
	s.registerTools()
	return s
}

// CallTool invokes a registered tool by name with the given arguments.
// Used by tests to exercise tool handlers without JSON-RPC framing.
func (s *Server) CallTool(name string, args map[string]interface{}) (*CallToolResult, error) {
	handler, ok := s.tools[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return handler(args)
}

func (s *Server) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "[mcp] invalid JSON: %s\n", err)
			continue
		}

		resp := s.dispatch(&req)
		if resp == nil {
			continue // notification, no response needed
		}

		data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[mcp] marshal error: %s\n", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "%s\n", data)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[mcp] stdin read error: %s\n", err)
	}
}

func (s *Server) dispatch(req *Request) *Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		return nil // notification
	case "notifications/initialized":
		return nil // notification
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "ping":
		return &Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]interface{}{}}
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32601, Message: "method not found: " + req.Method},
		}
	}
}

func (s *Server) handleInitialize(req *Request) *Response {
	result := InitializeResult{
		ProtocolVersion: protocolVersion,
		Capabilities: Capabilities{
			Tools: &ToolsCap{ListChanged: false},
		},
		ServerInfo: ServerInfo{
			Name:    "l8bugs-mcp",
			Version: "1.0.0",
		},
	}
	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) handleToolsList(req *Request) *Response {
	result := ToolsListResult{Tools: s.toolDefs}
	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) handleToolsCall(req *Request) *Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32602, Message: "invalid params: " + err.Error()},
		}
	}

	handler, ok := s.tools[params.Name]
	if !ok {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: -32602, Message: "unknown tool: " + params.Name},
		}
	}

	result, err := handler(params.Arguments)
	if err != nil {
		result = &CallToolResult{
			Content: []ContentBlock{{Type: "text", Text: err.Error()}},
			IsError: true,
		}
	}

	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}
