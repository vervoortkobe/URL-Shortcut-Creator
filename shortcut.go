package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func createDesktopShortcut(targetURL, shortcutName, iconLocation string) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()

	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()

	desktopPath, err := getDesktopPath(wshell)
	if err != nil {
		return err
	}
	shortcutPath := filepath.Join(desktopPath, shortcutName+".lnk")
	fmt.Printf("Creating shortcut at: %s\n", shortcutPath)

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcutPath)
	if err != nil {
		return err
	}
	iDispatch := cs.ToIDispatch()

	_, err = oleutil.PutProperty(iDispatch, "TargetPath", targetURL)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(iDispatch, "IconLocation", iconLocation)
	if err != nil {
		return err
	}

	_, err = oleutil.CallMethod(iDispatch, "Save")
	if err != nil {
		return err
	}

	return nil
}

func getDesktopPath(wshell *ole.IDispatch) (string, error) {
	desktop, err := oleutil.CallMethod(wshell, "SpecialFolders", "Desktop")
	if err != nil {
		return "", err
	}
	return desktop.ToString(), nil
}
