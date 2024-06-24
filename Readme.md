## Welcome to Frock. 

#### Easy tool to run your applications with Podman locally.

### Installation:
Go to Releases and install any version you want.  

### Json schema for yaml
Run `frock schema` to generate json schema for frock.yaml.
After that you can assign this schema to your frock.yaml file in your IDE.
Or use with any json schema validator.

### Yaml cheatsheet:
Create a file frock.yaml in the root of your project.
You also can create frock.override.yaml to override some settings.

```yaml
projectName: exampleBot # Name of the project
apps: # List of applications you want to run
  http: # Internal name of the application
    state: disabled # Default: "enabled" - Enable or disable application
    image: vladitot/php83-swow-ubuntu-local # Image name
    mountContainerDir: /var/www/ # Default: "/app" - Where to mount local directory to container.
    mountCode: enabled # Default: "enabled" - Enable mounting local directory to container. "enabled"
    tag: "latest" # Default: "latest" - Image tag
    privilegedMode: "enabled" # Default: "disabled" - Enable privileged mode
    labels: # Labels for container. Very usefully for exposing container to traefik load balancer
      traefikHost: "traefik.http.routers.example.rule=Host(`example.localhost`)"
      appPortToExpose: "traefik.http.services.example.loadbalancer.server.port=8080"
    command: ["/var/www/vendor/bin/gear-dev-server-arm64"] # Command to run in container. First element is entrypoint, the rest are arguments
boxes: # Boxes - containers that are not applications, but need to run your app. For example, database, redis, etc.
  postgresql: # settings are pretty same with apps, but without volume mount
    image: postgres
    tag: "13"
    state: enabled
    env:
      POSTGRES_DB: "example"
      POSTGRES_USER: "example"
      POSTGRES_PASSWORD: "example"
    labels:
    # any labels for container are here
    privilegedMode: disabled
    command:
      - postgres
commands: # List of usefully user defined commands
  - signature: "install-laravel" # Command will be called as "frock run install-laravel"
    type: container # Where to run the command. "container" or "local"
    command: # Command to run
      - composer
      - create-project
      - laravel/laravel
      - /var/www/php
  - signature: "bash"
    type: container
    command:
      - bash
      - "-c"
      - "-l"
      - bash
  - signature: "debug"
    type: container
    command: [ bash, "-c", "-l", "XDEBUG_SESSION=PHPSTORM XDEBUG_MODE=debug PHP_IDE_CONFIG=serverName=exampleBot bash" ] # you also can pass env variables like this
```

### Running your project
Run `frock up` to start all applications defined in frock.yaml.

### Stopping your project
Run `frock down` to stop all applications defined in frock.yaml.

### Running user defined command
Run `frock run <command_signature>` to run user defined command.

## Traefik
Frock has built-in traefik support. You can expose your applications to traefik load balancer by adding labels to your application.

To run traefik container you can call `frock upTraefik`. It will start traefik container with default settings. It will parse labels on other containers and expose them to you host machine.

To switch traefik down you can call `frock downTraefik`.

### Website
Website is under construction. We will provide more information soon.