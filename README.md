# spyglass
Spyglass: Magnifying the status of almost everything

## Local testing

Before running the server, make sure you have [Docker](www.docker.com) and [Docker Compose](www.docker.com/compose) installed. Run then the following command, from the root dir of the repo:

```
docker-compose up -d
```

If you want to run the API, use:

```
make api
```

For running unit tests:

```
make test
```

## E2E

Inside the `api/e2e` dir, there are several `.http` files that can be opened in Visual Studio Code,
and you can execute manual HTTP request tests using the `humao.rest-client` VS Code addon.
