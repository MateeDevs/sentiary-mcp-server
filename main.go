package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/MateeDevs/sentiary-mcp-server/internal/sentiary"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"

func main() {
	if isVersionCommand() {
		fmt.Println(versionString())
		return
	}

	if transportMode() == "http" {
		runHTTP()
		return
	}

	server := newServer(sentiary.ConfigFromEnv())
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("mcp server stopped: %v", err)
		os.Exit(1)
	}
}

func newServer(cfg sentiary.Config) *mcp.Server {
	client := sentiary.NewClient(cfg, &http.Client{Timeout: 30 * time.Second})

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sentiary",
		Title:   "Sentiary Strings",
		Version: version,
	}, nil)

	registerTools(server, client)
	return server
}

func isVersionCommand() bool {
	if len(os.Args) <= 1 {
		return false
	}
	arg := strings.ToLower(strings.TrimSpace(os.Args[1]))
	return arg == "version" || arg == "--version" || arg == "-v"
}

func versionString() string {
	if version != "" && version != "dev" {
		return version
	}
	buildInfo, ok := debug.ReadBuildInfo()
	if ok && buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	return version
}

func transportMode() string {
	if len(os.Args) > 1 {
		return strings.ToLower(strings.TrimSpace(os.Args[1]))
	}
	return strings.ToLower(strings.TrimSpace(os.Getenv("MCP_TRANSPORT")))
}

func runHTTP() {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	baseConfig := sentiary.ConfigFromEnv()

	// Build a fresh server per MCP session so a single hosted instance can serve
	// many users, each supplying their own API key from their MCP client. The
	// go-sdk calls this factory once per session (on initialize), where the
	// Authorization header is available. When no per-request key is present we
	// fall back to SENTIARY_USER_API_KEY from the environment.
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		cfg := baseConfig
		if key := apiKeyFromRequest(request); key != "" {
			cfg.APIKey = key
		}
		return newServer(cfg)
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.HandleFunc("GET /healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	addr := ":" + port
	log.Printf("mcp http server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("mcp http server stopped: %v", err)
		os.Exit(1)
	}
}

// apiKeyFromRequest extracts a per-user Sentiary API key from the incoming HTTP
// request. It accepts either a dedicated "X-Sentiary-User-Api-Key" header or a
// standard "Authorization" header ("Bearer <key>", "Ribbon <key>", or a bare
// "<key>"). Returns an empty string when no key is supplied, in which case the
// caller falls back to the server environment.
func apiKeyFromRequest(request *http.Request) string {
	if request == nil {
		return ""
	}
	if key := strings.TrimSpace(request.Header.Get("X-Sentiary-User-Api-Key")); key != "" {
		return key
	}
	authorization := strings.TrimSpace(request.Header.Get("Authorization"))
	if authorization == "" {
		return ""
	}
	if _, token, found := strings.Cut(authorization, " "); found {
		return strings.TrimSpace(token)
	}
	return authorization
}

func registerTools(server *mcp.Server, client *sentiary.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_project_strings",
		Title:       "List project strings",
		Description: "Lists project strings with pagination and optional translation languages. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.ListStringsInput) (*mcp.CallToolResult, sentiary.Paging, error) {
		output, err := client.ListStrings(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_project_string_by_name",
		Title:       "Search project string by name",
		Description: "Searches strings by name only. Use exact=true when you need a single exact key match. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.SearchStringInput) (*mcp.CallToolResult, sentiary.Paging, error) {
		output, err := client.SearchStringsByName(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project_string",
		Title:       "Get project string",
		Description: "Gets one string by ID, including translations. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.GetStringInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.GetString(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project_string_by_key",
		Title:       "Get project string by key",
		Description: "Gets one string by exact key/name, including requested translations. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.GetStringByKeyInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.GetStringByKey(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_project_string",
		Title:       "Add project string",
		Description: "Creates one string and optionally sets an initial translation. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.AddStringInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.AddString(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_project_string",
		Title:       "Edit project string",
		Description: "Edits string name, description, and context. Use clearDescription or clearContext to remove nullable fields. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.EditStringInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.EditString(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_project_string",
		Title:       "Remove project string",
		Description: "Removes one string and its translations. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.RemoveStringInput) (*mcp.CallToolResult, sentiary.DeleteOutput, error) {
		output, err := client.RemoveString(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_project_string_translation",
		Title:       "Set project string translation",
		Description: "Creates or updates one translation for a string. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.SetTranslationInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.SetStringTranslation(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_project_string_translation",
		Title:       "Remove project string translation",
		Description: "Removes one translation from a string. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.RemoveTranslationInput) (*mcp.CallToolResult, sentiary.Term, error) {
		output, err := client.RemoveStringTranslation(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project_strings_info",
		Title:       "Get project strings info",
		Description: "Returns batch string metadata for a project, including project name, languages, and last string modification date. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.ProjectInput) (*mcp.CallToolResult, sentiary.ProjectBatchInfo, error) {
		output, err := client.GetProjectBatchInfo(ctx, input.ProjectID)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "export_project_strings",
		Title:       "Export project strings",
		Description: "Exports project strings in json, android, apple, or compose format for one language. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.ExportInput) (*mcp.CallToolResult, sentiary.ExportOutput, error) {
		output, err := client.ExportStrings(ctx, input)
		return nil, output, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "import_project_strings",
		Title:       "Import project strings",
		Description: "Imports project strings in json, android, apple, or compose format for one language. Requires SENTIARY_USER_API_KEY.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input sentiary.ImportInput) (*mcp.CallToolResult, sentiary.ImportResult, error) {
		output, err := client.ImportStrings(ctx, input)
		return nil, output, err
	})
}
