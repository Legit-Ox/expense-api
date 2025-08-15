#!/usr/bin/env python3
"""Test script to verify environment setup."""

import os
from pathlib import Path

# Load environment variables from .env file
def load_env():
    env_file = Path('.env')
    if env_file.exists():
        with open(env_file, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#') and '=' in line:
                    key, value = line.split('=', 1)
                    os.environ[key] = value

# Test the setup
load_env()
api_key = os.getenv('GROQ_API_KEY')

print("Environment Setup Test:")
print(f"✓ .env file exists: {Path('.env').exists()}")
print(f"✓ GROQ_API_KEY loaded: {bool(api_key)}")
if api_key:
    print(f"✓ API key starts with: {api_key[:10]}...")
print(f"✓ Model: {os.getenv('GROQ_MODEL', 'llama3-70b-8192')}")
