# Project Roadmap (Go Fiber + SQLite + SQLC + dbmate-in-code)

## Phase 1: Concept and Environment Setup

Defining the foundation and initializing the Go module.

- [ ] **Data Model Definition**
  - [ ] Plan the Schema: User (id, email, password_hash, role) and 2-3 main entities.
  - [ ] Plan Relationships: Foreign keys for SQLite (e.g., `user_id` in items table).
- [ ] **Repository Initialization**
  - [X] `git init`.
  - [X] Create `.gitignore`:
    - [ ] Add `main` (binary).
    - [X] Add `*.db`, `*.db-journal`, `*.sqlite`.
    - [ ] Add `node_modules`, `.env`, `.output`, `dist`.
- [X] **Tooling Setup**
  - [X] Install **sqlc**: `brew install sqlc` (or go install).
  - [X] Install **dbmate** (CLI) only for generating files: `brew install dbmate`.
  - [X] Initialize Go Module: `go mod init <module-name>`.
  - [X] Install dependencies:
    - [X] Fiber: `go get github.com/gofiber/fiber/v2`.
    - [X] SQLite Driver: `go get github.com/mattn/go-sqlite3`.
    - [X] dbmate Library: `go get github.com/amacneil/dbmate/v2/pkg/dbmate`.

---

## Phase 2: Backend - Core and Database Setup

Setting up the programmatic migration runner and SQLC.

### 2.1 Database & Migrations (In-Code)

- [X] **Directory Setup**: Create a folder `db/migrations`.
- [X] **Migration Logic**:
  - [X] In `main.go` (or a dedicated `database` package), import `github.com/amacneil/dbmate/v2/pkg/dbmate`.
  - [X] Configure `dbmate.New` pointing to your SQLite file URL.
  - [X] Execute `dbmate.Migrate()` on application startup.
  - [X] Handle errors: If migration fails, the app should panic/exit.
- [X] **First Migration**:
  - [X] Run `dbmate new init_schema` (creates files in `db/migrations`).
  - [X] Write SQL DDL (CREATE TABLE Users...).
  - [X] Run `go run main.go` to verify the app applies the migration automatically.

### 2.2 SQLC Configuration

- [X] **Config**: Create `sqlc.yaml`.
  - [X] Set `path` to `internal/database` (or where you want generated code).
  - [X] Point `schema` to `db/migrations`.
  - [X] Point `queries` to `db/queries`.
- [ ] **User Queries**:
  - [ ] Create `db/queries/users.sql`.
  - [ ] Write: `CreateUser`, `GetUserByEmail`, `GetUserById`.
  - [ ] Run `sqlc generate`.

### 2.3 Authentication Logic

- [X] **Server**: Initialize `fiber.New()`.
- [X] **DB Connection**: Open `sql.Open` (separate from dbmate, for app usage) and pass to SQLC `New()`.
- [ ] **Endpoints**:
  - [ ] `POST /auth/register`: Validate -> Hash -> `queries.CreateUser`.
  - [ ] `POST /auth/login`: `queries.GetUserByEmail` -> Check Hash -> Generate JWT.
- [ ] **Middleware**: JWT validation parsing `Authorization` header.

---

## Phase 3: Backend - Business Logic (CRUD)

Standard workflow: Migration -> Query -> Generate -> Handler.

- [ ] **Entity 1 (e.g., Items)**
  - [ ] **Migration**: `dbmate new create_items_table`.
  - [ ] **Define**: Add SQL DDL to the new file.
  - [ ] **Restart App**: Verify auto-migration works.
  - [ ] **Queries**: Create `db/queries/items.sql`.
    - [ ] `-- name: CreateItem :one`
    - [ ] `-- name: ListItems :many`
    - [ ] `-- name: GetItem :one`
    - [ ] `-- name: UpdateItem :one`
    - [ ] `-- name: DeleteItem :exec`
  - [ ] **Generate**: Run `sqlc generate`.
  - [ ] **Handlers**: Create Fiber handlers using the generated code.
- [ ] **Entity 2**:
  - [ ] Repeat the cycle.

---

## Phase 4: Frontend - Nuxt SPA Setup

Visual layer setup.

- [ ] **Nuxt Initialization**:
  - [ ] `npx nuxi@latest init <project-name>`.
  - [ ] Set `ssr: false` in `nuxt.config.ts`.
- [ ] **Dependencies**: Install TailwindCSS and Pinia.
- [ ] **API Client**: Create a composable using `$fetch` pointing to localhost.

---

## Phase 5: Frontend - Integration

### 5.1 Auth Integration

- [ ] **Store**: Pinia store for User state.
- [ ] **Login/Register UI**: Forms connected to API.
- [ ] **Persistance**: Store token in Cookie/LocalStorage.

### 5.2 CRUD UI

- [ ] **Lists**: Fetch data on `onMounted`.
- [ ] **Actions**: Create/Update/Delete buttons connected to API.
- [ ] **Feedback**: Toast notifications for success/error.

---

## Phase 6: Hardening

### Backend (Go)

- [ ] **Validation**: Integrate `go-playground/validator` for request body validation.
- [ ] **SQLite Tuning**:
  - [ ] Execute `PRAGMA journal_mode=WAL;` after opening the DB connection.
  - [ ] Execute `PRAGMA foreign_keys=ON;` (Critical for SQLite constraints).
- [ ] **Security**: Helmet headers, Rate Limiting, CORS.
- [ ] **Logging**: Structured logging (slog) for requests and DB errors.

### Frontend

- [ ] **Types**: TypeScript interfaces matching Go structs.
- [ ] **Error Handling**: Global handling for 401 Unauthorized.

---

## Phase 7: Deployment

- [ ] **Embed Migrations (Optional but Recommended)**:
  - [ ] Use Go `embed` directive to bundle the `db/migrations` folder into the binary.
  - [ ] Update `dbmate` configuration in `main.go` to read from the embedded `fs.FS` instead of the file system.
- [ ] **Docker**:
  - [ ] Build the Go binary.
  - [ ] Setup volume for the SQLite file (`/app/data`).
  - [ ] Since migrations run in code, the container just needs to start the binary.
