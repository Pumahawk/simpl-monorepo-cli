package doctor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	if err := checkMiseDependency(); err != nil {
		fmt.Printf("Failed check mise dependency. %s\n", err)
		missing = append(missing, "mise")
	}
	if err := checkKubectlDependency(); err != nil {
		fmt.Printf("Failed check kubectl dependency. %s\n", err)
		missing = append(missing, "kubectl")
	}
	return missing
}

func checkMiseDependency() error {
	if _, err := exec.LookPath("mise"); err != nil {
		misePath, err := customMiseLocation()
		if err != nil {
			return fmt.Errorf("Invalid custom mise location. %w", err)
		} else {
			info, err := os.Stat(misePath)
			if err != nil {
				return fmt.Errorf("Invalid path info %s. %s\n", misePath, err)
			} else {
				mode := info.Mode()
				if mode&0111 == 0 {
					return fmt.Errorf("mise is not a exec file")
				}
			}
		}
	}
	return nil
}

func checkKubectlDependency() error {
	cmd, err := miseCommand("which", "kubectl")
	if err != nil {
		return fmt.Errorf("Unable to create mise command. %w", err)
	}
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Not found kubectl using mise where. %w", err)
	}
	return nil
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
		case "kubectl":
			if err := solveKubectlDependecy(); err != nil {
				errors = append(errors, fmt.Errorf("Unable to solve kubectl. %w", err))
			}
		}
	}
	return errors
}

func solveMiseDependecy() error {
	res, err := http.Get("https://github.com/jdx/mise/releases/download/v2026.3.5/mise-v2026.3.5-linux-x64")
	if err != nil {
		return fmt.Errorf("Unable to download mise. Error on get request. %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Unable to download mise. Status code not 200.")
	}
	miseFile, err := createMiseLocation()
	if err != nil {
		return fmt.Errorf("Unable to create mise file target. %s", err)
	}
	if _, err := io.Copy(miseFile, res.Body); err != nil {
		return fmt.Errorf("Unable to download mise. Error on write file. %w", err)
	}
	return nil
}

func solveKubectlDependecy() error {
	cmd, err := miseCommand("install", "kubectl@latest")
	if err != nil {
		return fmt.Errorf("Unable to create mise install command. %w", err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to install kubectl using mise. %w", err)
	}
	return nil
}

func createMiseLocation() (*os.File, error) {
	path, err := customMiseLocation()
	if err != nil {
		return nil, fmt.Errorf("Unable to get user home. %w", err)
	}

	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); err != nil {
		if err := os.Mkdir(parent, 0755); err != nil {
			return nil, fmt.Errorf("Unable to create configuration directory. %w", err)
		}
	}

	filepath.Dir(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("Unable to create file mise. %w", err)
	}
	return file, nil
}

func customMiseLocation() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Unable to get user home. %w", err)
	}
	return path + "/.simpl-monorepo-cli/mise", nil
}

func miseCommand(arg ...string) (*exec.Cmd, error) {
	var path string
	path, err := exec.LookPath("mise")
	if err != nil {
		path, err = customMiseLocation()
		if err != nil {
			return nil, fmt.Errorf("Unable to find mise location. %w", err)
		}
	}
	return exec.Command(path, arg...), nil
}
