#!/usr/bin/env python3
"""
PDF Transaction Extractor

This script extracts bank transactions from PDF statements.
It processes PDF files, extracts transaction data, and outputs JSON files.
"""

import argparse
import json
import math
import re
import time
from pathlib import Path
from typing import List, Dict, Any

import pdfplumber
import requests


class TransactionProcessor:
    """Handles PDF transaction extraction."""
    
    def __init__(self, api_key: str, model: str = "llama3-70b-8192"):
        self.api_key = api_key
        self.model = model
        self.api_url = "https://api.groq.com/openai/v1/chat/completions"
        self.rate_limit_delay = 8  # seconds between API calls
        
    def extract_text_from_pdf(self, pdf_path: str) -> List[str]:
        """Extract text from all pages of a PDF file."""
        pages_text = []
        try:
            with pdfplumber.open(pdf_path) as pdf:
                for page in pdf.pages:
                    text = page.extract_text()
                    if text:
                        pages_text.append(text)
        except Exception as e:
            print(f"Error reading PDF: {e}")
            return []
        
        return pages_text
    
    def extract_transactions_from_text(self, text_chunk: str) -> List[Dict[str, Any]]:
        """Extract transactions from a text chunk using AI."""
        prompt = f"""
        Extract all bank transactions from the following text.
        Output JSON array only.
        Each object must have: date (dd-mm-yyyy), details, amount (negative for debits, positive for credits).
        Ignore balance values. Skip headers.

        Text:
        {text_chunk}
        """
        
        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }
        
        data = {
            "model": self.model,
            "messages": [
                {"role": "system", "content": "You are a precise financial data extraction assistant."},
                {"role": "user", "content": prompt}
            ],
            "temperature": 0
        }
        
        try:
            response = requests.post(self.api_url, headers=headers, json=data)
            if response.status_code == 200:
                content = response.json()["choices"][0]["message"]["content"]
                return json.loads(content)
            else:
                print(f"API Error: {response.status_code} {response.text}")
                return []
        except Exception as e:
            print(f"Error extracting transactions: {e}")
            return []
    
    
    def process_pdf(self, pdf_path: str, output_path: str) -> bool:
        """Complete pipeline: PDF -> transactions -> JSON."""
        print(f"Processing PDF: {pdf_path}")
        
        # Step 1: Extract text from PDF
        pages_text = self.extract_text_from_pdf(pdf_path)
        if not pages_text:
            print("Failed to extract text from PDF")
            return False
        
        print(f"Extracted text from {len(pages_text)} pages")
        
        # Step 2: Extract transactions from each page
        all_transactions = []
        for i, page_text in enumerate(pages_text, 1):
            print(f"Processing page {i}/{len(pages_text)}")
            transactions = self.extract_transactions_from_text(page_text)
            all_transactions.extend(transactions)
            time.sleep(self.rate_limit_delay)  # Rate limiting
        
        if not all_transactions:
            print("No transactions found")
            return False
        
        print(f"Extracted {len(all_transactions)} transactions")
        
        # Step 3: Save to JSON file
        try:
            with open(output_path, "w", encoding="utf-8") as f:
                json.dump(all_transactions, f, indent=2, ensure_ascii=False)
            
            print(f"Saved {len(all_transactions)} transactions to {output_path}")
            return True
            
        except Exception as e:
            print(f"Error saving JSON file: {e}")
            return False


def main():
    """Main function with command line interface."""
    parser = argparse.ArgumentParser(
        description="Extract bank transactions from PDF statements"
    )
    parser.add_argument(
        "pdf_file",
        help="Path to the PDF bank statement file"
    )
    parser.add_argument(
        "-o", "--output",
        help="Output JSON file path (default: transactions.json)",
        default="transactions.json"
    )
    parser.add_argument(
        "--api-key",
        help="Groq API key (or set GROQ_API_KEY environment variable)",
        default=None
    )
    parser.add_argument(
        "--model",
        help="AI model to use (default: llama3-70b-8192)",
        default="llama3-70b-8192"
    )
    
    args = parser.parse_args()
    
    # Validate inputs
    pdf_path = Path(args.pdf_file)
    if not pdf_path.exists():
        print(f"PDF file not found: {pdf_path}")
        return 1
    
    # Get API key
    api_key = args.api_key
    if not api_key:
        import os
        api_key = os.getenv("GROQ_API_KEY")
        if not api_key:
            print("API key required. Use --api-key or set GROQ_API_KEY environment variable")
            return 1
    
    # Process the PDF
    processor = TransactionProcessor(api_key, args.model)
    success = processor.process_pdf(str(pdf_path), args.output)
    
    return 0 if success else 1


if __name__ == "__main__":
    exit(main())
