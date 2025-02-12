package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	owner    = "nikolanede"
	repo     = "employee-manager"
	filePath = "Jenkinsfile"
)

var rawFileURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", owner, repo, filePath)

func main() {
	fmt.Println("🚀 Starting GitHub Actions Importer Script...")

	fmt.Println("🔍 Fetching Jenkinsfile from GitHub...")
	content, err := fetchFile(rawFileURL)
	if err != nil {
		fmt.Println("❌ Failed to fetch Jenkinsfile:", err)
		return
	}
	fmt.Println("✅ Successfully fetched Jenkinsfile!")

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Println("❌ Error saving Jenkinsfile:", err)
		return
	}
	fmt.Println("✅ Jenkinsfile saved locally.")

	// Run GitHub Actions Importer
	fmt.Println("🚀 Running GitHub Actions Importer...")

	// Get GitHub Token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		fmt.Println("⚠️ Warning: GITHUB_TOKEN not set. Private repositories may not work.")
	}

	// Define Importer command
	cmd := exec.Command("gh", "actions-importer", "migrate", "jenkins",
		"--source-url", "https://github.com/"+owner+"/"+repo+"/blob/main/"+filePath,
		"--target-url", "https://github.com/"+owner+"/"+repo,
		"--output-dir", ".github/workflows",
		"--skip-fetch",
	)

	if githubToken != "" {
		cmd.Args = append(cmd.Args, "--jenkinsfile-access-token", githubToken)
	}

	// Run command and capture output
	fmt.Println("🔹 Running command:", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Error running GitHub Actions Importer:", err)
		fmt.Println("Output:", string(output))
		return
	}

	fmt.Println("Output:", string(output))
	fmt.Println("✅ GitHub Actions workflow generated successfully!")
	fmt.Println("📂 Check the output in .github/workflows")
}

// Fetches file content from GitHub
func fetchFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s\n%s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}
