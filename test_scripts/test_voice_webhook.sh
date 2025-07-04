#!/bin/bash

# Voice Webhook Test Script
# Tests the Voice processor with various call types and scenarios

BASE_URL="http://localhost:8080"
VOICE_ENDPOINT="$BASE_URL/webhook/voice"

echo "=== Voice Webhook Processor Tests ==="
echo

# Test 1: Tech support scam call
echo "Test 1: Tech support scam call"
curl -X POST "$VOICE_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "voice_test_001",
    "source": "twilio",
    "tag": "voice",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "call_001",
      "from": "+15551234567",
      "to": "+15559876543",
      "call_sid": "CA1234567890abcdef",
      "recording_url": "https://api.twilio.com/recordings/RE1234567890abcdef.mp3",
      "recording_duration": 180.5,
      "transcription": "Hello, this is Microsoft technical support. We have detected suspicious activity on your computer. Your Windows license has expired and your computer is infected with viruses. Please allow me to remote access your computer to fix these issues immediately.",
      "confidence": 0.92,
      "language": "en-US",
      "audio_format": "mp3",
      "file_size": 2890000,
      "direction": "inbound",
      "status": "completed",
      "country": "IN",
      "carrier": "Airtel",
      "voice_analysis": {
        "sentiment_score": -0.3,
        "emotion_scores": {
          "urgency": 0.8,
          "deception": 0.7,
          "confidence": 0.9
        },
        "stress_level": 0.4,
        "speech_rate": 150.0,
        "speaker_gender": "male",
        "estimated_age": 35,
        "accent_region": "South Asian",
        "background_noise": 0.2,
        "audio_quality": 0.8
      },
      "metadata": {
        "call_duration": "180",
        "price": "0.0085",
        "price_unit": "USD"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 2: IRS impersonation scam
echo "Test 2: IRS impersonation scam"
curl -X POST "$VOICE_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "voice_test_002",
    "source": "twilio",
    "tag": "voice",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "call_002",
      "from": "+15552345678",
      "to": "+15559876543",
      "call_sid": "CA2234567890abcdef",
      "recording_url": "https://api.twilio.com/recordings/RE2234567890abcdef.mp3",
      "recording_duration": 95.2,
      "transcription": "This is the Internal Revenue Service. You owe $3,247 in back taxes and penalties. If you do not pay immediately, a warrant will be issued for your arrest. Press 1 to speak with an agent or press 2 to make a payment now.",
      "confidence": 0.88,
      "language": "en-US",
      "audio_format": "wav",
      "file_size": 1520000,
      "direction": "inbound",
      "status": "completed",
      "country": "US",
      "carrier": "Unknown",
      "voice_analysis": {
        "sentiment_score": -0.6,
        "emotion_scores": {
          "threat": 0.9,
          "urgency": 0.95,
          "authority": 0.8
        },
        "stress_level": 0.3,
        "speech_rate": 120.0,
        "speaker_gender": "male",
        "estimated_age": 40,
        "accent_region": "American",
        "background_noise": 0.1,
        "audio_quality": 0.9
      },
      "metadata": {
        "call_duration": "95",
        "spoofed_caller_id": "true"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 3: Robocall with voice deepfake
echo "Test 3: Robocall with voice deepfake"
curl -X POST "$VOICE_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "voice_test_003",
    "source": "twilio",
    "tag": "voice",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "call_003",
      "from": "+15553456789",
      "to": "+15559876543",
      "call_sid": "CA3234567890abcdef",
      "recording_url": "https://api.twilio.com/recordings/RE3234567890abcdef.mp3",
      "recording_duration": 45.8,
      "transcription": "Hi, this is your grandson Tommy. I am in jail and need $5000 for bail money. Please do not tell mom and dad. Wire the money to this account immediately.",
      "confidence": 0.85,
      "language": "en-US",
      "audio_format": "mp3",
      "file_size": 732000,
      "direction": "inbound",
      "status": "completed",
      "country": "US",
      "carrier": "Verizon",
      "voice_analysis": {
        "sentiment_score": -0.4,
        "emotion_scores": {
          "desperation": 0.8,
          "urgency": 0.9,
          "deception": 0.85
        },
        "stress_level": 0.7,
        "speech_rate": 140.0,
        "voiceprint_id": "deepfake_detected",
        "speaker_gender": "male",
        "estimated_age": 22,
        "accent_region": "American",
        "background_noise": 0.05,
        "audio_quality": 0.95
      },
      "metadata": {
        "call_duration": "45",
        "ai_generated": "suspected"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 4: Trusted number call (should have lower threat threshold)
echo "Test 4: Trusted number call"
curl -X POST "$VOICE_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "voice_test_004",
    "source": "twilio",
    "tag": "voice",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "call_004",
      "from": "+15554567890",
      "to": "+15559876543",
      "call_sid": "CA4234567890abcdef",
      "recording_url": "https://api.twilio.com/recordings/RE4234567890abcdef.mp3",
      "recording_duration": 120.0,
      "transcription": "Hello, this is Sarah from First National Bank calling to inform you about new security features we have added to your account. Please visit our website or call our official number to learn more.",
      "confidence": 0.95,
      "language": "en-US",
      "audio_format": "wav",
      "file_size": 1920000,
      "direction": "inbound",
      "status": "completed",
      "country": "US",
      "carrier": "AT&T",
      "voice_analysis": {
        "sentiment_score": 0.2,
        "emotion_scores": {
          "professional": 0.9,
          "helpful": 0.8,
          "trustworthy": 0.85
        },
        "stress_level": 0.1,
        "speech_rate": 130.0,
        "speaker_gender": "female",
        "estimated_age": 28,
        "accent_region": "American",
        "background_noise": 0.05,
        "audio_quality": 0.95
      },
      "metadata": {
        "call_duration": "120",
        "verified_caller": "true"
      }
    }
  }' | jq '.'
echo -e "\n"

# Test 5: Cryptocurrency investment scam
echo "Test 5: Cryptocurrency investment scam"
curl -X POST "$VOICE_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "voice_test_005",
    "source": "twilio",
    "tag": "voice",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "TwilioProxy/1.1"
    },
    "body": {
      "id": "call_005",
      "from": "+15555678901",
      "to": "+15559876543",
      "call_sid": "CA5234567890abcdef",
      "recording_url": "https://api.twilio.com/recordings/RE5234567890abcdef.mp3",
      "recording_duration": 240.3,
      "transcription": "Congratulations! You have been selected for an exclusive cryptocurrency investment opportunity. Bitcoin is going to reach $500,000 by next month. Invest just $1000 today and you will make $50,000 guaranteed. This offer expires in 24 hours. Do not miss this once in a lifetime opportunity.",
      "confidence": 0.90,
      "language": "en-US",
      "audio_format": "mp3",
      "file_size": 3845000,
      "direction": "inbound",
      "status": "completed",
      "country": "US",
      "carrier": "T-Mobile",
      "voice_analysis": {
        "sentiment_score": 0.6,
        "emotion_scores": {
          "excitement": 0.9,
          "urgency": 0.85,
          "greed": 0.8,
          "deception": 0.75
        },
        "stress_level": 0.2,
        "speech_rate": 160.0,
        "speaker_gender": "male",
        "estimated_age": 30,
        "accent_region": "American",
        "background_noise": 0.15,
        "audio_quality": 0.85
      },
      "metadata": {
        "call_duration": "240",
        "high_pressure_sales": "detected"
      }
    }
  }' | jq '.'
echo -e "\n"

echo "=== Voice Webhook Tests Complete ==="
