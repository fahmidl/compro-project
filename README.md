# Company Profile with CMS

A company profile website built with **React** (frontend), **Go + Gin** (backend), and **MongoDB** (database), featuring a CMS admin panel to manage site content.

## Features

- **Public Site**: Home (Hero), About Us, and Contact Us pages
- **CMS Admin Panel**: Edit text, images, and toggle visibility for all sections
- **Image Upload**: Upload and preview images directly from the admin panel
- **JWT Authentication**: Secure admin access with token-based auth

## Tech Stack

| Layer     | Technology                        |
|-----------|-----------------------------------|
| Frontend  | React 18, Vite, React Router, Axios |
| Backend   | Go, Gin Framework                 |
| Database  | MongoDB (official Go driver)      |
| Auth      | JWT tokens, bcrypt password hashing |

## Prerequisites

- [Go](https://go.dev/dl/) (1.21+)
- [Node.js](https://nodejs.org/) (18+)
- [MongoDB](https://www.mongodb.com/docs/manual/installation/) running locally or a MongoDB Atlas URI

## Getting Started

### 1. Clone the Repository

```bash
git clone git@github.com:fahmidl/compro-project.git
cd compro-project
```

### 2. Configure the Backend

```bash
cd backend
```

Edit the `.env` file with your settings:

```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=compro
JWT_SECRET=change-this-to-a-secure-secret-key
PORT=8080
```

### 3. Start MongoDB

Make sure MongoDB is running. If installed as a service:

```bash
sudo systemctl start mongod
```

Or run manually:

```bash
mongod --dbpath /path/to/data/db
```

### 4. Start the Backend

```bash
cd backend
go run main.go
```

The API server will run on `http://localhost:8080`.

### 5. Install Frontend Dependencies and Start

```bash
cd frontend
npm install
npm run dev
```

The frontend will run on `http://localhost:5173` (proxied to the backend).

### 6. Seed the Admin Account

Run this once to create the initial admin user:

```bash
curl -X POST http://localhost:8080/api/auth/seed \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

> The seed endpoint only works once. After an admin exists, it will return an error.

## Using the CMS

1. Go to `http://localhost:5173/admin/login`
2. Login with the credentials you seeded
3. From the dashboard, you can:
   - **Edit Hero Section**: title, subtitle, background image, visibility
   - **Edit About Section**: title, description, image, visibility
   - **Edit Contact Section**: title, address, phone, email, map embed URL, visibility

## API Endpoints

| Method | Endpoint           | Auth   | Description                   |
|--------|--------------------|--------|-------------------------------|
| GET    | `/api/content`     | No     | Fetch all public site content |
| PUT    | `/api/content`     | Yes    | Update site content (CMS)     |
| POST   | `/api/auth/login`  | No     | Admin login, returns JWT      |
| POST   | `/api/auth/seed`   | No     | Seed initial admin (one-time) |
| POST   | `/api/upload`      | Yes    | Upload image, returns URL     |

## Project Structure

```
compro-project/
├── backend/
│   ├── handlers/       # Route handlers (content, auth, upload)
│   ├── middleware/     # JWT auth middleware
│   ├── models/         # MongoDB models (SiteContent, Admin)
│   ├── routes/         # Route definitions
│   ├── uploads/        # Uploaded images (served as static)
│   ├── main.go         # Entry point
│   └── .env            # Environment config
├── frontend/
│   ├── src/
│   │   ├── components/ # Layout (Navbar, Footer)
│   │   ├── pages/      # Public pages
│   │   │   └── admin/  # CMS admin pages
│   │   ├── services/   # API service (Axios)
│   │   ├── App.jsx     # Routes
│   │   └── main.jsx    # Entry point
│   ├── vite.config.js  # Vite config with API proxy
│   └── index.html
└── .gitignore
```

## Build for Production

### Frontend

```bash
cd frontend
npm run build
```

The production build will be in `frontend/dist/`.

### Backend

```bash
cd backend
go build -o server main.go
./server
```

## License

This project is for demonstration purposes.
