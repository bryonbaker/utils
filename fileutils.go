package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// TODO: Replace stdout messages to a logging framework.

// BootConfig holds the startup config used to bootstrap all other config.
type BootConfig struct {
	ApplicationConfig string `json:"applicationConfig,omitempty"`
}

// LoadBootConfig loads the config file that the program uses. at startup. Returns error is the file cannot be loaded
func LoadBootConfig(file string) (BootConfig, error) {
	var bootConfig BootConfig

	// Find the directory that this app is running in and use it as the base directory for the config file.
	configFileAbsPath := GetExePath() + "/" + file
	fmt.Fprintf(os.Stdout, "Loading configuration from: ", configFileAbsPath)

	jsonFile, err := os.Open(configFileAbsPath)
	if err != nil {
		var msg = "WARNING: Configuration file: " + configFileAbsPath + " does not exist\n"
		fmt.Fprintf(os.Stderr, msg)
		return bootConfig, err
	}

	defer jsonFile.Close()
	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot read the json file contents.\n")
		return bootConfig, err
	}

	Must(json.Unmarshal(byteVal, &bootConfig))

	return bootConfig, nil
}

// GetExePath retrieves the fully qualified of the executable. This needs to consider symlinks - so they
// are traversed to make sure we find the real executable.
func GetExePath() string {
	var exePath string
	var err error

	exePath, err = os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Exe path = ", exePath)

	// Get the file info to see if it is a symlink.
	fileInfo, err := os.Lstat(exePath)
	if err != nil {
		log.Fatal(err)
	}

	// Double check the path is not a symlink. If so, get the target.
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// Get the path the symlink points to.

		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "Resolved symlink to be: ", exePath)
	}

	// Get the directory from the full path - removing any trailing / etc
	dir := filepath.Dir(exePath)

	return dir
}

// Must tests if error is nil and if not then the program terminates.s
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
