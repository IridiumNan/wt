package command

import (
	_ "embed"
	"fmt"
)

const (
	CommandIndex      = 0
	FirstTargetIndex  = 0
	SecondTargetIndex = 1
	ThirdTargetIndex  = 2
	FourthTargetIndex = 3
	FifthTargetIndex  = 4
	SixthTargetIndex  = 5
)

//go:embed Default.txt
var DefaultManual string

//go:embed Simple.txt
var SimpleManual string

//go:embed Advance.txt
var AdvanceManual string

func Usage(command string) {
	fmt.Println("\n=== Usage Guide ===")

	switch command {
	case "search":
		fmt.Println("  wt search <package-name>")
		fmt.Println("  Example: wt search my-app")

	case "info":
		fmt.Println("  wt info <package-name>")
		fmt.Println("  Example: wt info my-app-v1")

	case "install":
		fmt.Println("  wt install <package-name>")
		fmt.Println("  Example: wt install my-app-v1")

	case "upload":
		fmt.Println("  wt upload <file-path> [package-name]")
		fmt.Println("  Example: wt upload ./build/app.tar.gz")
		fmt.Println("  Example: wt upload ./build/app.tar.gz my-app-v1")

	case "mv":
		fmt.Println("  wt mv <old-name> <new-name>")
		fmt.Println("  Example: wt mv my-app-v1 my-app-v2")

	case "rm":
		fmt.Println("  wt rm <package-name>")
		fmt.Println("  Example: wt rm my-app-v1")

	case "list", "ls":
		fmt.Println("  wt list [tag]")
		fmt.Println("  Example: wt list")
		fmt.Println("  Example: wt list latest")

	case "sync":
		fmt.Println("  wt sync")
		fmt.Println("  Description: Sync local metadata with server")

	case "tag":
		fmt.Println("  wt tag <subcommand> [args]")
		fmt.Println("  Subcommands: add, rm, list, <pkg> <tag>")
		fmt.Println("  Example: wt tag add stable")
		fmt.Println("  Example: wt tag my-app stable")

	case "public":
		fmt.Println("  wt public <package-name>")
		fmt.Println("  Example: wt public my-app")
		fmt.Println("  Description: Make package publicly accessible")

	case "private":
		fmt.Println("  wt private <package-name>")
		fmt.Println("  Example: wt private my-app")
		fmt.Println("  Description: Revoke public access for package")

	case "links":
		fmt.Println("  wt links")
		fmt.Println("  Description: List all public package links")

	case "list-servers":
		fmt.Println("  wt list-servers")
		fmt.Println("  Description: Show configured client servers")

	case "change-server":
		fmt.Println("  wt change-server <server-name>")
		fmt.Println("  Example: wt change-server prod")
		fmt.Println("  Description: Change the active client server")

	case "add-server":
		fmt.Println("  wt add-server <name> <server-url>")
		fmt.Println("  Example: wt add-server prod http://192.168.1.2:12212")
		fmt.Println("  Description: Add a named client server")

	case "del-server":
		fmt.Println("  wt del-server <server-name>")
		fmt.Println("  Example: wt del-server prod")
		fmt.Println("  Description: Remove a configured client server")

	case "reload":
		fmt.Println("  wt reload")
		fmt.Println("  Description: Reload server configuration")

	case "help", "--help", "-h":
		fmt.Println("  wt help")
		fmt.Println("  wt help <command>")
		fmt.Println("  wt help simple")
		fmt.Println("  wt help advance")
		fmt.Println("  Description: Show this help manual")

	default:
		fmt.Println("  Unknown command:", command)
		fmt.Printf("\nAvailable commands:\n")
		fmt.Println("  search   - Search for packages")
		fmt.Println("  info     - Show package information")
		fmt.Println("  install  - Download and install a package")
		fmt.Println("  upload   - Upload a package to server")
		fmt.Println("  mv       - Rename a package")
		fmt.Println("  rm       - Remove a package")
		fmt.Println("  list     - List packages by tag")
		fmt.Println("  sync     - Sync metadata with server")
		fmt.Println("  tag      - Manage tags and package assignments")
		fmt.Println("  public   - Make package publicly accessible")
		fmt.Println("  private  - Revoke public access")
		fmt.Println("  links    - List public package links")
		fmt.Println("  list-servers - List configured client servers")
		fmt.Println("  change-server - Change the active client server")
		fmt.Println("  add-server - Add a new client server")
		fmt.Println("  del-server - Remove a configured client server")
		fmt.Println("  reload   - Reload server configuration")
		fmt.Println("  help     - Show help information")
	}

	fmt.Println("\nFor more information, use: wt help")
}
