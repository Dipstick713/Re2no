# Contributing to Re2no

The following is a set of guidelines for contributing to Re2no. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it is governed by a Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Re2no. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

- **Use a clear and descriptive title** for the issue to identify the problem.
- **Describe the exact steps to reproduce the problem** in as much detail as possible.
- **Provide specific examples** to demonstrate the steps.
- **Describe the behavior you observed** after following the steps and point out what exactly is the problem with that behavior.
- **Explain which behavior you expected to see instead** and why.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Re2no, including completely new features and minor improvements to existing functionality.

- **Use a clear and descriptive title** for the issue to identify the suggestion.
- **Provide a step-by-step description of the suggested enhancement** in as much detail as possible.
- **Explain why this enhancement would be useful** to most Re2no users.

### Pull Requests

- Fill in the required template
- Do not include issue numbers in the PR title
- Include screenshots and animated GIFs in your pull request whenever possible.
- Follow the [styleguides](#styleguides)

## Development Guide

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) (Recommended)
- OR
- [Go](https://golang.org/doc/install) (v1.20+)
- [Node.js](https://nodejs.org/) (v18+)
- [PostgreSQL](https://www.postgresql.org/download/)

### Setting Up the Development Environment

#### Option 1: Using Docker (Recommended)

1.  **Clone the repository**
    ```bash
    git clone https://github.com/Dipstick713/Re2no.git
    cd Re2no
    ```

2.  **Configure Environment Variables**
    Create a `.env` file in the root directory based on the example below:
    ```env
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=re2no
    JWT_SECRET=your_jwt_secret
    NOTION_CLIENT_ID=your_notion_client_id
    NOTION_CLIENT_SECRET=your_notion_client_secret
    NOTION_REDIRECT_URI=http://localhost:3000/dashboard
    FRONTEND_URL=http://localhost:3000
    ```

3.  **Start the Application**
    ```bash
    docker-compose up --build
    ```
    The application will be available at `http://localhost:3000`.

#### Option 2: Manual Setup

1.  **Database Setup**
    - Install and start PostgreSQL.
    - Create a database named `re2no`.

2.  **Backend Setup**
    ```bash
    cd server
    cp .env.example .env # Ensure you configure your .env file
    go mod download
    go run main.go
    ```

3.  **Frontend Setup**
    ```bash
    cd client
    cp .env.example .env # Ensure you configure your .env file
    npm install
    npm run dev
    ```

## Project Structure

```
Re2no/
├── client/                 # Vue.js Frontend
│   ├── src/
│   │   ├── components/     # Reusable UI components
│   │   ├── views/          # Page views
│   │   ├── lib/            # API clients and helpers
│   │   └── ...
│   └── ...
├── server/                 # Go Backend
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # Auth and other middleware
│   ├── models/             # Database models
│   ├── notion/             # Notion API integration
│   ├── reddit/             # Reddit API integration
│   └── ...
└── ...
```

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### JavaScript/TypeScript Styleguide

- Use Prettier for formatting.
- Follow the ESLint configuration provided in the project.

### Go Styleguide

- Use `gofmt` to format your code.
- Follow standard Go conventions (Effective Go).