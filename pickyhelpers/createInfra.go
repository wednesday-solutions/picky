package pickyhelpers

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func IsInfraFilesExist() bool {
	path := fileutils.CurrentDirectory()
	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigJsFile,
		constants.Stacks,
	}
	for _, file := range files {
		status, _ := fileutils.IsExists(filepath.Join(path, file))
		if !status {
			return false
		}
	}
	return true
}

func CreateInfraSetup() error {

	infraFiles := make(map[string]string)
	// package.json file
	infraFiles[constants.PackageDotJsonFile] = sources.PackageDotJsonSource()
	// .env file
	infraFiles[constants.EnvFile] = sources.EnvFileSource()

	var err error
	var path string
	for file, source := range infraFiles {
		path = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), file)
		err = fileutils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)
	}
	return nil
}

func CreateInfraStacks(service, stack, dirName string) error {
	var err error
	var stackFileName string
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.Stacks)
	folderExist, _ := fileutils.IsExists(path)
	if !folderExist {
		err = fileutils.MakeDirectory(fileutils.CurrentDirectory(), constants.Stacks)
		errorhandler.CheckNilErr(err)
	}
	stackFileName = fmt.Sprintf("%s%s", dirName, ".js")
	path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(), constants.Stacks, stackFileName)
	var source string

	switch service {
	case constants.Web:
		source = sources.WebStackJsSource(dirName)
	case constants.Mobile:
		// not implemented
	case constants.Backend:
		source = sources.BackendStackJsSource(dirName)
	}
	err = fileutils.WriteToFile(path, source)
	return err
}

func CreateSstConfigFile(stackInfo map[string]interface{},
	all bool, dirName string, directories []string,
) error {
	sstConfigSource := sources.SstConfigJsSource()
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.SstConfigJsFile)
	stackInfo[constants.SstConfigStack] = dirName
	if all {
		// SST config file for all existing stacks.
		camelCaseDirectories := utils.ToCamelCase(directories)
		stackInfo[constants.ExistingDirectories] = camelCaseDirectories
	}
	err := hbs.ParseAndWriteToFile(sstConfigSource, path, stackInfo)
	return err
}

func UpdateEnvDevelopment(dirName string) error {
	path := fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
		dirName,
		constants.EnvDevFile,
	)
	err := fileutils.WriteToFile(path, sources.EnvDevSource())
	errorhandler.CheckNilErr(err)
	return nil
}
