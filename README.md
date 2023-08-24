# STAUVISE

# GO - An Audio/Video Streaming Service
This service implements these main features:

1. Stream media files.
2. Create and manage users.
3. Create and manage categories.
4. Create and manage videos.

Supported transport protocols
- [x] HLS
- [ ] RTMP
- [ ] FLV

Infrastructure:
- App context: [AppContext](https://github.com/hoangtk0100/app-context) (library packages common components)
- ffmpeg: `ffmpeg`
- Token authentication: `PASETO`
- Database migrations: `golang-migrate`
- Generate CRUD operations from SQL: `GORM`
- Database: `PostgreSQL`
- Database documentation: `dbdocs`
- Generate SQL schema: `dbml2sql`
- Web framework: `Gin`
- Containerize: `docker` | `docker compose`


## 1. Quick start
- Setup ENV:

    ```bash
    make build
    make outenv     # Show ENV on terminal
    make outenvfile # Extract ENV into .env

    # Replace ENV variables in `.env` with your configuration
    ```

- Run in background (all in one):

    ```bash
    make up
    ```

- Run on current terminal:

    ```bash
    # Start PostgreSQL database
    make updb

    # Run server
    make server
    ```

## 2. Setup local development

### Install tools
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```

- [ffmpeg](https://www.geekbits.io/how-to-install-ffmpeg-on-mac-os)

    ```bash
    brew install ffmpeg
    ```

- [Swagger](https://github.com/swaggo/gin-swagger)

    ``` bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

### Setup infrastructure

- Start PostgreSQL container:

    ```bash
    make updb
    ```

- Run DB migration up all versions:

    ```bash
    make migrateup
    ```

- Run DB migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run DB migration down all versions:

    ```bash
    make migratedown
    ```

- Run DB migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DB documentation:

    ```bash
    make dbdocs
    ```

- Access the DB documentation at [THIS ADDRESS](https://dbdocs.io/hoangtk.0100/stauvise).

### How to generate code

- List ENV variables on terminal:

    ```bash
    make outenv
    ```

- Extract ENV variables into `.env` file:

    ```bash
    make outenvfile
    ```

- Generate schema SQL file with DBML:

    ```bash
    make dbschema
    ```

- Create a new DB migration:

    ```bash
    make new_migration name=<migration_name>
    ```

- Generate swagger documents:

    ```bash
    make swagger
    ```

### How to run

- Run server on terminal:

    ```bash
    make server
    ```

- Run containers (verbose mode):

    ```bash
    make upv
    ```

- Run containers (detach mode):

    ```bash
    make up
    ```

- Remove containers:

    ```bash
    make down
    ```


## 3. API Endpoints - [Swagger](http://localhost:8080/api/v1/swagger/index.html)
| Method | Path                                    | Description                   | Notes                             |
| ------ | --------------------------------------- | ----------------------------- | --------------------------------- |
| POST   | `/api/v1/users/register`                    | Register a new user           | Receives user details             |
| POST   | `/api/v1/auth/login`                        | Login to get access token     | Receives username and password   |
| GET    | `/api/v1/users/me`                          | Retrieve user profile         | Requires valid access token      |
| POST   | `/api/v1/categories`                        | Create a new category         | Receives category details        |
| GET    | `/api/v1/categories`                        | Retrieve all categories       | Requires valid access token      |
| GET    | `/api/v1/videos/streams/{id}/{id}.m3u8`      | Stream video segments         | Requires valid video ID          |
| GET    | `/api/v1/videos`                            | Retrieve videos               | Supports pagination              |
| GET    | `/api/v1/videos/{id}`                       | Retrieve video by ID          | Requires valid video ID          |
| GET    | `/api/v1/videos/{id}/segments`              | Retrieve video segments       | Requires valid video ID, supports pagination |
| POST   | `/api/v1/videos`                            | Create a new video            | Receives video details and file  |


## 4. Notes
- To test streaming media use: [HLS Demo](https://hlsjs.video-dev.org/demo)
- Input format: `{host}/api/v1/videos/streams/{video_path}` (video_path: "path" in getting video details response)
- For example: `http://localhost:8080/api/v1/videos/streams/1692794790339777000/1692794790339777000.m3u8`