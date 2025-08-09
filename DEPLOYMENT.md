# 🚀 Deployment Guide for Expense API

This guide covers multiple deployment options for your Go Expense Tracking API.

## 📋 Prerequisites

- GitHub repository (already set up)
- GitHub Actions enabled
- Platform-specific accounts (Render, Fly.io, DigitalOcean, etc.)

## 🎯 Quick Start: Deploy to Render (Recommended)

### 1. Create Render Account
- Go to [render.com](https://render.com)
- Sign up with your GitHub account

### 2. Connect Repository
- Click "New +" → "Web Service"
- Connect your GitHub repository: `Legit-Ox/expense-api`
- Render will auto-detect your `render.yaml` configuration

### 3. Configure Environment Variables
- `DB_URL`: Will be automatically set by Render's PostgreSQL
- `PORT`: 8080 (already set)
- `ENV`: production (already set)

### 4. Deploy
- Click "Create Web Service"
- Your API will be live in minutes!

## 🚂 Deploy to Railway

### 1. Create Railway Account
- Go to [railway.app](https://railway.app)
- Sign up with GitHub

### 2. Deploy
- Click "Start a New Project"
- Choose "Deploy from GitHub repo"
- Select your repository
- Railway will use your `railway.json` configuration

## 🐳 Deploy to Fly.io

### 1. Install Flyctl
```bash
# macOS
brew install flyctl

# Or download from https://fly.io/docs/hands-on/install-flyctl/
```

### 2. Login to Fly.io
```bash
flyctl auth login
```

### 3. Deploy
```bash
flyctl launch
flyctl deploy
```

### 4. Set up GitHub Actions
Add these secrets to your GitHub repository:
- `FLY_API_TOKEN`: Get from `flyctl auth token`

## ☁️ Deploy to DigitalOcean App Platform

### 1. Create DigitalOcean Account
- Go to [digitalocean.com](https://digitalocean.com)
- Sign up and add payment method

### 2. Create App
- Go to "Apps" → "Create App"
- Connect your GitHub repository
- DigitalOcean will use your `.do/app.yaml` configuration

### 3. Set up GitHub Actions
Add these secrets to your GitHub repository:
- `DIGITALOCEAN_ACCESS_TOKEN`: Create in API section

## 🐙 Deploy with Docker Hub

### 1. Create Docker Hub Account
- Go to [hub.docker.com](https://hub.docker.com)
- Sign up and create a repository

### 2. Set up GitHub Actions
Add these secrets to your GitHub repository:
- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token

### 3. Deploy
- Push to main branch or create a tag
- GitHub Actions will build and push Docker image
- Use the image in any platform that supports Docker

## 🔐 Setting up GitHub Secrets

### For Render:
1. Go to your GitHub repository
2. Settings → Secrets and variables → Actions
3. Add these secrets:
   - `RENDER_TOKEN`: Get from Render dashboard
   - `RENDER_SERVICE_ID`: Your Render service ID

### For Fly.io:
1. Add `FLY_API_TOKEN` secret
2. Get token with: `flyctl auth token`

### For DigitalOcean:
1. Add `DIGITALOCEAN_ACCESS_TOKEN` secret
2. Create token in DigitalOcean API section

### For Docker Hub:
1. Add `DOCKERHUB_USERNAME` secret
2. Add `DOCKERHUB_TOKEN` secret
3. Create access token in Docker Hub account settings

## 🌐 Custom Domain Setup

### Render:
- Go to your service → Settings → Custom Domains
- Add your domain and follow DNS instructions

### Fly.io:
```bash
flyctl certs add yourdomain.com
flyctl domains add yourdomain.com
```

### DigitalOcean:
- Go to your app → Settings → Domains
- Add custom domain and configure DNS

## 📊 Monitoring & Health Checks

All deployment configurations include:
- ✅ Health check endpoint: `/health`
- ✅ Automatic restarts on failure
- ✅ Logging and monitoring
- ✅ SSL certificates (automatic)

## 🔄 Auto-Deployment

Your GitHub Actions workflows will:
- ✅ Run tests on every push
- ✅ Build the application
- ✅ Deploy automatically to your chosen platform
- ✅ Trigger on main branch pushes and pull requests

## 🆘 Troubleshooting

### Common Issues:
1. **Database Connection**: Ensure `DB_URL` is set correctly
2. **Port Binding**: Make sure `PORT` environment variable is set
3. **Build Failures**: Check Go version compatibility
4. **Deployment Timeouts**: Increase timeout values in platform settings

### Getting Help:
- Check GitHub Actions logs for detailed error messages
- Review platform-specific documentation
- Check the `/health` endpoint for API status

## 🎉 Success!

Once deployed, your API will be available at:
- **Render**: `https://your-app-name.onrender.com`
- **Fly.io**: `https://your-app-name.fly.dev`
- **DigitalOcean**: `https://your-app-name.ondigitalocean.app`
- **Railway**: `https://your-app-name.railway.app`

Test your deployment:
```bash
curl https://your-app-url/health
curl https://your-app-url/api/categories
```

## 📚 Next Steps

1. **Set up monitoring**: Add logging and metrics
2. **Configure backups**: Set up database backups
3. **Add CI/CD**: Enhance GitHub Actions workflows
4. **Scale up**: Upgrade to paid plans as needed
5. **Custom domain**: Add your own domain name 