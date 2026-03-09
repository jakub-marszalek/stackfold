# Stackfold

```
    _____ __             __   ____      __    __ 
   / ___// /_____ ______/ /__/ __/___  / /___/ / 
   \__ \/ __/ __ '/ ___/ //_/ /_/ __ \/ / __  /  
  ___/ / /_/ /_/ / /__/ ,< / __/ /_/ / / /_/ /   
 /____/\__/\__,_/\___/_/|_/_/  \____/_/\__,_/
```

**Stackfold** is a CLI automation tool that streamlines the setup, build, and management of new Laravel projects using Docker. It handles container orchestration, environment configuration, health checks, and database operations — so you can start a new project in minutes.

---

## Requirements

- [Docker](https://www.docker.com/) with Docker Compose

---

## Installation

```bash
git clone https://gitlab.com/hesoyamcc/stackr.git
cd stackr
```

The compiled binary is included in `bin/stackfold`. Make it executable and move it to a directory on your `PATH`:

```bash
chmod +x bin/stackfold
sudo mv bin/stackfold /usr/local/bin/stackfold
```

---

## Quickstart

```bash
# 1. Initialize a new project (runs the configuration wizard)
stackfold init

# 2. Build Docker images and start all containers
stackfold build --start

# 3. Open your browser at http://localhost:8080
```

That's it. Stackfold installs a fresh Laravel project, configures the environment, and starts all services automatically.

---

## Commands

### `stackfold init`
Guides you through first-time project setup. Runs an interactive wizard to collect project name, database credentials, and port configuration. Installs a fresh Laravel project via the official Composer Docker image and patches `vite.config.js` for Docker HMR support.

```bash
stackfold init
```

### `stackfold build`
Builds Docker images from scratch using `--no-cache`. Optionally starts the stack immediately after a successful build.

```bash
stackfold build           # build only
stackfold build --start   # build and start
stackfold build -s        # shorthand
```

### `stackfold up`
Starts all containers and runs a health check sequence to verify all services are ready.

```bash
stackfold up
```

### `stackfold down`
Stops and removes all containers. By default, volumes are preserved so your database data persists.

```bash
stackfold down             # stop, keep volumes
stackfold down --remove    # stop and destroy volumes
stackfold down -r          # shorthand
```

> **Warning:** `--remove` permanently deletes all database data. Run `stackfold backup` first.

### `stackfold test`
Runs a full health check of the stack. Verifies container status, Nginx HTTP response, PHP internals, Vite dev server, and open ports.

```bash
stackfold test
```

### `stackfold config`
Updates project name and database credentials in `.env.docker` and `app/.env`. Intended for reconfiguring an existing environment.

```bash
stackfold config
```

> **Note:** Changing database credentials does not migrate existing data. If you change the database name or password, run `stackfold reset --force` afterward.

### `stackfold reset`
Performs a full environment reset: shuts down the stack, removes all volumes, rebuilds images, and starts all services. **All database data will be lost.**

```bash
stackfold reset --force
stackfold reset -f
```

### `stackfold backup`
Dumps the current PostgreSQL database to a timestamped `.sql` file in the `backups/` directory.

```bash
stackfold backup
# creates: backups/backup_2026-01-15_14-30-00.sql
```

### `stackfold restore`
Lists available backups and lets you choose which one to restore interactively. Requires `--force` to prevent accidental overwrites.

```bash
stackfold restore --force
stackfold restore -f
```

### `stackfold --version`

```bash
stackfold --version
```

---

## Stack

| Service     | Image                | Default Port |
|-------------|----------------------|--------------|
| Nginx       | `nginx:alpine`       | 8080         |
| PHP-FPM     | `php:8.4-fpm-alpine` | —            |
| PostgreSQL  | `postgres:15`        | 5432         |
| Redis       | `redis:7-alpine`     | 6379         |
| Node / Vite | `node:22-alpine`     | 5173         |

---

## Project Structure

```
stackfold/
├── app/                  # Laravel application (created by stackfold init)
├── backups/              # Database backups (created by stackfold backup)
├── bin/                  # Compiled CLI binary
├── cli/                  # Go CLI source code
│   ├── cmd/
│   └── main.go
├── docker/               # Dockerfiles and config
│   ├── nginx/
│   ├── node/
│   └── php/
├── scripts/              # Bash stack scripts
│   └── bash/
├── src/                  # PHP health check source
├── docker-compose.yaml
└── .env.docker           # Docker environment config (created by stackfold init)
```

---

## License

MIT © 2026 Jakub Marszałek