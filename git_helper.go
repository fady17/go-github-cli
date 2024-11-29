package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Check if the current directory is a Git repository
	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
		fmt.Println("This directory is not a Git repository. Initializing a new Git repository...")
		if err := exec.Command("git", "init").Run(); err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized.")
	}

	// Get repository name or full URL from command-line arguments
	args := os.Args[1:]
	var repoArg string
	if len(args) > 0 {
		repoArg = args[0]
	} else {
		// Prompt user for repository name or URL if not provided
		fmt.Print("Enter the repository name or full URL: ")
		repoArg = readInput()
	}

	// Process repository argument to construct the remote URL
	var remoteURL string
	if strings.HasPrefix(repoArg, "git@") || strings.HasPrefix(repoArg, "https://") {
		// If a full URL is provided, use it directly
		remoteURL = repoArg
	} else {
		// Otherwise, treat it as a repository name and construct the URL
		remoteURL = fmt.Sprintf("git@github-personal:fady17/%s.git", repoArg)
	}

	// Set the Git remote URL
	if err := exec.Command("git", "remote", "add", "origin", remoteURL).Run(); err != nil {
		fmt.Println("Warning: Remote 'origin' might already exist. Attempting to update URL...")
		if err := exec.Command("git", "remote", "set-url", "origin", remoteURL).Run(); err != nil {
			fmt.Println("Error setting Git remote URL:", err)
			return
		}
	}
	fmt.Printf("Git remote set to %s.\n", remoteURL)

	// Test connection
	if err := exec.Command("git", "remote", "-v").Run(); err != nil {
		fmt.Println("Error testing Git remote:", err)
		return
	}

	// Add changes
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

	// Confirm and push
	if confirm("Are you sure you want to push to this repository? (y/n): ") {
		if err := exec.Command("git", "push", "-u", "origin", "main").Run(); err != nil {
			fmt.Println("Error pushing changes:", err)
			return
		}
		fmt.Println("Changes pushed successfully!")
	} else {
		fmt.Println("Push operation cancelled.")
	}
}

// readInput reads a line of input from the user
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// confirm asks the user for a yes/no confirmation
func confirm(prompt string) bool {
	fmt.Print(prompt)
	response := readInput()
	return strings.ToLower(response) == "y"
}


// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func main() {
// 	// Verify if we're in a Git repository
// 	if _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output(); err != nil {
// 		fmt.Println("Error: This directory is not a Git repository.")
// 		return
// 	}

// 	// Get repository name or full URL from command-line arguments
// 	args := os.Args[1:]
// 	var repoArg string
// 	if len(args) > 0 {
// 		repoArg = args[0]
// 	} else {
// 		// Prompt user for repository name or URL if not provided
// 		fmt.Print("Enter the repository name or full URL: ")
// 		repoArg = readInput()
// 	}

// 	// Process repository argument to construct the remote URL
// 	var remoteURL string
// 	if strings.HasPrefix(repoArg, "git@") {
// 		// If a full URL is provided, use it directly
// 		remoteURL = repoArg
// 	} else {
// 		// Otherwise, treat it as a repository name and construct the URL
// 		remoteURL = fmt.Sprintf("git@github-personal:fady17/%s.git", repoArg)
// 	}

// 	// Set the Git remote URL
// 	if err := exec.Command("git", "remote", "set-url", "origin", remoteURL).Run(); err != nil {
// 		fmt.Println("Error setting Git remote URL:", err)
// 		return
// 	}
// 	fmt.Printf("Git remote set to %s.\n", remoteURL)

// 	// Test connection
// 	if err := exec.Command("git", "remote", "-v").Run(); err != nil {
// 		fmt.Println("Error testing Git remote:", err)
// 		return
// 	}

// 	// Add changes
// 	fmt.Println("Adding all changes...")
// 	if err := exec.Command("git", "add", ".").Run(); err != nil {
// 		fmt.Println("Error adding changes:", err)
// 		return
// 	}

// 	// Confirm and push
// 	if confirm("Are you sure you want to push to this repository? (y/n): ") {
// 		if err := exec.Command("git", "push", "origin", "main").Run(); err != nil {
// 			fmt.Println("Error pushing changes:", err)
// 			return
// 		}
// 		fmt.Println("Changes pushed successfully!")
// 	} else {
// 		fmt.Println("Push operation cancelled.")
// 	}
// }

// // readInput reads a line of input from the user
// func readInput() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	input, _ := reader.ReadString('\n')
// 	return strings.TrimSpace(input)
// }

// // confirm asks the user for a yes/no confirmation
// func confirm(prompt string) bool {
// 	fmt.Print(prompt)
// 	response := readInput()
// 	return strings.ToLower(response) == "y"
// }
