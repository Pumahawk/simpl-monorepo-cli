package doctor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var DependencyList = [...]Dependency{
	miseD,
	NewMiseDependency("kubectl"),
	NewMiseDependency("helm"),
}

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

func FindMissingDependencies() []Dependency {
	var missing []Dependency
	for _, dep := range DependencyList {
		if err := dep.Check(); err != nil {
			fmt.Printf("Failed check %s dependency. %s\n", dep.Name, err)
			missing = append(missing, dep)
		}
	}
	return missing
}

func SolveDependecies(dependecies []Dependency) []error {
	var errors []error
	fmt.Print("Solve dependecies\n")
	for _, dep := range dependecies {
		fmt.Printf("Solve %s\n", dep.Name)
		if err := dep.Solve(); err != nil {
			errors = append(errors, fmt.Errorf("Unable to solve %s. %w", dep.Name, err))
		}
	}
	return errors
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

type Dependency struct {
	Name  string
	Check func() error
	Solve func() error
}

var miseD = Dependency{
	Name: "mise",
	Check: func() error {
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
	},
	Solve: func() error {
		res, err := http.Get("https://github.com/jdx/mise/releases/download/v2026.3.6/mise-v2026.3.6-linux-x64")
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
		defer miseFile.Close()
		if _, err := io.Copy(miseFile, res.Body); err != nil {
			return fmt.Errorf("Unable to download mise. Error on write file. %w", err)
		}
		return nil
	},
}

func NewMiseDependency(name string) Dependency {
	return Dependency{
		Name: name,
		Check: func() error {
			cmd, err := miseCommand("where", name)
			if err != nil {
				return fmt.Errorf("Unable to create mise command. %w", err)
			}
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("Not found %s using mise where. %w", name, err)
			}
			return nil
		},
		Solve: func() error {
			cmd, err := miseCommand("install", name)
			if err != nil {
				return fmt.Errorf("Unable to create mise install command. %w", err)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("Unable to install %s using mise. %w", name, err)
			}
			return nil
		},
	}
}
