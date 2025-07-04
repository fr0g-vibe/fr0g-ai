#!/bin/bash

# ESMTP Webhook Test Script
# Tests the ESMTP processor with various email types and scenarios

BASE_URL="http://localhost:8080"
ESMTP_ENDPOINT="$BASE_URL/webhook/esmtp"

echo "=== ESMTP Webhook Processor Tests ==="
echo

# Test 1: Phishing email impersonating bank
echo "Test 1: Phishing email impersonating bank"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_001",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_001",
      "from": "security@fake-bank.com",
      "to": ["victim@example.com"],
      "cc": [],
      "bcc": [],
      "subject": "URGENT: Suspicious Activity Detected on Your Account",
      "body": "Dear Valued Customer,\n\nWe have detected suspicious activity on your account. Please click the link below to verify your identity immediately:\n\nhttp://secure-bank-verification.malicious-site.com/login\n\nFailure to verify within 24 hours will result in account suspension.\n\nBest regards,\nSecurity Team",
      "html_body": "<html><body><p>Dear Valued Customer,</p><p>We have detected suspicious activity on your account. Please <a href=\"http://secure-bank-verification.malicious-site.com/login\">click here</a> to verify your identity immediately.</p><p>Failure to verify within 24 hours will result in account suspension.</p><p>Best regards,<br>Security Team</p></body></html>",
      "headers": {
        "Return-Path": "bounce@fake-bank.com",
        "Received": "from mail.fake-bank.com",
        "Message-ID": "<20231201120000.12345@fake-bank.com>",
        "MIME-Version": "1.0",
        "Content-Type": "text/html; charset=UTF-8"
      },
      "attachments": [],
      "message_id": "20231201120000.12345@fake-bank.com",
      "in_reply_to": "",
      "references": [],
      "priority": "high",
      "metadata": {
        "spf_result": "fail",
        "dkim_result": "fail",
        "dmarc_result": "fail"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 2: Business Email Compromise (BEC) attack
echo "Test 2: Business Email Compromise (BEC) attack"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_002",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_002",
      "from": "ceo@company-fake.com",
      "to": ["finance@company.com"],
      "cc": [],
      "bcc": [],
      "subject": "URGENT: Wire Transfer Request",
      "body": "Hi,\n\nI need you to process an urgent wire transfer for a confidential acquisition. Please transfer $50,000 to the following account:\n\nBank: International Trust Bank\nAccount: 1234567890\nRouting: 987654321\nBeneficiary: Acquisition Holdings LLC\n\nThis is time sensitive and confidential. Please process immediately and confirm.\n\nThanks,\nJohn Smith\nCEO",
      "headers": {
        "Return-Path": "ceo@company-fake.com",
        "Received": "from mail.company-fake.com",
        "Message-ID": "<20231201130000.67890@company-fake.com>",
        "MIME-Version": "1.0",
        "Content-Type": "text/plain; charset=UTF-8"
      },
      "attachments": [],
      "message_id": "20231201130000.67890@company-fake.com",
      "priority": "urgent",
      "metadata": {
        "spf_result": "softfail",
        "dkim_result": "none",
        "dmarc_result": "fail"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 3: Malware attachment email
echo "Test 3: Malware attachment email"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_003",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_003",
      "from": "hr@legitimate-company.com",
      "to": ["employee@company.com"],
      "cc": [],
      "bcc": [],
      "subject": "Updated Employee Handbook - Please Review",
      "body": "Dear Employee,\n\nPlease find attached the updated employee handbook. All employees must review and acknowledge receipt by end of week.\n\nBest regards,\nHR Department",
      "headers": {
        "Return-Path": "hr@legitimate-company.com",
        "Received": "from mail.legitimate-company.com",
        "Message-ID": "<20231201140000.11111@legitimate-company.com>",
        "MIME-Version": "1.0",
        "Content-Type": "multipart/mixed"
      },
      "attachments": [
        {
          "filename": "Employee_Handbook_2023.pdf.exe",
          "content_type": "application/octet-stream",
          "size": 2048576,
          "hash": "d41d8cd98f00b204e9800998ecf8427e",
          "quarantined": false
        },
        {
          "filename": "Salary_Information.xlsx",
          "content_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
          "size": 1024000,
          "hash": "5d41402abc4b2a76b9719d911017c592",
          "quarantined": false
        }
      ],
      "message_id": "20231201140000.11111@legitimate-company.com",
      "metadata": {
        "spf_result": "pass",
        "dkim_result": "pass",
        "dmarc_result": "pass",
        "virus_scan": "suspicious"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 4: Trusted domain email (should have lower threat threshold)
echo "Test 4: Trusted domain email"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_004",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_004",
      "from": "notifications@github.com",
      "to": ["developer@company.com"],
      "cc": [],
      "bcc": [],
      "subject": "Security alert: New sign-in to your account",
      "body": "Hi developer,\n\nWe noticed a new sign-in to your GitHub account from a new device.\n\nDevice: Chrome on Windows\nLocation: San Francisco, CA\nTime: December 1, 2023 at 2:00 PM PST\n\nIf this was you, you can safely ignore this email. If not, please secure your account immediately.\n\nThanks,\nThe GitHub Team",
      "headers": {
        "Return-Path": "noreply@github.com",
        "Received": "from mail.github.com",
        "Message-ID": "<20231201150000.22222@github.com>",
        "MIME-Version": "1.0",
        "Content-Type": "text/plain; charset=UTF-8"
      },
      "attachments": [],
      "message_id": "20231201150000.22222@github.com",
      "metadata": {
        "spf_result": "pass",
        "dkim_result": "pass",
        "dmarc_result": "pass",
        "trusted_sender": "true"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 5: Cryptocurrency scam email
echo "Test 5: Cryptocurrency scam email"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_005",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_005",
      "from": "elon@tesla-crypto.com",
      "to": ["investor@example.com"],
      "cc": [],
      "bcc": [],
      "subject": "🚀 Exclusive Bitcoin Giveaway - Limited Time! 🚀",
      "body": "Greetings,\n\nI am excited to announce an exclusive Bitcoin giveaway! To participate:\n\n1. Send 0.1 to 2 BTC to: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\n2. Receive 2x to 20x back within 30 minutes!\n\nThis is a limited time offer to celebrate Tesla reaching $1 trillion market cap.\n\nDont miss this opportunity!\n\nElon Musk\nCEO, Tesla & SpaceX",
      "html_body": "<html><body><h1>🚀 Exclusive Bitcoin Giveaway! 🚀</h1><p>Send Bitcoin and receive 2x-20x back!</p><p><strong>Bitcoin Address:</strong> 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa</p></body></html>",
      "headers": {
        "Return-Path": "bounce@tesla-crypto.com",
        "Received": "from mail.tesla-crypto.com",
        "Message-ID": "<20231201160000.33333@tesla-crypto.com>",
        "MIME-Version": "1.0",
        "Content-Type": "text/html; charset=UTF-8"
      },
      "attachments": [],
      "message_id": "20231201160000.33333@tesla-crypto.com",
      "metadata": {
        "spf_result": "fail",
        "dkim_result": "fail",
        "dmarc_result": "fail",
        "reputation_score": "0.1"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 6: Ransomware email
echo "Test 6: Ransomware email"
curl -X POST "$ESMTP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "esmtp_test_006",
    "source": "mail_server",
    "tag": "esmtp",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "MailServer/2.1"
    },
    "body": {
      "id": "email_006",
      "from": "legal@law-firm-notice.com",
      "to": ["target@company.com"],
      "cc": [],
      "bcc": [],
      "subject": "Legal Notice - Immediate Action Required",
      "body": "LEGAL NOTICE\n\nYou have been served with a legal notice. Please open the attached document immediately to avoid legal consequences.\n\nTime sensitive - must respond within 24 hours.\n\nLegal Department",
      "headers": {
        "Return-Path": "legal@law-firm-notice.com",
        "Received": "from mail.law-firm-notice.com",
        "Message-ID": "<20231201170000.44444@law-firm-notice.com>",
        "MIME-Version": "1.0",
        "Content-Type": "multipart/mixed"
      },
      "attachments": [
        {
          "filename": "Legal_Notice_URGENT.pdf",
          "content_type": "application/pdf",
          "size": 512000,
          "hash": "098f6bcd4621d373cade4e832627b4f6",
          "quarantined": true
        }
      ],
      "message_id": "20231201170000.44444@law-firm-notice.com",
      "metadata": {
        "spf_result": "fail",
        "dkim_result": "none",
        "dmarc_result": "fail",
        "malware_detected": "true"
      }
    }
  }' | jq '.'
echo -e "\n"

echo "=== ESMTP Webhook Tests Complete ==="
