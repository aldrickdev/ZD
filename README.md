# ZD

This is the Zendesk Service that is responsible for producing User Events.

This project makes use of [Microsoft Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers). This allows you to create reproducible development environments for all developers working on this project. The only requirements are [Docker](https://www.docker.com/), [Visual Studio Code](https://code.visualstudio.com/) and the **Dev Containers VSCode Extension**.

![devcontainer extension](docs/images/devcontainer-extension.jpg)

## Using the Devcontainer

Before we start, this application depends on 2 environment variables; `ENV` and `PORT`. These can be set by copying the `.env.example` file found in the `.devcontainer` directory and name the copy `.env`. Inside the `.env` file, provide the desired values for `ENV` and `PORT`.

Once you have provided the environment variables, open the `Command Palette`, type `Dev Containers` and select the command `Dev Containers: Reopen in Container`. This should trigger VSCode to build a Devcontainer with everything you need to start working on the project. Once the Devcontainer is created, VSCode will close and reopen with a remote connection to the Devcontainer. From here you have everything you need to work on the project. If you need to close the remote connection between VSCode and the Devcontainer open the `Command Palette` and select the command `Dev Containers: Reopen Folder Locally`.

Some key tools that were installed into this Devcontainer are:

- [Golang v1.21](https://github.com/devcontainers/images/tree/main/src/go)
- [NodeJS v18](https://github.com/devcontainers/features/tree/main/src/node)
- [Task](https://github.com/eitsupi/devcontainer-features/tree/main/src/go-task)
- [JSON Server](https://www.npmjs.com/package/json-server)
- [Air](https://github.com/cosmtrek/air/tree/master)

## Running the Zendesk Service

Before continuing, make sure that you are in VSCode, have a remote connection to the Devcontainer and are in the `/workspaces` directory. Getting a remote connection to the Devcontainer is covered in section [Using the Devcontainer](#using-the-devcontainer).

The application has a dependency on the User Service which currently, is responsible for providing all available users and event data. Currently, the User Service is under development so to mock this User Service, we are using [JSON Server](https://www.npmjs.com/package/json-server). 

To run the Mocked User Service, run the command below, this will startup the Mocked User Service locally so that the Zendesk Service can get the required data for it to run.

``` bash
task run-us
```

To run the Zendesk Service, run the command below. This will run the Zendesk Service locally.

``` bash
task run-zd
```

Since running the Zendesk Service depends on the Mock User Service, both need to be running. Since both the Zendesk Service and the Mock User Service hijacks the terminal you will need separate terminals to run them both with the commands above. So the commands above are good if you would like to run them independently, however if you would like to run them both at the same time, without needing to create another terminal instance, you can user the `--parallel` flag that Task provides. The command below will run both of the tasks `run-us` and `run-zd` at the same time, in the same terminal.

``` bash
task --parallel run-us run-zd
```

Note that VSCode will automatically port forward the ports that both the Zendesk Service and the Mock User Service expose, to your local machine so that you can access them from outside of the Devcontainer.

If you would like to see if messages are actually being added to the queue, you can also run the consumer app located in `cmd/consumer`. This is a very simple Go application that receives messages from the queue so that you can see that it is working.
