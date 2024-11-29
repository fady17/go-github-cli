package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Parse command-line flags
	branch := flag.String("branch", "main", "Specify the branch to use (default: main)")
	flag.Parse()

	// Check if the current directory is a Git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
		fmt.Println("This directory is not a Git repository. Initializing a new Git repository...")
		if err := exec.Command("git", "init").Run(); err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized.")
	}

	// Ensure the specified branch exists
	ensureBranch(*branch)

	// Get repository name or full URL from remaining arguments
	args := flag.Args()
	var repoArg string
	if len(args) > 0 {
		repoArg = args[0]
	} else {
		fmt.Print("Enter the repository name or full URL: ")
		repoArg = readInput()
	}

	// Construct the remote URL
	remoteURL := constructRemoteURL(repoArg)

	// Add or update the Git remote
	if err := exec.Command("git", "remote", "add", "origin", remoteURL).Run(); err != nil {
		fmt.Println("Remote 'origin' already exists. Attempting to update...")
		if err := exec.Command("git", "remote", "set-url", "origin", remoteURL).Run(); err != nil {
			fmt.Println("Error setting Git remote URL:", err)
			return
		}
	}
	fmt.Printf("Git remote set to %s.\n", remoteURL)

	// Add all changes
	fmt.Println("Adding all changes...")
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		fmt.Println("Error adding changes:", err)
		return
	}

	// Commit changes
	fmt.Print("Enter commit message: ")
	commitMessage := readInput()
	if err := exec.Command("git", "commit", "-m", commitMessage).Run(); err != nil {
		fmt.Println("Error committing changes:", err)
		return
	}
	fmt.Println("Changes committed.")

	// Push to remote repository
	fmt.Printf("Pushing changes to branch '%s'...\n", *branch)
	if err := exec.Command("git", "push", "-u", "origin", *branch).Run(); err != nil {
		fmt.Println("Error pushing changes. Attempting to set upstream and retry...")
		if setupBranchAndPush(*branch) {
			fmt.Println("Changes pushed successfully!")
		} else {
			fmt.Println("Failed to push changes. Please check your remote settings and authentication.")
		}
	} else {
		fmt.Println("Changes pushed successfully!")
	}
}

// ensureBranch ensures the specified branch exists and switches to it if needed
func ensureBranch(branchName string) {
	output, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		fmt.Println("Error checking the current branch:", err)
		return
	}

	currentBranch := strings.TrimSpace(string(output))
	if currentBranch != branchName {
		fmt.Printf("Switching to branch '%s'...\n", branchName)
		if err := exec.Command("git", "checkout", "-B", branchName).Run(); err != nil {
			fmt.Println("Error creating or switching to branch:", err)
		}
	}
}

// constructRemoteURL generates a proper remote URL from user input
func constructRemoteURL(repoArg string) string {
	if strings.HasPrefix(repoArg, "git@") || strings.HasPrefix(repoArg, "https://") {
		return repoArg
	}
	return fmt.Sprintf("git@github-personal:fady17/%s.git", repoArg)
}

// setupBranchAndPush tries to set upstream branch and push
func setupBranchAndPush(branch string) bool {
	if err := exec.Command("git", "branch", "-M", branch).Run(); err != nil {
		fmt.Println("Error renaming branch to '"+branch+"':", err)
		return false
	}
	if err := exec.Command("git", "push", "-u", "origin", branch).Run(); err != nil {
		fmt.Println("Error pushing changes to '"+branch+"':", err)
		return false
	}
	return true
}

// readInput reads a line of input from the user
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}


