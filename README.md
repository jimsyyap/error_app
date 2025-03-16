# Tennis Error Tracker

Tennis Error Tracker is a web application designed for tennis players to log unforced errors during matches and analyze their performance. By tracking errors by type (e.g., forehand, backhand, serve, volley), players can identify patterns and areas for improvement. The app supports multiple users, each with their own secure login and personal data.

This project is currently set up for local development, with plans to host it on Google Cloud in the future.

## Features

- **User Authentication**: Register and log in securely with JWT-based authentication.
- **Match Sessions**: Start and end match sessions to organize error logs.
- **Error Logging**: Quickly log unforced errors by type during a match with a simple, mobile-friendly interface.
- **Undo Functionality**: Easily undo the last error logged in the current session.
- **Error Summaries**: View detailed summaries of errors for each match session, including counts by error type and total errors.
- **Responsive Design**: Optimized for mobile devices, ensuring ease of use during matches.

## Tech Stack

- **Backend**: Go (using [Gin](https://github.com/gin-gonic/gin) or [Echo](https://echo.labstack.com/) framework)
- **Database**: PostgreSQL
- **Frontend**: ReactJS with [Axios](https://axios-http.com/) for API calls and [Material-UI](https://mui.com/) for styling

## Installation and Setup

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.16 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) (version 12 or higher)
- [Node.js](https://nodejs.org/en/download/) (version 14 or higher)
- [npm](https://www.npmjs.com/get-npm) (version 6 or higher)

### Backend Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/tennis-error-tracker.git
   cd tennis-error-tracker/backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:
   Create a `.env` file in the `backend` directory with the following content:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_pg_username
   DB_PASSWORD=your_pg_password
   DB_NAME=tennis_errors
   JWT_SECRET=your_jwt_secret_key
   ```
   Replace `your_pg_username`, `your_pg_password`, and `your_jwt_secret_key` with your actual PostgreSQL credentials and a secure JWT secret.

4. **Run the backend server**:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

### Database Setup

1. **Start PostgreSQL**:
   Ensure PostgreSQL is running locally on your machine.

2. **Create the database**:
   ```bash
   createdb tennis_errors
   ```

3. **Run migrations**:
   If you're using a migration tool like [golang-migrate](https://github.com/golang-migrate/migrate), run:
   ```bash
   migrate -path ./migrations -database "postgres://your_pg_username:your_pg_password@localhost:5432/tennis_errors?sslmode=disable" up
   ```
   Replace `your_pg_username` and `your_pg_password` with your PostgreSQL credentials.

4. **Seed error types**:
   Manually insert predefined error types (e.g., Forehand, Backhand, Serve, Volley) into the `error_types` table using a SQL client or a seed script.

### Frontend Setup

1. **Navigate to the frontend directory**:
   ```bash
   cd ../frontend
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Set up environment variables**:
   Create a `.env` file in the `frontend` directory with:
   ```env
   REACT_APP_BACKEND_URL=http://localhost:8080
   ```

4. **Start the development server**:
   ```bash
   npm start
   ```
   The frontend will be available at `http://localhost:3000`.

## Usage

1. **Register a new user**:
   - Open your browser and go to `http://localhost:3000/register`.
   - Fill out the form to create an account.

2. **Log in**:
   - Navigate to `http://localhost:3000/login` and enter your credentials.

3. **Start a match session**:
   - After logging in, click "Start New Session" to begin tracking errors for a match.

4. **Log errors**:
   - During the session, use the buttons labeled "Forehand", "Backhand", "Serve", or "Volley" to log unforced errors.
   - Click "Undo" to remove the most recently logged error.

5. **End the session**:
   - Click "End Session" to save and complete the current match session.

6. **View summaries**:
   - From the session list, select a past session to view a detailed summary of errors, including counts by type and totals.

## Future Plans

- **Hosting on Google Cloud**:
  - Deploy the backend to Google Cloud Run.
  - Use Cloud SQL for the PostgreSQL database.
  - Host the frontend as a static site on Cloud Storage.
- **Additional Features**:
  - Add support for custom error types.
  - Implement advanced analytics for deeper performance insights.

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new branch for your changes.
3. Submit a pull request with your updates.

For major changes, please open an issue first to discuss your ideas.

## Issues

If you encounter bugs or have feature requests, please report them on the [GitHub repository](https://github.com/yourusername/tennis-error-tracker/issues).

---

