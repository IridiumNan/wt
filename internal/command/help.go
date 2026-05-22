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

	case "help", "--help", "-h":
		fmt.Println("  wt help")
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
		fmt.Println("  help     - Show help information")
	}

	fmt.Println("\nFor more information, use: wt help")
}
