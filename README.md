# Gin API Template

This is a template that can be used to get a running Golang Gin API quickly.
This repo takes advantage of [Microsoft Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
to create a reproducible environment for any developer. All you need is [Docker](https://www.docker.com/), the Dev 
Containers VSCode Extension and the [Dev Container CLI](https://github.com/devcontainers/cli).

## Running the Devcontainer

To start the Devcontainer copy the `.devcontainer/.env.example`, name the copy `.env` and leave it in the
`.devcontainer` directory. Provide a value for `ENV` and `PORT`.

- `ENV` tells the application what environment it is running in, `PROD` would run Gin in Release Mode
- `PORT` tells the application what port it should use to receive requests

After the environment variables, open the `Makefile` and set the `MOD_NAME` variables, this will be the name of the 
the application. After setting that you can install all the dependencies with the command:

```bash
make init
```

After setting that up you can start up the Dev Container with the command:

```bash
make dc_up
```

Once it has finished, you can open a VSCode instance in that container with the command:

```bash
make dc_open
```

To run the application, run the command:

```bash
make run
```
