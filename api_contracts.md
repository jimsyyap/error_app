To ensure alignment between backend and frontend development for the **Tennis Error Tracker** web application, I’ve defined a comprehensive set of API contracts below. These contracts specify the endpoints, HTTP methods, request payloads, and response formats needed to support the app’s functionality. The app allows users to track unforced errors during tennis matches, requiring features like user authentication, session management, error logging, and summary reporting. Here’s the complete API specification:

---

## API Contracts for Tennis Error Tracker

### General Notes
- **Authentication**: Most endpoints require a JSON Web Token (JWT) included in the `Authorization` header as `Bearer <token>`. Exceptions are `/register` and `/login`.
- **Timestamps**: All timestamps are in UTC and returned in ISO 8601 format (e.g., `"2023-10-05T14:48:00Z"`).
- **IDs**:
  - Session IDs (`session_id`) are UUIDs (strings).
  - Error type IDs (`error_type_id`) are integers.
- **Error Responses**: Include a JSON object with an `"error"` field (e.g., `{ "error": "Invalid credentials" }`).

---

### 1. Authentication Endpoints

#### **POST /register**
- **Description**: Register a new user.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Responses**:
  - `201 Created`: User registered successfully.
  - `400 Bad Request`: Invalid input (e.g., missing or malformed fields).
    ```json
    {
      "error": "Missing username or password"
    }
    ```
  - `409 Conflict`: Username already exists.
    ```json
    {
      "error": "Username already exists"
    }
    ```

#### **POST /login**
- **Description**: Authenticate a user and return a JWT token.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Responses**:
  - `200 OK`: Authentication successful.
    ```json
    {
      "token": "jwt_token_string"
    }
    ```
  - `401 Unauthorized`: Invalid credentials.
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

---

### 2. Match Session Endpoints

#### **POST /sessions**
- **Description**: Start a new match session.
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: None
- **Responses**:
  - `201 Created`: Session started.
    ```json
    {
      "session_id": "uuid"
    }
    ```
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```

#### **PUT /sessions/{session_id}**
- **Description**: End an active match session.
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: None
- **Responses**:
  - `200 OK`: Session ended successfully.
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - `404 Not Found`: Session not found.
    ```json
    {
      "error": "Session not found"
    }
    ```

#### **GET /sessions**
- **Description**: Retrieve a list of the user’s past match sessions.
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: List of sessions.
    ```json
    [
      {
        "session_id": "uuid",
        "start_time": "2023-10-05T14:48:00Z",
        "end_time": "2023-10-05T15:30:00Z"
      },
      ...
    ]
    ```
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```

#### **GET /sessions/active**
- **Description**: Retrieve the user’s current active session, if any.
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: Active session found.
    ```json
    {
      "session_id": "uuid",
      "start_time": "2023-10-05T14:48:00Z"
    }
    ```
  - `204 No Content`: No active session.
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```

---

### 3. Error Logging Endpoints

#### **POST /errors**
- **Description**: Log an unforced error in a session. (Timestamp is set server-side.)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "session_id": "uuid",
    "error_type_id": 1  // Integer (e.g., 1 = "Forehand")
  }
  ```
- **Responses**:
  - `201 Created`: Error logged successfully.
  - `400 Bad Request`: Invalid input (e.g., missing fields).
    ```json
    {
      "error": "Missing session_id or error_type_id"
    }
    ```
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - `404 Not Found`: Session not found.
    ```json
    {
      "error": "Session not found"
    }
    ```

#### **DELETE /errors/last**
- **Description**: Undo the last error logged in a session.
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "session_id": "uuid"
  }
  ```
- **Responses**:
  - `200 OK`: Last error deleted.
  - `400 Bad Request`: No errors to undo.
    ```json
    {
      "error": "No errors to undo in this session"
    }
    ```
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - `404 Not Found`: Session not found.
    ```json
    {
      "error": "Session not found"
    }
    ```

---

### 4. Summary Endpoints

#### **GET /sessions/{session_id}/summary**
- **Description**: Retrieve an error summary for a specific session.
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: Summary of errors.
    ```json
    {
      "total_errors": 10,
      "errors_by_type": {
        "Forehand": 3,
        "Backhand": 2,
        "Serve": 4,
        "Volley": 1
      }
    }
    ```
  - `401 Unauthorized`: Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - `404 Not Found`: Session not found.
    ```json
    {
      "error": "Session not found"
    }
    ```

---

### 5. Error Type Endpoints

#### **GET /error-types**
- **Description**: Retrieve the list of predefined error types.
- **Headers**: `Authorization: Bearer <token>` (optional, depending on access restrictions)
- **Responses**:
  - `200 OK`: List of error types.
    ```json
    [
      {
        "error_type_id": 1,
        "name": "Forehand"
      },
      {
        "error_type_id": 2,
        "name": "Backhand"
      },
      {
        "error_type_id": 3,
        "name": "Serve"
      },
      {
        "error_type_id": 4,
        "name": "Volley"
      }
    ]
    ```
  - `401 Unauthorized` (if restricted): Invalid or missing token.
    ```json
    {
      "error": "Unauthorized"
    }
    ```

---

## Summary of Functionality Covered
- **User Authentication**: Register (`POST /register`) and login (`POST /login`).
- **Session Management**: Start (`POST /sessions`), end (`PUT /sessions/{session_id}`), list (`GET /sessions`), and check active session (`GET /sessions/active`).
- **Error Logging**: Log an error (`POST /errors`) and undo the last error (`DELETE /errors/last`).
- **Summaries**: View error summary for a session (`GET /sessions/{session_id}/summary`).
- **Error Types**: Retrieve predefined error types (`GET /error-types`).

---

## Additional Notes
- **Error Timestamps**: Automatically generated server-side when logging errors via `POST /errors`.
- **Scalability**: The API supports multiple users, with data isolated by JWT authentication.
- **Frontend Support**: The `/sessions/active` endpoint simplifies checking for an ongoing session.

These API contracts provide a clear blueprint for backend and frontend teams to implement the Tennis Error Tracker web app, ensuring smooth integration and consistent data handling. Let me know if you need adjustments or additional endpoints!
