---
name: server-middleware
description: >-
  Runs project middleware on the user's Tencent Cloud server with Docker
  instead of local Docker Desktop. Use when a project needs MySQL, Redis,
  RabbitMQ, or other middleware, when local Docker will not start, or when
  the user mentions the server, 124.221.130.183, middleware host, or SSH host
  middleware.
---

# Server middleware

Local Docker Desktop is not the place to run project databases and queues. Use the server below.

## Connection

```text
Host middleware
HostName 124.221.130.183
User ubuntu
Port 22
IdentityFile ~/.ssh/id_ed25519
```

On a new computer, add that host to `~/.ssh/config` and use the same private key. Do not copy the private key into this repo.

```powershell
ssh middleware
```

`ubuntu` has passwordless `sudo`. Docker commands on the server use `sudo docker`.

`.env` stays gitignored. This repo's `.env.example` is the tracked connection file. `config.Load` reads `.env` first and falls back to `.env.example`, so another computer can run after `git pull`.

## What is already running

- OS: Ubuntu 22.04. Docker Engine 26.1.3 is installed, enabled, and active.
- Image pulls use `https://mirror.ccs.tencentyun.com` in `/etc/docker/daemon.json`. Leave that mirror in place while pulls succeed.
- Container `mysql-33603` (`mysql:8.0`) publishes `33603 -> 3306`. From the dev PC, `124.221.130.183:33603` is reachable. Do not recreate or delete that container.
- GO-Demo uses database `awesome_project` on that instance. Tables `todos` and `todo_status_logs` are already created. The DSN is in `.env.example`. Port `13306` is not open on the cloud firewall.
- UFW is inactive. A new host port still needs to be open in the Tencent Cloud security group if the PC cannot connect.

## New middleware

Give each new project its own container and host port. GO-Demo already has its database on `mysql-33603`; do not add another project's tables there.

1. `ssh middleware` and confirm `sudo docker info` works.
2. Pick a free host port. Check with `sudo ss -lnt`, then confirm the dev PC can open that port.
3. Start the container with Docker, publish only the needed port, and keep data in a named volume.
4. Put the connection string in `.env.example` for this learning repo, and in local `.env` when overriding it. Point the app at `124.221.130.183` and the published port.
5. Apply the project's migrations on that new database.
6. From the project machine, test the TCP port, then run the project's tests.

If an image pull fails, check `sudo docker info` for the registry mirror before changing `/etc/docker/daemon.json`. Do not remove the Tencent mirror until a replacement pull succeeds.
