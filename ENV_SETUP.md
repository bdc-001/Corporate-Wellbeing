# Environment Variables Setup

This project requires environment variables to be configured. **Never commit `.env` files to git** - they contain sensitive information.

## Quick Setup

1. **Create a `.env` file** in the project root:
   ```bash
   cp .env.example .env  # If .env.example exists
   # OR create .env manually
   ```

2. **Set required environment variables** in your `.env` file:

   ```bash
   # Database Configuration (REQUIRED)
   POSTGRES_PASSWORD=your-secure-password-here
   DATABASE_URL=postgres://postgres:your-secure-password-here@localhost:5432/convin_crae?sslmode=disable

   # Security Secrets (REQUIRED - Generate strong random strings)
   JWT_SECRET=generate-a-strong-random-secret-here
   API_KEY_SECRET=generate-a-strong-random-secret-here

   # CORS Configuration
   CORS_ALLOW_ORIGINS=http://localhost:3000

   # Rate Limiting
   RATE_LIMIT_ENABLED=true
   RATE_LIMIT_RPS=100

   # Logging
   LOG_LEVEL=info
   LOG_FORMAT=json
   ```

## Generating Secure Secrets

Use one of these methods to generate secure random secrets:

**Using OpenSSL:**
```bash
openssl rand -base64 32  # For JWT_SECRET
openssl rand -base64 32  # For API_KEY_SECRET
openssl rand -base64 32  # For POSTGRES_PASSWORD
```

**Using Python:**
```python
import secrets
print(secrets.token_urlsafe(32))
```

**Using Node.js:**
```bash
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

## Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | PostgreSQL database password | `mySecurePassword123!` |
| `DATABASE_URL` | Full database connection string | `postgres://postgres:password@localhost:5432/convin_crae?sslmode=disable` |
| `JWT_SECRET` | Secret for JWT token signing | `generated-random-secret-32-chars` |
| `API_KEY_SECRET` | Secret for API key validation | `generated-random-secret-32-chars` |

## Optional Variables

- `CONVIN_API_KEY` - Convin API integration key
- `CONVIN_WEBHOOK_SECRET` - Webhook signature secret
- `TWILIO_ACCOUNT_SID` - Twilio account SID
- `TWILIO_AUTH_TOKEN` - Twilio auth token
- `SENTRY_DSN` - Sentry error tracking DSN
- `NEW_RELIC_LICENSE_KEY` - New Relic monitoring key

## Security Notes

⚠️ **IMPORTANT:**
- Never commit `.env` files to version control
- Use strong, unique passwords for production
- Rotate secrets regularly
- Use different secrets for development and production
- Store production secrets in a secure secret management system (AWS Secrets Manager, HashiCorp Vault, etc.)

## Docker Compose

When using `docker-compose`, environment variables are automatically loaded from your `.env` file. Make sure your `.env` file is in the project root directory.

```bash
docker-compose up -d
```

## Production Deployment

For production deployments:
1. Set all required environment variables in your deployment platform
2. Use environment-specific secret management
3. Never hardcode secrets in code or configuration files
4. Regularly audit and rotate secrets

