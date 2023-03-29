package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateDockerComposeFile(stackInfo map[string]interface{}, forceCreate bool) error {

	filePath := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(),
		constants.DockerComposeFile,
	)
	if !forceCreate {
		status, _ := fileutils.IsExists(filePath)
		if status {
			return errorhandler.ErrExist
		}
	}
	// Don't make any changes in the below source string.
	source := `version: '3'
services:
{{#if backendStatus}}
  # Setup {{database}}
  {{dbServiceName stack database}}:
    image: '{{dbVersion database}}' 
    ports:
      - {{portConnection database}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - {{envFileBackend database}}
    volumes:
      - {{projectName}}_db_volume:/var/lib/{{databaseVolume database}}
{{#equal stack GolangPostgreSQL}}
    environment:
      POSTGRES_USER: ${PSQL_USER}
      POSTGRES_PASSWORD: ${PSQL_PASS}
      POSTGRES_DB: ${PSQL_DBNAME}
      POSTGRES_PORT: ${PSQL_PORT}
{{/equal}}
{{#equal stack GolangMySQL}}
    environment:
      MYSQL_DATABASE: ${MYSQL_DBNAME}
      MYSQL_PASSWORD: ${MYSQL_PASS}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
{{/equal}}

  # Setup Redis
  redis:
    image: 'redis:6-alpine'
    ports:
      - {{portConnection redis}}
    # Default command that redis will execute at start
    command: ['redis-server']

{{#equal stack GolangPostgreSQL}}
{{{waitForDBService database}}}

{{/equal}}
{{#equal stack GolangMySQL}}
{{{waitForDBService database}}}

{{/equal}}
  # Setup {{projectName}} API
  {{projectName}}_api:
    build:
      context: './backend'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - {{envFileBackend database}}
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
{{/if}}
{{#if webStatus}} 
  # Setup {{projectName}} web
  {{projectName}}_web:
    build:
      context: './web'
    ports:
      - {{portConnection web}}
    env_file:
      - ./web/.env.docker
{{/if}}

# Setup Volumes
volumes:
  {{projectName}}_db_volume:
`

	err := hbs.ParseAndWriteToFile(source, filePath, stackInfo)
	errorhandler.CheckNilErr(err)

	return nil
}
