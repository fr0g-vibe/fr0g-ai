#!/bin/bash

# SMS Webhook Test Script
# Tests the SMS processor with various message types and scenarios

BASE_URL="http://localhost:8080"
SMS_ENDPOINT="$BASE_URL/webhook/sms"

echo "=== SMS Webhook Processor Tests ==="
echo

# Test 1: Basic SMS phishing attempt
echo "Test 1: SMS phishing attempt"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_001",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_001",
      "from": "+15551234567",
      "to": "+15559876543",
      "body": "URGENT: Your bank account has been compromised. Click here to secure it immediately: http://fake-bank.com/secure?token=abc123",
      "message_sid": "SM1234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "US",
      "region": "CA",
      "carrier": "Verizon",
      "message_type": "sms",
      "metadata": {
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 2: MMS with suspicious media
echo "Test 2: MMS with suspicious media"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_002",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_002",
      "from": "+15552345678",
      "to": "+15559876543",
      "body": "Check out these exclusive photos! 😉",
      "media_urls": [
        "https://suspicious-site.com/media/image1.jpg",
        "https://malware-host.net/download/app.apk"
      ],
      "message_sid": "MM1234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "US",
      "region": "NY",
      "carrier": "AT&T",
      "message_type": "mms",
      "metadata": {
        "num_media": "2",
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 3: Cryptocurrency scam
echo "Test 3: Cryptocurrency scam"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_003",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_003",
      "from": "+15553456789",
      "to": "+15559876543",
      "body": "🚀 BITCOIN GIVEAWAY! 🚀 Send 0.1 BTC to 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa and receive 2 BTC back! Limited time offer from Elon Musk! 💰",
      "message_sid": "SM2234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "US",
      "region": "TX",
      "carrier": "T-Mobile",
      "message_type": "sms",
      "metadata": {
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 4: Trusted number message (should have lower threat threshold)
echo "Test 4: Trusted number message"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_004",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_004",
      "from": "+15554567890",
      "to": "+15559876543",
      "body": "Hi! This is your bank. Your monthly statement is ready for review in your online banking portal.",
      "message_sid": "SM3234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "US",
      "region": "CA",
      "carrier": "Verizon",
      "message_type": "sms",
      "metadata": {
        "verified_sender": "true",
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 5: Romance scam
echo "Test 5: Romance scam"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_005",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_005",
      "from": "+15555678901",
      "to": "+15559876543",
      "body": "Hello beautiful! I am US Army soldier deployed in Afghanistan. I have $2.5 million that I need to transfer. Can you help me? I will share 50% with you for your assistance. Please reply with your bank details.",
      "message_sid": "SM4234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "NG",
      "region": "Lagos",
      "carrier": "MTN",
      "message_type": "sms",
      "metadata": {
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 6: 2FA bypass attempt
echo "Test 6: 2FA bypass attempt"
curl -X POST "$SMS_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sms_test_006",
    "source": "twilio",
    "tag": "sms",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "sms_006",
      "from": "+15556789012",
      "to": "+15559876543",
      "body": "Security Alert: Someone is trying to access your account. Reply with your 6-digit verification code to block this attempt. Do not share this code with anyone.",
      "message_sid": "SM5234567890abcdef",
      "status": "received",
      "direction": "inbound",
      "country": "US",
      "region": "FL",
      "carrier": "Sprint",
      "message_type": "sms",
      "metadata": {
        "price": "0.0075",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

echo "=== SMS Webhook Tests Complete ==="
