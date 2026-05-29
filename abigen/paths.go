package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// expandPath expands a user-facing file path into a clean local path.
func expandPath(p string) string {
	if strings.HasPrefix(p, `\\.\pipe`) {
		return p
	}
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := HomeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return filepath.Clean(os.ExpandEnv(p))
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
