# Archon
Archon is a simple configuration management tool comprised of the archon agent and archon controller.

The name was inspired by the [Dark Archon](https://starcraft.fandom.com/wiki/Dark_archon_(StarCraft)) in StarCraft and their mind control power.

## Agent
The archon agent runs on each server in the fleet and accepts http calls to perform tasks. It's possible to control the agent directly through api calls and its endpoints are abstracted out to allow for future implementations on additional operating systems.

## Controller
The archon controller parses instructions from the configuration files and concurrently runs these tasks over the fleet. These configuration tasks are just http calls in disguise allowing for homogeny in their structure and a one to one relation with their underlying api calls.

## Configuration
The first configuration file type is the archon script yaml, it contains the tasks to be executed over the remote fleet. These tasks are executed sequentially and a task failure stops script execution on the given server. Tasks mimic the underlying api calls with `resource` being the path, `action` being the method, `parameters` being the query values, and `payload` being the file used as the request body. An example is provided below.
```
- resource: package
  action: post
  parameters:
    package: apache2

- resource: service
  action: delete
  parameters:
    service: apache2

- resource: file
  action: post
  parameters:
    file: /var/www/html/index.html
    mode: "644"
  payload: hello-world.html
```

The second configuration file type is the archon fleet yaml, it contains the list of servers in the fleet. These servers can either be ip or dns addresses. An example is provided below.
```
- ec2-12-34-567-890.us-west-2.compute.amazonaws.com
- ec2-12-34-567-980.us-west-2.compute.amazonaws.com
- ec2-12-34-567-809.us-west-2.compute.amazonaws.com
```

## Linux
Precompiled Linux x86-64 binaries are available under the linux directory.

## Variables
The configuration of `agent` and `archon` are set with the following environment variables.

- TOKEN: the token used for request authentication, can be anything but must match; required
- PORT: the port used for api communication; defaults to 8080

## Installation
Before a script can be executed over a fleet the archon agent must be installed on each remote server. To do this one must run the `agent` process on the remote servers in a detached state under the root user.

Ensure the `TOKEN` environment variable is set before running.

## Execution
To execute an archon script over a fleet run the following command.
```
archon <example.script.yaml> <example.fleet.yaml>
```

(note: the payload paths in the script are set relative to the working directory)

Ensure the `TOKEN` environment variable is set before running.