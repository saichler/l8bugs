/*
 * L8Bugs MCP Server — Stdio binary for AI coding agents.
 * Communicates via JSON-RPC 2.0 over stdin/stdout.
 */
package main

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/mcp"
	"github.com/saichler/l8bugs/go/bugs/website"
	"github.com/saichler/l8bugs/go/bugs/common"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "[mcp] L8Bugs MCP Server starting...")
	nic := website.CreateVnic(common.BUGS_VNET)
	fmt.Fprintln(os.Stderr, "[mcp] Connected to L8Bugs VNet")
	server := mcp.NewServer(nic)
	server.Run()
}
