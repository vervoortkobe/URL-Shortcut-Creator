package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func getDesktopPath(wshell *ole.IDispatch) (string, error) {
	desktop, err := oleutil.CallMethod(wshell, "SpecialFolders", "Desktop")
	if err != nil {
		return "", fmt.Errorf("failed to get desktop folder: %v", err)
	}
	return desktop.ToString(), nil
}

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
	fmt.Printf("Creating internet shortcut at: %s\n", shortcutPath)

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
