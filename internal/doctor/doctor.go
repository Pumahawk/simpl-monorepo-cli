package doctor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func FindMissingRequirements() []error {
	_, err := exec.LookPath("docker")
	var errors []error
	if err != nil {
		errors = append(errors, fmt.Errorf("Not found docker. %w", err))
	}
	if err = checkDockerServerUp(); err != nil {
		errors = append(errors, err)
	}
	return errors
}

func checkDockerServerUp() error {
	cmd := exec.Command("docker", "--version", "{{.Server.Version}}")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Unable to get docker version. %w", err)
	}
	return nil
}

func FindMissingDependencies() []string {
	var missing []string
	if _, err := exec.LookPath("mise"); err != nil {
		return append(missing, "mise")
	}
	return missing
}

func SolveDependecies(names []string) []error {
	var errors []error
	fmt.Print("Solve dependecies\n")
	for _, name := range names {
		fmt.Printf("Solve %s\n", name)
		switch name {
		case "mise":
			if err := solveMiseDependecy(); err != nil {
				errors = append(errors, fmt.Errorf("Unable to solve mise. %w", err))
			}
		}
	}
	return errors
}

func solveMiseDependecy() error {
	res, err := http.Get("https://github.com/jdx/mise/releases/download/v2026.3.5/mise-v2026.3.5-linux-x64")
	if err != nil {
		return fmt.Errorf("Unable to download mise. %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Unable to download mise. Status code not 200.")
	}
	miseFile, err := createMiseLocation()
	if _, err := io.Copy(miseFile, res.Body); err != nil {
		return fmt.Errorf("Unable to download mise. %w", err)
	}
	return nil
}

func createMiseLocation() (*os.File, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Unable to get user home. %w", err)
	}
	dirPath := path + "/.simpl-monorepo-cli"
	if err := os.Mkdir(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("Unable to create configuration directory. %w", err)
	}

	file, err := os.Create(dirPath + "/mise")
	if err != nil {
		return nil, fmt.Errorf("Unable to create file mise. %w", err)
	}
	return file, nil
}
