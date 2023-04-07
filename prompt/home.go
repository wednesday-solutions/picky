package prompt

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func PromptHome() {
	var p PromptInput
	p.Label = "Pick an option"
	p.GoBack = PromptHome
	var initService bool
	stacks, databases, _ := utils.ExistingStacksDatabasesAndDirectories()
	if len(stacks) > 0 {
		p.Items = []string{constants.InitService, constants.CreateCD}
		showCreateDC := ShowCreateDockerCompose(databases)
		if showCreateDC {
			p.Items = append(p.Items, constants.DockerCompose)
		}
		p.Items = append(p.Items,
			constants.SetupInfra,
			constants.Deploy,
		)
		showRemoveDeploy := ShowRemoveDeploy()
		if showRemoveDeploy {
			p.Items = append(p.Items, constants.RemoveDeploy)
		}
		p.Items = append(p.Items, constants.Exit)
		response := p.PromptSelect()
		switch response {
		case constants.InitService:
			initService = true
		case constants.DockerCompose:
			PromptDockerCompose()
		case constants.CreateCD:
			PromptCreateCD()
		case constants.SetupInfra:
			PromptSetupInfra()
		case constants.Deploy:
			PromptDeploy()
		case constants.RemoveDeploy:
			PromptRemoveDeploy()
		case constants.Exit:
			PromptExit()
		}
	}
	if len(stacks) == 0 || initService {
		PromptSelectService()
	}
}