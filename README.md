# Re2no - Save Reddit Posts to Notion
<video src="https://github.com/user-attachments/assets/3e1ea527-657e-424f-afb0-dd97ae9ef868" height="600"></video>
---

## What is Re2no?

Re2no is a powerful tool that bridges the gap between Reddit and Notion. It allows users to browse Reddit, filter content, and save posts directly to their Notion databases with a single click. Whether you're a researcher, content creator, or just an avid Reddit user, Re2no helps you curate and organize information efficiently.

### Key Features

- **Smart Filtering**: Browse Reddit posts by subreddit, time range, and sort order.
- **One-Click Save**: Instantly save posts to your selected Notion database.
- **Seamless Integration**: Automatically maps Reddit post data (title, content, author, score, URL) to Notion properties.
- **Saved State Tracking**: Visual indicators for posts you've already saved.
- **Direct Access**: Open saved posts in Notion directly from the dashboard.
- **Secure Authentication**: Uses Notion's official OAuth flow for secure access.

---

## Tech Stack

- **Frontend**: Vue.js 3, TypeScript, Tailwind CSS, Vite
- **Backend**: Go (Golang), Gin Framework
- **Database**: PostgreSQL
- **Infrastructure**: Docker, Docker Compose
- **APIs**: Reddit API (Public), Notion API

---

## Getting Started

The easiest way to run Re2no locally is using Docker Compose.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- A Notion Integration (for API credentials)

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/Dipstick713/Re2no.git
    cd Re2no
    ```

2.  **Configure Environment Variables**
    Create a `.env` file in the root directory:
    ```bash
    cp .env.example .env
    ```
    Update the `.env` file with your credentials (see [Environment Variables](#environment-variables) below).

3.  **Run with Docker**
    ```bash
    docker-compose up --build
    ```

4.  **Access the App**
    - Frontend: `http://localhost:3000`
    - Backend: `http://localhost:8080`

For manual setup instructions, please refer to [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Environment Variables

You need to set up the following environment variables in your `.env` file:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `POSTGRES_USER` | Database user | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `postgres` |
| `POSTGRES_DB` | Database name | `re2no` |
| `JWT_SECRET` | Secret key for JWT tokens | **Required** |
| `NOTION_CLIENT_ID` | Notion Integration Client ID | **Required** |
| `NOTION_CLIENT_SECRET` | Notion Integration Client Secret | **Required** |
| `NOTION_REDIRECT_URI` | OAuth Redirect URI | `http://localhost:3000/dashboard` |
| `FRONTEND_URL` | URL of the frontend application | `http://localhost:3000` |

---

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and set up your development environment manually.