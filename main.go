package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	filetype "version-generator/fileType"
	gittype "version-generator/gitType"
	"version-generator/versionSchemes"

	"github.com/alecthomas/kong"
)

// Version information - can be set at build time
var (
	Version   = "dev"     // Set via -ldflags at build time
	GitCommit = "unknown" // Set via -ldflags at build time
	BuildDate = "unknown" // Set via -ldflags at build time
)

type VersionInfo struct {
	Branch       string
	LastTag      string
	CommitsSince int
	ShortHash    string
	Version      string
}

type CLI struct {
	Version    kong.VersionFlag `kong:"short='v',help='Show version information'"`
	Semver     bool             `kong:"help='Use Semantic Versioning format'"`
	CalVer     bool             `kong:"help='Use Calendar Versioning format'"`
	Simple     bool             `kong:"help='Use simple version format (no branch info)'"`
	Hash       bool             `kong:"help='Include short hash in version'"`
	InBuiltGit bool             `kong:"short='i',help='Use built-in go-git library instead of system git'"`
	Go         bool             `kong:"short='g',help='Generate Go format version file'"`
	GoPath     string           `kong:"help='Path for Go file (default: version.go)',placeholder='PATH'"`
	Cpp        bool             `kong:"short='c',help='Generate C++ format version file'"`
	CppPath    string           `kong:"help='Path for C++ file (default: version.h)',placeholder='PATH'"`
	Yaml       bool             `kong:"short='y',help='Generate YAML format version file'"`
	YamlPath   string           `kong:"help='Path for YAML file (default: version.yaml)',placeholder='PATH'"`
	File       bool             `kong:"short='f',help='Write version to file'"`
	FilePath   string           `kong:"help='Path for file (default: .VERSION)',placeholder='PATH'"`
}

// getAppVersion returns the version of the application
func getAppVersion() string {
	// If version was set at build time, use it
	if Version != "dev" {
		return Version
	}

	// Fallback to git-based version detection for development
	gitHandler, err := gittype.GetGitHandler(false, ".")
	if err != nil {
		return "dev-unknown"
	}

	versionInfo, err := gitHandler.GenerateVersionInfo(false)
	if err != nil {
		return "dev-unknown"
	}

	return versionInfo.Version
}

func main() {
	var cli CLI

	// Get version for help display
	version := getAppVersion()

	kong.Parse(&cli,
		kong.Name("version-generator"),
		kong.Description(fmt.Sprintf("Git Version Generator - Generate version numbers from git repository state\n\nVersion: %s", version)),
		kong.Vars{"version": version},
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

	// Determine versioning options
	options := versionSchemes.VersioningOptions{
		Semver: cli.Semver,
		CalVer: cli.CalVer,
		Simple: cli.Simple,
		Hash:   cli.Hash,
	}

	// Generate version information based on options
	var versionInfo *gittype.VersionInfo
	if options.Semver || options.CalVer || options.Simple || options.Hash {
		versionInfo, err = gitHandler.GenerateVersionInfoWithOptions(options)
	} else {
		// Fallback to original method for backward compatibility
		versionInfo, err = gitHandler.GenerateVersionInfo(false)
	}
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
