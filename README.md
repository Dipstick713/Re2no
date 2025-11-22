# Re2no - Save Reddit Posts to Notion

<div align="center">

**Save Reddit Smarter. Build Better in Notion.**

Discover, filter, and save Reddit posts directly to your Notion workspace.

[![Live Demo](https://img.shields.io/badge/demo-live-brightgreen)](https://re2no-site.vercel.app)
[![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-35495E?logo=vue.js&logoColor=4FC08D)](https://vuejs.org/)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://golang.org/)

</div>

---

## Demo Video

![demo](client/src/assets/re2no.gif)

---

## What is Re2no?

Re2no lets you save Reddit posts directly to your Notion workspace. Perfect for researchers, content creators, or anyone who wants to organize interesting Reddit content in one place.

### Key Features

- Browse and filter Reddit posts by subreddit, time, and sort order
- Save any post to Notion with one click
- Automatically organized in your Notion database
- Track saved posts and open them directly in Notion
- Delete posts from both the app and Notion
- Secure OAuth authentication with Notion

---

## Tech Stack

**Frontend:** Vue.js, TypeScript, Tailwind CSS  
**Backend:** Go, Gin, PostgreSQL  
**APIs:** Reddit API, Notion API  
**Hosting:** Vercel (Frontend), Render (Backend)

---

## How to Use

**Live App:** [https://re2no-site.vercel.app](https://re2no-site.vercel.app)

1. **Connect Notion** - Click "Connect Notion" and authorize access
2. **Select Database** - Choose which Notion database to save posts to
3. **Add Subreddits** - Enter subreddit names you want to browse
4. **Filter Posts** - Set time range and sort preferences, then fetch posts
5. **Save to Notion** - Click "Save to Notion" on any post you like
6. **Manage Posts** - View saved posts, open them in Notion, or delete them

---

## Contributing

Contributions are welcome! Here's how you can help:

### Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/Re2no.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Commit your changes: `git commit -m 'Add some feature'`
6. Push to your fork: `git push origin feature/your-feature-name`
7. Open a Pull Request

### Development Setup

**Backend (Go)**
```bash
cd server
go mod download
cp .env.example .env
# Update .env with your credentials
go run main.go
```

**Frontend (Vue.js)**
```bash
cd client
npm install
cp .env.example .env
# Set VITE_API_URL to your backend URL
npm run dev
```

### What to Contribute

- Bug fixes
- New features
- Documentation improvements
- Code quality improvements
- UI/UX enhancements

### Guidelines

- Write clear commit messages
- Follow existing code style
- Test your changes before submitting
- Update documentation if needed

---

## License

MIT License - see LICENSE file for details