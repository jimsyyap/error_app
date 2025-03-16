# Tracker: Project Reference Guide

This document serves as a comprehensive reference for the **Tennis Error Tracker** web application, a tool designed for tennis players to log and analyze unforced errors during matches. It includes the project overview, technical stack, database schema, API contracts, file/folder structure, setup instructions, and a development roadmap—everything needed to understand, contribute to, or troubleshoot the project.

## 1. Project Overview
The Tennis Error Tracker helps tennis players track unforced errors (e.g., forehand, backhand, serve, volley) during matches to identify patterns and improve performance. It supports multiple users with secure logins and provides a mobile-friendly interface for real-time error logging.

### Key Features
- **User Authentication**: Secure registration and login using JWT.
- **Match Sessions**: Start/end sessions to group errors by match.
- **Error Logging**: Log errors by type with an undo option.
- **Error Summaries**: View detailed stats per session.
- **Responsive Design**: Optimized for mobile use during matches.

---

## 2. Technical Stack
- **Backend**: Go (using Gin or Echo framework)
- **Database**: PostgreSQL
- **Frontend**: ReactJS with Axios (API calls) and Material-UI (styling)

---

## 3. Database Schema
The database includes four tables:

### Users Table
- `user_id` (PK, UUID): Unique user identifier
- `username` (VARCHAR): User’s chosen username
- `password_hash` (VARCHAR): Hashed password

### Match_Sessions Table
- `session_id` (PK, UUID): Unique session identifier
- `user_id` (FK, UUID): References Users table
- `start_time` (TIMESTAMP): Session start (UTC)
- `end_time` (TIMESTAMP): Session end (UTC, nullable)

### Error_Types Table
- `error_type_id` (PK, INT): Unique error type identifier
- `name` (VARCHAR): Error type (e.g., "Forehand", "Backhand")

### Error_Logs Table
- `error_id` (PK, UUID): Unique error log identifier
- `session_id` (FK, UUID): References Match_Sessions
- `error_type_id` (FK, INT): References Error_Types
- `timestamp` (TIMESTAMP): Error log time (UTC)

---

## 4. API Contracts
The backend exposes RESTful endpoints. All except `/register` and `/login` require JWT authentication.

### Authentication
- **POST /register**: Create a new user
  - Request: `{ "username": "string", "password": "string" }`
  - Response: `201 Created` or `400 Bad Request` / `409 Conflict`
- **POST /login**: Get a JWT token
  - Request: `{ "username": "string", "password": "string" }`
  - Response: `{ "token": "jwt_token_string" }` or `401 Unauthorized`

### Match Sessions
- **POST /sessions**: Start a session
  - Response: `{ "session_id": "uuid" }`
- **PUT /sessions/{session_id}**: End a session
  - Response: `200 OK` or `404 Not Found`
- **GET /sessions**: List user’s sessions
  - Response: Array of `{ "session_id": "uuid", "start_time": "timestamp", "end_time": "timestamp" }`
- **GET /sessions/active**: Get active session
  - Response: `{ "session_id": "uuid", "start_time": "timestamp" }` or `204 No Content`

### Error Logging
- **POST /errors**: Log an error
  - Request: `{ "session_id": "uuid", "error_type_id": 1 }`
  - Response: `201 Created` or `400 Bad Request` / `404 Not Found`
- **DELETE /errors/last**: Undo last error
  - Request: `{ "session_id": "uuid" }`
  - Response: `200 OK` or `400 Bad Request` / `404 Not Found`

### Summaries
- **GET /sessions/{session_id}/summary**: Get session stats
  - Response: `{ "total_errors": 10, "errors_by_type": { "Forehand": 3, "Backhand": 2, ... } }`

### Error Types
- **GET /error-types**: List error types
  - Response: Array of `{ "error_type_id": 1, "name": "Forehand" }`

---

## 5. File/Folder Structure
```
error_app/
├── backend/                    # Go backend
│   ├── cmd/                    # Main entry point
│   ├── pkg/                    # Shared packages
│   ├── internal/               # Handlers, middleware
│   ├── migrations/             # DB migrations
│   ├── config/                 # Config files
│   ├── test/                   # Tests
│   └── go.mod                  # Go module
├── frontend/                   # React frontend
│   ├── public/                 # Static files
│   ├── src/                    # Components, pages
│   ├── config/                 # Config files
│   ├── test/                   # Tests
│   └── package.json            # Node.js config
├── docs/                       # Documentation
│   └── README.md               # Project README
└── .gitignore                  # Git ignore
```

---

## 6. Setup Instructions

### Backend
1. Go to `/backend`.
2. Run `go mod tidy` to install dependencies.
3. Set `.env` vars (e.g., `DB_HOST`, `JWT_SECRET`).
4. Start server: `go run cmd/main.go`.

### Database
1. Start PostgreSQL.
2. Create DB: `createdb tennis_errors`.
3. Run migrations to set up tables.
4. Seed `error_types` with predefined values.

### Frontend
1. Go to `/frontend`.
2. Run `npm install`.
3. Set `.env` vars (e.g., `REACT_APP_BACKEND_URL`).
4. Start dev server: `npm start`.

---

## 7. Roadmap
1. **Setup and Planning** (1-2 days): Git, tools, structure
2. **Database Design** (1-2 days): Schema and seeding
    **Tasks**:
    - Design the database schema:
    - `Users` table: `user_id` (PK), `username`, `password_hash`.
    - `Match_Sessions` table: `session_id` (PK), `user_id` (FK), `start_time`, `end_time`.
    - `Error_Types` table: `error_type_id` (PK), `name` (e.g., "Forehand", "Backhand").
    - `Error_Logs` table: `error_id` (PK), `session_id` (FK), `error_type_id` (FK), `timestamp`.
    - Write SQL scripts or use an ORM (e.g., GORM) to create the tables.
    - Seed the `Error_Types` table with predefined error types: "Forehand", "Backhand", "Serve", "Volley".
    **Deliverables**:
    - Fully implemented database schema in PostgreSQL.
    - Predefined error types populated in the database.
3. **Backend Development** (1-2 weeks): APIs, auth
    **Tasks**:
    - Set up a Go backend framework (e.g., Gin or Echo).
    - Implement authentication:
    - `POST /register`: Register a new user with username and password.
    - `POST /login`: Authenticate user and return a JWT token.
    - Middleware to protect endpoints with JWT validation.
    - Create API endpoints:
    - `POST /sessions`: Start a new match session.
    - `PUT /sessions/{session_id}`: End an active session.
    - `GET /sessions`: Retrieve a list of the user’s past sessions.
    - `POST /errors`: Log an error for the current session.
    - `DELETE /errors/last`: Undo the last error in the current session.
    - `GET /sessions/{session_id}/summary`: Get an error summary for a session.
    - `GET /error-types`: Retrieve the list of error types.
    - Implement database interactions for each endpoint using GORM or raw SQL.
    - Add input validation and error handling (e.g., invalid session IDs, unauthorized access).
    **Deliverables**:
    - A fully functional backend with secure authentication and API endpoints.
    - Backend connected to the PostgreSQL database.
4. **Frontend Development** (1-2 weeks): UI, API integration
    **Tasks**:
    - Set up the ReactJS project using `create-react-app` or a similar tool.
    - Create key components:
    - `Login.js`: Form for user login.
    - `Register.js`: Form for user registration.
    - `SessionStart.js`: Button to start a new session.
    - `ErrorLogging.js`: Interface with buttons for each error type and an undo option.
    - `SessionList.js`: Display a list of past sessions.
    - `SummaryView.js`: Show a summary of errors for a selected session.
    - Implement routing between pages (e.g., using React Router).
    - Use Axios or Fetch to connect to backend API endpoints.
    - Manage state (e.g., with React Context or Redux) for user sessions and error data.
    - Style the app for mobile-friendliness using CSS frameworks (e.g., Material-UI) or custom CSS.
    **Deliverables**:
    - A responsive ReactJS frontend with all core features.
    - Functional API integration for user interactions.
5. **Integration and Testing** (1 week): End-to-end tests
    **Tasks**:
    - Connect the frontend to the backend APIs.
    - Test key functionalities:
    - Authentication (register, login, logout).
    - Session management (start, end, list).
    - Error logging and undo actions.
    - Summary generation and display.
    - Perform end-to-end testing to ensure seamless operation.
    - Debug and resolve any issues (e.g., API errors, UI glitches).
    **Deliverables**:
    - A fully integrated and tested web application.
    - Bug-free functionality across all features.
6. **Documentation and Refinement** (2-3 days): Polish and docs
    **Tasks**:
    - Write or update `README.md` with:
    - Project overview.
    - Setup instructions.
    - Usage guide.
    - Refine the UI/UX based on testing feedback (e.g., adjust button sizes, add loading states).
    - Optimize performance if needed (e.g., minimize API calls, improve load times).
    - Document steps for future Google Cloud deployment.
    **Deliverables**:
    - Comprehensive project documentation.
    - A polished, user-friendly application ready for use.
**Future**: Deploy to Google Cloud (Cloud Run, Cloud SQL).
    **Tasks**:
    - Set up Google Cloud services:
    - Cloud SQL for PostgreSQL hosting.
    - Cloud Run for the Go backend.
    - Cloud Storage for the React frontend (static site).
    - Configure environment variables and secrets management.
    - Implement CI/CD with Cloud Build for automated deployments.
    - Test the deployed app in the cloud environment.
    **Deliverables**:
    - A live, hosted version of the Tennis Error Tracker accessible online.


### **Key Considerations**
- **Security**: Use hashed passwords (e.g., bcrypt) and secure JWT handling.
- **Performance**: Ensure the error logging interface is fast and responsive for real-time use during matches.
- **Scalability**: Design the database and backend to support multiple users efficiently.
- **Time Zones**: Store timestamps in UTC and handle conversions on the frontend.

---

### Recommendation
For future help requests, attach this **"Tennis Error Tracker: Project Reference Guide"**. It provides a complete picture of the project—purpose, tech details, structure, and plans—making it easier for others to assist you. If you need adjustments or additional sections, let me know!

---

Based on your query, you’re asking which of the previous responses I’ve provided would be most helpful to submit as an attachment if you were to ask for help in the future. After reviewing the options, I’ve determined that while several responses cover useful aspects of your project (e.g., project overview, API contracts, roadmap), none of them alone serve as a single, comprehensive reference that would provide all the essential details someone might need to assist you. Since no single existing response fully meets this need, I’ve created a new document below that consolidates the most critical information into one complete guide. This document, titled **"Tennis Error Tracker: Project Reference Guide"**, is what I recommend you submit as an attachment for future help requests.
