# CI/CD Pipeline Setup Guide

This guide explains how to configure the CI/CD pipeline for automated testing and Docker image building.

## Overview

The CI/CD pipeline consists of three jobs:
1. **test-backend**: Runs Go tests with PostgreSQL
2. **test-frontend**: Runs npm tests and builds the frontend
3. **build**: Builds and pushes Docker images to Docker Hub (optional)

## Required Setup

### 1. GitHub Secrets Configuration

For the Docker build and push step to work, you need to configure GitHub secrets:

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secrets:

#### Required for Docker Hub Push:
- **`DOCKER_USERNAME`**: Your Docker Hub username
- **`DOCKER_PASSWORD`**: Your Docker Hub password or Personal Access Token (recommended)

#### Optional Secrets:
- None required for basic CI/CD (tests will run without Docker secrets)

### 2. Docker Hub Setup

#### Option A: Using Password (Less Secure)
1. Use your Docker Hub account password
2. ⚠️ **Not recommended** - Use Personal Access Token instead

#### Option B: Using Personal Access Token (Recommended)
1. Log in to [Docker Hub](https://hub.docker.com/)
2. Go to **Account Settings** → **Security** → **New Access Token**
3. Create a token with **Read & Write** permissions
4. Copy the token and use it as `DOCKER_PASSWORD` secret

### 3. Pipeline Behavior

#### Without Docker Hub Secrets:
- ✅ Tests will run successfully
- ✅ Docker images will be built locally (not pushed)
- ✅ Pipeline will complete successfully
- ℹ️ No Docker images will be pushed to registry

#### With Docker Hub Secrets:
- ✅ Tests will run successfully
- ✅ Docker images will be built
- ✅ Docker images will be pushed to Docker Hub
- ✅ Images will be available at:
  - `$DOCKER_USERNAME/convin-backend:latest`
  - `$DOCKER_USERNAME/convin-frontend:latest`

## Pipeline Jobs

### test-backend
- **Trigger**: On push to `main` or `develop`, or on pull requests
- **Requirements**: None
- **Actions**:
  - Sets up PostgreSQL test database
  - Installs Go dependencies
  - Runs Go tests with race detection
  - Uploads test coverage

### test-frontend
- **Trigger**: On push to `main` or `develop`, or on pull requests
- **Requirements**: None
- **Actions**:
  - Sets up Node.js 18
  - Installs npm dependencies (`npm ci`)
  - Runs linter (warnings don't fail)
  - Runs tests (warnings don't fail)
  - Builds frontend

### build
- **Trigger**: Only on push to `main` branch
- **Requirements**: 
  - `test-backend` must pass
  - `test-frontend` must pass
- **Actions**:
  - Sets up Docker Buildx
  - Logs into Docker Hub (if secrets are configured)
  - Builds backend Docker image
  - Builds frontend Docker image
  - Pushes images to Docker Hub (if secrets are configured)

## Troubleshooting

### Docker Hub Login Fails

**Error**: `Error: Process completed with exit code 1` at "Login to Docker Hub"

**Solutions**:
1. **Check secrets are set**:
   - Go to repository Settings → Secrets → Actions
   - Verify `DOCKER_USERNAME` and `DOCKER_PASSWORD` exist

2. **Verify credentials**:
   - Test login manually: `docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD`
   - Use Personal Access Token instead of password

3. **Check token permissions**:
   - Personal Access Token must have **Read & Write** permissions

4. **Skip Docker push** (if not needed):
   - The pipeline will still build images locally
   - Just don't set the Docker secrets
   - Pipeline will complete successfully

### Build Fails

**Error**: Docker build fails

**Solutions**:
1. Check Dockerfile syntax
2. Verify all dependencies are available
3. Check build logs for specific errors
4. Test Docker build locally: `docker build -t test ./backend`

### Tests Fail

**Error**: Backend or frontend tests fail

**Solutions**:
1. Run tests locally to reproduce
2. Check test logs for specific failures
3. Ensure database migrations are up to date
4. Verify environment variables are set correctly

## Manual Testing

### Test Backend Locally
```bash
cd backend
go test -v ./...
```

### Test Frontend Locally
```bash
cd frontend
npm ci
npm test
npm run build
```

### Build Docker Images Locally
```bash
# Backend
docker build -t convin-backend:latest ./backend

# Frontend
docker build -t convin-frontend:latest ./frontend
```

## Best Practices

1. **Use Personal Access Tokens** instead of passwords for Docker Hub
2. **Rotate secrets regularly** for security
3. **Test locally** before pushing to avoid CI failures
4. **Monitor pipeline** for failures and fix promptly
5. **Use branch protection** to require CI passes before merge

## Security Notes

⚠️ **Important**:
- Never commit secrets to the repository
- Use GitHub Secrets for all sensitive data
- Rotate Docker Hub tokens regularly
- Use least-privilege tokens (only necessary permissions)
- Review pipeline logs for exposed secrets (though GitHub masks them)

## Next Steps

1. Set up GitHub Secrets (if you want Docker push)
2. Push to `main` branch to trigger pipeline
3. Monitor pipeline execution in **Actions** tab
4. Verify Docker images are pushed (if secrets configured)

