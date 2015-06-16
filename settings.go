package main

import (
	"fmt"
	"github.com/pilu/config"
	"os"
	"syscall"
	"path/filepath"
	"strconv"
	"time"
)

const (
	mainSettingsSection = "Settings"
)

type Settings map[string]string

var settings = Settings {
	"config_path":       "./runner.conf",
	"root":              ".",
	"output_path":       "./tmp",
	"build_name":        "runner-build",
	"build_log":         "runner-build-errors.log",
	"valid_ext":         ".go, .tpl, .tmpl, .html",
	"build_delay":       "600",
	"colors":            "1",
	"shutdown_signal":   "TERM",
	"log_color_main":    "cyan",
	"log_color_build":   "yellow",
	"log_color_runner":  "green",
	"log_color_watcher": "magenta",
	"log_color_app":     "",
}

var colors = map[string]string{
	"reset":          "0",
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bold_black":     "30;1",
	"bold_red":       "31;1",
	"bold_green":     "32;1",
	"bold_yellow":    "33;1",
	"bold_blue":      "34;1",
	"bold_magenta":   "35;1",
	"bold_cyan":      "36;1",
	"bold_white":     "37;1",
	"bright_black":   "30;2",
	"bright_red":     "31;2",
	"bright_green":   "32;2",
	"bright_yellow":  "33;2",
	"bright_blue":    "34;2",
	"bright_magenta": "35;2",
	"bright_cyan":    "36;2",
	"bright_white":   "37;2",
}

func (s Settings) logColor(logName string) string {
	settingsKey := fmt.Sprintf("log_color_%s", logName)
	colorName := s[settingsKey]

	return colors[colorName]
}

func (s Settings) load() {
	if _, err := os.Stat(s.configPath()); err != nil {
		return
	}

	logger.Printf("Loading settings from %s", s.configPath())
	sections, err := config.ParseFile(s.configPath(), mainSettingsSection)
	if err != nil {
		return
	}

	for key, value := range sections[mainSettingsSection] {
		s[key] = value
	}
}

func (s Settings) root() string {
	return s["root"]
}

func (s Settings) configPath() string {
	return s["config_path"]
}

func (s Settings) shutdownSignal() syscall.Signal {
	var signal syscall.Signal
	switch s["shutdown_signal"] {
	case "TERM":
		signal = syscall.SIGTERM
	case "KILL":
		signal = syscall.SIGKILL
	}
	return signal
}

func (s Settings) outputPath() string {
	return s["output_path"]
}

func (s Settings) postBuildScript() Script {
	return Script(s["post_build_script"])
}

func (s Settings) preBuildScript() Script {
	return Script(s["pre_build_script"])
}

func (s Settings) buildName() string {
	return s["build_name"]
}

func (s Settings) buildPath() string {
	return filepath.Join(s.outputPath(), s.buildName())
}

func (s Settings) buildErrorsFileName() string {
	return s["build_log"]
}

func (s Settings) buildErrorsFilePath() string {
	return filepath.Join(s.outputPath(), s.buildErrorsFileName())
}

func (s Settings) buildDelay() time.Duration {
	value, _ := strconv.Atoi(s["build_delay"])

	return time.Duration(value)
}
