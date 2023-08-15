package atlantis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

const tfvarsExtension = ".tfvars"

func multiWorkspaceGetProjectScope(relPath, patternDetector string, changedFiles []string) string {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", relPath)) &&
			!strings.Contains(file, patternDetector) {
			return "crossWorkspace"
		}
	}
	return "workspace"
}

func multiWorkspaceGenWorkspaceList(relPath string, changedFiles []string, scope string) (workspaceList []string, err error) {
	err = filepath.Walk(relPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), tfvarsExtension) {
			if helpers.IsStringInList(path, changedFiles) || (scope == "crossWorkspace") {
				workspaceList = append(workspaceList, helpers.TrimFileExtension(info.Name()))
			}
		}
		return nil
	})
	return workspaceList, err
}

func multiWorkspaceDetectProjectWorkspaces(changedFiles []string, foldersList []ProjectFolder, patternDetector string) (updatedFolderList []ProjectFolder, err error) {

	for i := range foldersList {
		scope := multiWorkspaceGetProjectScope(foldersList[i].Path, patternDetector, changedFiles)
		workspaceList, err := multiWorkspaceGenWorkspaceList(foldersList[i].Path, changedFiles, scope)
		if err != nil {
			return foldersList, err
		}
		foldersList[i].WorkspaceList = workspaceList
	}
	return foldersList, nil
}

func multiWorkspaceWorkflowFilter(info os.FileInfo, path, patternDetector string) bool {
	return info.IsDir() &&
		info.Name() == patternDetector &&
		!strings.Contains(path, ".terraform")

}