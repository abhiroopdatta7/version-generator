package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	filetype "version-generator/fileType"
	gittype "version-generator/gitType"

	"github.com/alecthomas/kong"
)

type VersionInfo struct {
	Branch       string
	LastTag      string
	CommitsSince int
	ShortHash    string
	Version      string
}

type CLI struct {
	Docker     bool   `kong:"short='d',help='Use Docker version format'"`
	InBuiltGit bool   `kong:"short='i',help='Use built-in go-git library instead of system git'"`
	Go         bool   `kong:"short='g',help='Generate Go format version file'"`
	GoPath     string `kong:"help='Path for Go file (default: version.go)',placeholder='PATH'"`
	Cpp        bool   `kong:"short='c',help='Generate C++ format version file'"`
	CppPath    string `kong:"help='Path for C++ file (default: version.h)',placeholder='PATH'"`
	Yaml       bool   `kong:"short='y',help='Generate YAML format version file'"`
	YamlPath   string `kong:"help='Path for YAML file (default: version.yaml)',placeholder='PATH'"`
	File       bool   `kong:"short='f',help='Write version to file'"`
	FilePath   string `kong:"help='Path for file (default: .VERSION)',placeholder='PATH'"`
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("version-generator"),
		kong.Description("Git Version Generator - Generate version numbers from git repository state"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	// Get git handler based on inBuiltGit flag
	gitHandler, err := gittype.GetGitHandler(cli.InBuiltGit, ".")
	if err != nil {
		log.Fatalf("Failed to initialize git handler: %v", err)
	}

	// Generate version information
	versionInfo, err := gitHandler.GenerateVersionInfo(cli.Docker)
	if err != nil {
		log.Fatalf("Failed to generate version info: %v", err)
	}

	// Determine output file and file type
	var filename string
	var fileTypeHandler filetype.FileType

	// Helper function to determine final path
	getFilePath := func(providedPath, defaultFilename string) string {
		if providedPath == "" {
			return defaultFilename
		}
		// Check if provided path is a directory (ends with /)
		if strings.HasSuffix(providedPath, "/") {
			return providedPath + defaultFilename
		}
		return providedPath
	}

	// Determine file type based on flags
	switch {
	case cli.Go:
		fileTypeHandler = &filetype.GoType{}
		filename = getFilePath(cli.GoPath, "version.go")
	case cli.Cpp:
		fileTypeHandler = &filetype.CPPType{}
		filename = getFilePath(cli.CppPath, "version.h")
	case cli.Yaml:
		fileTypeHandler = &filetype.YAMLFile{}
		filename = getFilePath(cli.YamlPath, "version.yaml")
	case cli.File:
		fileTypeHandler = &filetype.BasicFile{}
		filename = getFilePath(cli.FilePath, ".VERSION")
	}

	// Print only the version string (unless file type format is used)
	if fileTypeHandler == nil {
		fmt.Println(versionInfo.Version)
	}

	// Write to file if requested or file type format is specified
	if filename != "" && fileTypeHandler != nil {
		err := fileTypeHandler.WriteVersion(filename, versionInfo.Version)
		if err != nil {
			log.Fatalf("Failed to write version to file %s: %v", filename, err)
		}
	} else if filename != "" {
		// Fallback to basic file writing
		err := writeVersionToFile(filename, versionInfo.Version)
		if err != nil {
			log.Fatalf("Failed to write version to file %s: %v", filename, err)
		}
	}
}

func writeVersionToFile(filename, version string) error {
	return os.WriteFile(filename, []byte(version+"\n"), 0644)
}
