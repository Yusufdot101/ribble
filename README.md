<div align="center">
# 〰️ Ripple

**Real-time chat. Minimal. Fast. Always live.**
</div>

---

## 🌊 What is Ripple?

Ripple is a real-time chat application where messages flow instantly — no refresh, no delay. Create groups, invite people, and just talk. Built as a deep dive into WebSockets, event-driven architecture, and full-stack development.

> Think Discord, but stripped down to what actually matters.

---

## ✨ Features

- ⚡ **Real-time messaging** — powered by WebSockets, zero polling
- 🏠 **Chats** — create and join chats
- 👥 **Online indicators** — see who's active in real time
- 🔐 **Auth** — JWT-based login and registration
- 📜 **Message history** — persisted and loaded on chat join
- 🔔 **Typing indicators** — "someone is typing..." done right
- 📱 **Responsive UI** — works on mobile and desktop

---

## 🛠 Tech Stack

| Layer | Tech |
|---|---|
| Frontend | React, TailwindCSS |
| Backend | Golang |
| Database | PostgreSQL |
| Auth | OAuth2 + JWT + bcrypt |

---

## 🚀 Getting Started

### Prerequisites

- Node.js v18+
- Docker
- Golang v1.26
- PostgreSQL running locally or via Docker
- `npm`

### Installation

```bash
# Clone the repo
git clone https://github.com/Yusufdot101/ripple.git
cd ripple
```

### Environment Variables

Create a `.env` in `/server`:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/ripple
JWT_SECRET=your_super_secret_key
PORT=4000
```

Create a `.env` in `/client`:

```env
VITE_API_URL=http://localhost:4000
VITE_SOCKET_URL=http://localhost:4000
```

### Run locally

```bash
# Terminal 1 — start the server
docker-compose up
```

App will be at `http://localhost:5173`

---

## 📁 Project Structure

```
ripple/
└── README.md
```

---


## 🤝 Contributing

PRs are welcome. For major changes, open an issue first so we can talk about it.

```bash
git checkout -b feature/your-feature
git commit -m "feat: add your feature"
git push origin feature/your-feature
```

---

## 📄 License

MIT — do whatever you want with it.

---

<div align="center">
  Built with ☕ and too many WebSocket debugging sessions.
</div>
