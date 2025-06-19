<h1 align="center">go_api_demo (v2)</h1>

<p align="center">
  The ultimate single-file demonstration of a production-ready, idiomatic HTTP API server built with pure Go.
  <br />
  It showcases CRUD operations, repository pattern, middleware, validation, and concurrent in-memory storage.
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Language-Go-blue.svg" alt="Go"></a>
  <a href="https://github.com/dunamismax/golang/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://github.com/dunamismax/golang/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square" alt="PRs Welcome"></a>
  <a href="https://github.com/dunamismax/golang/stargazers"><img src="https://img.shields.io/github/stars/dunamismax/golang?style=social" alt="GitHub Stars"></a>
</p>

---

## ✨ Features & Philosophy

`go_api_demo` is a comprehensive showcase of modern Go backend development patterns, engineered for clarity and correctness in a single, self-contained file.

- **Standard Library First**: The foundation is Go's robust standard library, including `net/http` for the server, `log/slog` for structured logging, and `sync` for concurrency control.
- **Full CRUD API**: Implements complete Create, Read, Update, and Delete operations for a `User` resource, demonstrating RESTful principles with versioned endpoints (`/api/v1/...`).
- **Repository Pattern**: Decouples business logic from the data layer using a `UserRepository` interface. This includes a concurrent-safe, in-memory implementation that mimics a real database with a `sync.RWMutex`.
- **Middleware**: Features a `loggingMiddleware` to demonstrate how to handle cross-cutting concerns like logging every incoming request's method, path, and duration.
- **Struct Validation**: Employs `validator/v10` to enforce strict validation rules on incoming JSON request bodies, a critical practice for API security and data integrity.
- **Graceful Shutdown**: The server listens for OS signals (`SIGINT`, `SIGTERM`) and performs a graceful shutdown, allowing in-flight requests to complete before exiting.
- **Dependency Injection**: Utilizes a central `application` struct to hold and inject dependencies (like the logger and repository) into handlers, promoting clean, testable code.
- **Advanced Routing**: Leverages the Go 1.22 `http.ServeMux` to handle RESTful path parameters (e.g., `/users/{id}`) without needing a third-party router.

---

## 🚀 Getting Started

You will need the Go toolchain installed (version 1.22+ recommended).

### 1. Set Up The Project

Navigate to the project's root directory (`go_api_demo/`), then initialize the Go module and fetch the required dependency.

```sh
# Navigate to your project folder
cd go_api_demo

# Initialize the module
go mod init go_api_demo

# Get the validator dependency
go get github.com/go-playground/validator/v10

# Tidy up the module file
go mod tidy
```

### 2. Run the Server

Run the application from the project's root directory. By default, it will listen on port `8080`.

```sh
go run main.go
```

The server will log that it has started:

```json
{
  "time": "2025-06-19T22:30:00.123Z",
  "level": "INFO",
  "msg": "server starting",
  "address": ":8080"
}
```

---

## 🔬 Interacting with the API: A Guided Tour

The following is a complete walkthrough of the API's functionality using `curl`. Open a new terminal to run these commands.

### Step 1: Get All Users (Initial State)

First, let's see the list of users. Since we just started the server, it will be empty.

```sh
curl http://localhost:8080/api/v1/users
```

**Response:** An empty list.

```json
{ "status": "success", "data": [] }
```

### Step 2: Create a User (Success)

Now, let's create our first user, "Alice".

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}' \
  http://localhost:8080/api/v1/users
```

**Response:** The server confirms the creation and returns the new user object, complete with a server-assigned `id` and `createdAt` timestamp.

```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "user_1718843400000000000",
    "createdAt": "2025-06-19T23:10:00.00Z",
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

**Important:** For the next steps, copy the `id` from the response. We'll store it in a shell variable to make things easier.

```sh
# Replace the value with the actual ID you received
ALICE_ID="user_1718843400000000000"
```

### Step 3: Test Validation and Error Handling

Our API validates incoming data. Let's see what happens when we send invalid requests.

**A. Missing Required Field (`email`)**

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "Bad Request"}' \
  http://localhost:8080/api/v1/users
```

**Response:** A `400 Bad Request` with a clear error message.

```json
{
  "error": "validation failed: Key: 'input.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

**B. Invalid Email Format**

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "Bob", "email": "bob@invalid"}' \
  http://localhost:8080/api/v1/users
```

**Response:** The API rejects the invalid email address.

```json
{
  "error": "validation failed: Key: 'input.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**C. Name Too Short**

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "A", "email": "a@example.com"}' \
  http://localhost:8080/api/v1/users
```

**Response:** The API enforces the minimum length for the `name` field.

```json
{
  "error": "validation failed: Key: 'input.Name' Error:Field validation for 'Name' failed on the 'min' tag"
}
```

### Step 4: Create a Second User

Let's create another valid user, "Bob".

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "Bob", "email": "bob@example.com"}' \
  http://localhost:8080/api/v1/users
```

**Response:** Success. Note the new unique ID.

```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "user_1718843460000000000"
    /* ... */
  }
}
```

Again, let's save this new ID.

```sh
# Replace the value with the actual ID you received for Bob
BOB_ID="user_1718843460000000000"
```

### Step 5: Get All Users (Populated)

Now if we get the list of all users, we should see both Alice and Bob.

```sh
curl http://localhost:8080/api/v1/users
```

**Response:** A list containing two user objects.

```json
{
  "status": "success",
  "data": [
    {
      "id": "user_1718843400000000000",
      "name": "Alice"
      /* ... */
    },
    {
      "id": "user_1718843460000000000",
      "name": "Bob"
      /* ... */
    }
  ]
}
```

### Step 6: Get a Specific User

Let's retrieve only Alice's details using her ID.

```sh
curl http://localhost:8080/api/v1/users/$ALICE_ID
```

**Response:** A single user object for Alice.

```json
{ "status": "success", "data": { "id": "user_1718843400000000000" /* ... */ } }
```

### Step 7: Update a User

Let's change Alice's email address using a `PUT` request.

```sh
curl -X PUT -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice.smith@example.com"}' \
  http://localhost:8080/api/v1/users/$ALICE_ID
```

**Response:** Success, with the updated user object returned.

```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": "user_1718843400000000000",
    "name": "Alice",
    "email": "alice.smith@example.com"
  }
}
```

### Step 8: Delete a User

Now, let's delete Bob from the system.

```sh
curl -X DELETE http://localhost:8080/api/v1/users/$BOB_ID
```

**Response:** A confirmation message.

```json
{ "status": "success", "message": "User deleted successfully" }
```

### Step 9: Verify Deletion

If we get all users again, only Alice should remain.

```sh
curl http://localhost:8080/api/v1/users
```

**Response:** The list now contains only one user.

```json
{
  "status": "success",
  "data": [{ "id": "user_1718843400000000000" /* ... */ }]
}
```

### Step 10: Attempt to Get a Non-Existent User

Trying to `GET` or `DELETE` a user that no longer exists (like Bob) will result in an error.

```sh
curl http://localhost:8080/api/v1/users/$BOB_ID
```

**Response:** A `404 Not Found` error.

```json
{ "error": "user not found" }
```

### Step 11: Graceful Shutdown

To stop the server, return to the terminal where it's running and press `Ctrl+C`. You will see shutdown logs as the server gracefully terminates.

---

## 🤝 Contribute

**This project is built by the community, for the community. We need your help!**

If this project excites you, please **give it a star!** ⭐ It helps us gain visibility and attract more talented contributors like you.

### Connect

Connect with the author, **dunamismax**, on:

- **Twitter:** [@dunamismax](https://twitter.com/dunamismax)
- **Bluesky:** [@dunamismax.bsky.social](https://bsky.app/profile/dunamismax.bsky.social)
- **Reddit:** [u/dunamismax](https://www.reddit.com/user/dunamismax)
- **Discord:** `dunamismax`
- **Signal:** `dunamismax.66`

## 📜 License

This project is licensed under the **MIT License**. See the `LICENSE` file for details.
