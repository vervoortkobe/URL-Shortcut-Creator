package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createDesktopShortcut(targetURL, shortcutName, iconLocation string) error {
	userProfile, exists := os.LookupEnv("USERPROFILE")
	if !exists {
		return fmt.Errorf("USERPROFILE environment variable not found")
	}

	desktopPath := filepath.Join(userProfile, "Desktop")

	safeShortcutName := strings.Map(func(r rune) rune {
		if strings.ContainsRune(`\/:*?"<>|`, r) {
			return '_'
		}
		return r
	}, shortcutName)

	shortcutPath := filepath.Join(desktopPath, safeShortcutName+".url")
	fmt.Printf("✅ Creating desktop shortcut at: %s\n", shortcutPath)

	urlFileContent := fmt.Sprintf(`[InternetShortcut]
URL=%s
IconFile=%s
IconIndex=0
`, targetURL, iconLocation)

	err := os.WriteFile(shortcutPath, []byte(urlFileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create shortcut file: %v", err)
	}

	return nil
}
