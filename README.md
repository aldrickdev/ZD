# ZD

This is the Zendesk Service that is responsible for producing User Events.

This project makes use of [Microsoft Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers). This allows you to create reproducible development environments for all developers working on this project. The only requirements are [Docker](https://www.docker.com/), [Visual Studio Code](https://code.visualstudio.com/) and the **Dev Containers VSCode Extension**.

![devcontainer extension](docs/images/devcontainer-extension.jpg)

## Using the Devcontainer

Before we start, this application depends on 2 environment variables; `ENV` and `PORT`. These can be set by copying the `.env.example` file found in the `.devcontainer` directory and name the copy `.env`. Inside the `.env` file, provide the desired values for `ENV` and `PORT`.

Once you have provided the environment variables, open the `Command Palette`, type `Dev Containers` and select the command `Dev Containers: Reopen in Container`. This should trigger VSCode to build a Devcontainer with everything you need to start working on the project. Once the Devcontainer is created, VSCode will close and reopen with a remote connection to the Devcontainer. From here you have everything you need to work on the project.

Some key tools that were installed into this Devcontainer are:

- [Golang v1.21](https://github.com/devcontainers/images/tree/main/src/go)
- [NodeJS v18](https://github.com/devcontainers/features/tree/main/src/node)
- [Task](https://github.com/eitsupi/devcontainer-features/tree/main/src/go-task)
- [JSON Server](https://www.npmjs.com/package/json-server)
- [Air](https://github.com/cosmtrek/air/tree/master)

## Running the Application

Once your VSCode has connected to the Devcontainer, make sure you are in the directory `/workspaces` and run the command:

``` bash
task run
```

This will run the application but if you would like to the application to reload when you make changes then you will need to run the command:

``` bash
task dev-run
```
