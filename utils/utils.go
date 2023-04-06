package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func DirectoryName(dirName, stack, database string) string {
	switch stack {
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiMySqlTemplate)
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlMySqlTemplate)
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeExpressMongoTemplate)
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangMySqlTemplate)
		}
	case constants.ReactJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactTemplate)
	case constants.NextJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.NextTemplate)
	case constants.ReactNative:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactNativeTemplate)
	case constants.Android:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.AndroidTemplate)
	case constants.IOS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.IOSTemplate)
	case constants.Flutter:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.FlutterTemplate)
	}
	return dirName
}

func ExistingStacksDatabasesAndDirectories() ([]string, []string, []string) {
	var stacks, databases, dirNames []string
	var stack, database string
	directories, err := fileutils.ReadAllContents(fileutils.CurrentDirectory())
	errorhandler.CheckNilErr(err)

	for _, dirName := range directories {
		stack, database = FindStackAndDatabase(dirName)
		if stack != "" {
			stacks = append(stacks, stack)
			databases = append(databases, database)
			dirNames = append(dirNames, dirName)
		}
	}
	return stacks, databases, dirNames
}

func FindStackAndDatabase(dirName string) (string, string) {
	var first, second, stack, database string
	splitDirName := strings.Split(dirName, "-")
	if len(splitDirName) > 2 {
		first = splitDirName[len(splitDirName)-1]
		second = splitDirName[len(splitDirName)-2]
		switch first {
		case "pg":
			database = constants.PostgreSQL
			if second == "hapi" {
				stack = constants.NodeHapiTemplate
			} else if second == "graphql" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if second == "golang" {
				stack = constants.GolangEchoTemplate
			}
		case "mysql":
			database = constants.MySQL
			if second == "hapi" {
				stack = constants.NodeHapiTemplate
			} else if second == "graphql" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if second == "golang" {
				stack = constants.GolangEchoTemplate
			}
		case "mongo":
			database = constants.MongoDB
			if second == "express" {
				stack = constants.NodeExpressTemplate
			}
		case "web":
			if second == "react" {
				stack = constants.ReactJS
			} else if second == "next" {
				stack = constants.NextJS
			}
		case "mobile":
			if second == "reactnative" {
				stack = constants.ReactNative
			} else if second == "android" {
				stack = constants.Android
			} else if second == "ios" {
				stack = constants.IOS
			} else if second == "flutter" {
				stack = constants.Flutter
			}
		}
	}
	return stack, database
}

func ExistingStackAndDatabase(dirName string) (string, string) {
	stack, database := FindStackAndDatabase(dirName)
	return stack, database
}

func FindService(dirName string) string {
	splitDirName := strings.Split(dirName, "-")
	if len(splitDirName) > 2 {
		suffix := splitDirName[len(splitDirName)-1]
		switch suffix {
		case "pg", "mysql":
			return constants.Backend
		case "web":
			return constants.Web
		case "mobile":
			return constants.Mobile
		}
	}
	return ""
}

func ToCamelCase(slice []string) []string {
	camelSlice := []string{}
	for _, str := range slice {
		camelSlice = append(camelSlice, strcase.ToCamel(str))
	}
	return camelSlice
}

func CreateTemplate(name, text string) *template.Template {
	tpl, err := template.New(name).Parse(text)
	errorhandler.CheckNilErr(err)
	tpl = template.Must(tpl, err)
	return tpl
}

func PrintMultiSelectMessage(messages []string) error {
	var message, coloredMessage string
	var tpl *template.Template
	if len(messages) > 0 {
		var templateText string
		if len(messages) == 1 {
			templateText = fmt.Sprintf("%s %d option selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		} else {
			templateText = fmt.Sprintf("%s %d options selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		}
		for _, option := range messages {
			message = fmt.Sprintf("%s%s ", message, option)
		}
		coloredMessage = color.GreenString("%s", message)
		tpl = CreateTemplate("message", templateText)
	} else {
		message = "No options selected"
		coloredMessage = color.YellowString("%s", message)
		tpl = CreateTemplate("responseMessage", fmt.Sprintf("%s {{ . }}\n", constants.IconWarn))
	}
	err := tpl.Execute(os.Stdout, coloredMessage)
	return err
}

func PrintWarningMessage(message string) error {
	tpl := CreateTemplate("warningMessage", fmt.Sprintf("\n%s {{ . }}\n", constants.IconWarn))
	message = color.YellowString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}

func GetSuffixOfStack(stack, database string) string {
	var suffix string
	switch stack {
	case constants.ReactJS:
		suffix = constants.ReactTemplate
	case constants.NextJS:
		suffix = constants.NextTemplate
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeHapiPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeHapiMySqlTemplate
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeGraphqlPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeGraphqlMySqlTemplate
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			suffix = constants.NodeExpressMongoTemplate
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.GolangPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.GolangMySqlTemplate
		}
	case constants.ReactNative:
		suffix = constants.ReactNativeTemplate
	case constants.Android:
		suffix = constants.AndroidTemplate
	case constants.IOS:
		suffix = constants.IOSTemplate
	case constants.Flutter:
		suffix = constants.FlutterTemplate
	}
	return suffix
}

type StackDetails struct {
	Name      string
	Language  string
	Framework string
	Type      string
	Databases string
}

func GetStackDetails(service string) []StackDetails {
	var stacksDetails []StackDetails
	switch service {
	case constants.Backend:
		stacksDetails = []StackDetails{
			{
				Name:      constants.NodeHapiTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Hapi",
				Type:      "REST API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
			{
				Name:      constants.NodeExpressGraphqlTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Express",
				Type:      "GraphQL API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
			{
				Name:      constants.NodeExpressTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Express",
				Type:      "REST API",
				Databases: constants.MongoDB,
			},
			{
				Name:      constants.GolangEchoTemplate,
				Language:  "Golang",
				Framework: "Echo",
				Type:      "GraphQL API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
		}
	case constants.Web:
		stacksDetails = []StackDetails{
			{
				Name:      constants.ReactJS,
				Language:  "JavaScript",
				Framework: "React",
			},
			{
				Name:      constants.NextJS,
				Language:  "JavaScript",
				Framework: "Next.js",
			},
		}
	case constants.Mobile:
		stacksDetails = []StackDetails{
			{
				Name:      constants.ReactNative,
				Language:  "JavaScript",
				Framework: "React Native",
			},
			{
				Name:      constants.Android,
				Language:  "Kotlin",
				Framework: "-",
			},
			{
				Name:      constants.IOS,
				Language:  "Swift",
				Framework: "-",
			},
			{
				Name:      constants.Flutter,
				Language:  "Dart",
				Framework: "Flutter",
			},
		}
	}
	return stacksDetails
}

func FindConfigStacks(configLine string) []string {
	var stack string
	var stacks []string
	stackFound := false

	for _, char := range configLine {
		if char == '(' {
			stackFound = true
			continue
		} else if char == ')' {
			stackFound = false
			stacks = append(stacks, stack)
			stack = ""
		}
		if stackFound {
			stack = fmt.Sprintf("%s%s", stack, string(char))
		}
	}
	return stacks
}

func FindStacksByConfigStacks(configStacks []string) []string {
	var stacks []string

	_, _, directories := ExistingStacksDatabasesAndDirectories()
	camelCaseDirectories := ToCamelCase(directories)

	for _, configStack := range configStacks {
		for idx, camelCaseDirName := range camelCaseDirectories {
			if configStack == camelCaseDirName {
				stacks = append(stacks, directories[idx])
			}
		}
	}
	return stacks
}

func IsYarnOrNpmInstalled() string {
	var pkgManager string
	err := exec.Command(constants.Yarn, "-v").Run()
	if err != nil {
		err = exec.Command(constants.Npm, "-v").Run()
		if err != nil {
			errorhandler.CheckNilErr(fmt.Errorf("Please install 'yarn' or 'npm' in your machine."))
		} else {
			pkgManager = constants.Npm
		}
	} else {
		pkgManager = constants.Yarn
	}
	return pkgManager
}
