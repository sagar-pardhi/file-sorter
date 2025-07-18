package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// FileCategory defines the mapping between file extensions and their categories
type FileCategory struct {
	Extensions []string
	FolderName string
}

var (
	categories = []FileCategory{
		{
			Extensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".svg"},
			FolderName: "Pictures",
		},
		{
			Extensions: []string{".mp4", ".mov", ".avi", ".mkv", ".flv", ".wmv", ".webm"},
			FolderName: "Videos",
		},
		{
			Extensions: []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma"},
			FolderName: "Music",
		},
		{
			Extensions: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf"},
			FolderName: "Documents",
		},
		{
			Extensions: []string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2"},
			FolderName: "Archives",
		},
		{
			Extensions: []string{".exe", ".msi", ".dmg", ".pkg", ".deb", ".rpm"},
			FolderName: "Programs",
		},
	}

	// Other files will go to this folder
	otherFolder = "Other"
)

func main() {
	// Parse command line flags
	dirPath := flag.String("dir", ".", "Directory path to sort files")
	dryRun := flag.Bool("dry-run", false, "Dry run - show what would be done without actually moving files")
	flag.Parse()

	// Validate directory path
	if *dirPath == "" {
		fmt.Println("Please specify a directory path with -dir flag")
		os.Exit(1)
	}

	// Get absolute path
	absPath, err := filepath.Abs(*dirPath)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	fmt.Printf("Sorting files in: %s\n", absPath)
	if *dryRun {
		fmt.Println("DRY RUN MODE - No files will be actually moved")
	}

	// Read directory
	files, err := os.ReadDir(absPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))

		// Find the category for this file
		category := findCategory(ext)
		targetFolder := filepath.Join(absPath, category)

		// Create target folder if it doesn't exist
		if !*dryRun {
			if err := os.MkdirAll(targetFolder, 0755); err != nil {
				log.Printf("Error creating folder %s: %v", targetFolder, err)
				continue
			}
		}

		oldPath := filepath.Join(absPath, filename)
		newPath := filepath.Join(targetFolder, filename)

		if *dryRun {
			fmt.Printf("[Dry Run] Would move: %s -> %s\n", oldPath, newPath)
		} else {
			if err := os.Rename(oldPath, newPath); err != nil {
				log.Printf("Error moving file %s: %v", filename, err)
			} else {
				fmt.Printf("Moved: %s -> %s\n", oldPath, newPath)
			}
		}
	}

	fmt.Println("File sorting completed!")
}

// findCategory returns the folder name for the given file extension
func findCategory(ext string) string {
	for _, cat := range categories {
		if slices.Contains(cat.Extensions, ext) {
			return cat.FolderName
		}
	}
	return otherFolder
}
