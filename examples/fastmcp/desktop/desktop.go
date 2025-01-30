package desktop

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/fastmcp"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

func main() {
	// Create a new FastMCP app
	app := fastmcp.New("Desktop Integration Demo")

	// Add a tool to open a file with the default application
	app.Tool("openFile", func(path string) error {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", path)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", path)
		default: // Linux and others
			cmd = exec.Command("xdg-open", path)
		}
		return cmd.Run()
	}, "Open a file with the default application")

	// Add a tool to get system information
	app.Tool("systemInfo", func() map[string]interface{} {
		hostname, _ := os.Hostname()
		workingDir, _ := os.Getwd()
		userCacheDir, _ := os.UserCacheDir()
		userConfigDir, _ := os.UserConfigDir()

		return map[string]interface{}{
			"os":              runtime.GOOS,
			"arch":            runtime.GOARCH,
			"num_cpus":        runtime.NumCPU(),
			"hostname":        hostname,
			"working_dir":     workingDir,
			"temp_dir":        os.TempDir(),
			"path_separator":  string(os.PathSeparator),
			"path_list_sep":   string(os.PathListSeparator),
			"user_home_dir":   os.Getenv("HOME"),
			"user_cache_dir":  userCacheDir,
			"user_config_dir": userConfigDir,
		}
	}, "Get system information")

	// Add a tool to manage clipboard (mock implementation)
	app.Tool("setClipboard", func(text string) error {
		log.Printf("Setting clipboard text: %s", text)
		return nil
	}, "Set clipboard text")

	// Add a resource to access environment variables
	app.Resource("env/{name}", func(name string) string {
		return os.Getenv(name)
	}, "Get environment variable value")

	// Add a prompt for file operations
	app.Prompt("fileOp", func(path string, operation string) []protocol.PromptMessage {
		return []protocol.PromptMessage{
			{
				Role: protocol.RoleAssistant,
				Content: protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("You are about to %s the file at %s.", operation, path),
				},
			},
			{
				Role: protocol.RoleUser,
				Content: protocol.TextContent{
					Type: "text",
					Text: "This operation may modify your system. Please confirm with 'yes' or 'no'.",
				},
			},
		}
	}, "Prompt for file operations")

	// Run the server with stdio transport
	log.Println("Starting Desktop Integration Demo...")
	if err := app.RunStdio(); err != nil {
		log.Fatal(err)
	}
}
