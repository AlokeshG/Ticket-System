# Ticket System API

A RESTful backend API for managing support tickets.

This project provides user registration and login with JWT authentication, ticket creation, ticket listing, ticket details, and ticket status management with ownership-based authorization.

## Assumptions

- SQLite is used as the database for simplicity.
- Users can only access and modify tickets that belong to their account.
- New tickets start with `open` status.
- Valid status flow is `open -> in_progress -> closed`.
- Closed tickets cannot be reopened.
- JWT tokens are required for all ticket endpoints.
- `JWT_SECRET` is provided through an environment variable.
- The deployed free Render service uses an ephemeral filesystem, so SQLite data should be considered suitable for demonstration/testing rather than persistent production storage.

## 🚀 Live Deployment

**Base URL:**

https://ticket-system-dhmg.onrender.com

**Health Check:**

https://ticket-system-dhmg.onrender.com/health

Health response:

```json
{
  "status": "ok"
}
