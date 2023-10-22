package services

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateSDKApp(projectName string) error {
	baseDir := "./cloned"

	// Ensure the base directory exists
	if err := createFolderIfNotExist(baseDir); err != nil {
		return err
	}

	// Use npx command with project name
	cmd := exec.Command("npx", "packlify-start-app", projectName)
	cmd.Dir = baseDir
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("SDK App created.")

	return nil
}

func PushToGitHubRepo(repoName, cloneURL string) error {
	projectPath := "./cloned/" + repoName

	// Change directory to the cloned project
	err := os.Chdir(projectPath)
	if err != nil {
		return err
	}

	// Git Init
	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Git Add all files
	cmd = exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Git Commit
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	remoteURL := fmt.Sprintf(
		"https://%s:%s@github.com/%s/%s.git",
		os.Getenv("GITHUB_USERNAME"),
		os.Getenv("GITHUB_ACCESS_TOKEN"),
		os.Getenv("GITHUB_OWNER"),
		repoName)

	// Git Remote Add
	cmd = exec.Command("git", "remote", "add", "origin", remoteURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Git Push
	cmd = exec.Command("git", "push", "-u", "origin", "main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
