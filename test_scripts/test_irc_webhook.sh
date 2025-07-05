#!/bin/bash

# IRC Webhook Test Script
# Tests the IRC processor with various message types and scenarios

BASE_URL="http://localhost:8080"
IRC_ENDPOINT="$BASE_URL/webhook/irc"

echo "=== IRC Webhook Processor Tests ==="
echo

# Test 1: Basic IRC channel message
echo "Test 1: Basic IRC channel message"
curl -X POST "$IRC_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "irc_test_001",
    "source": "irc_bot",
    "tag": "irc",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "IRC-Bot/1.0"
    },
    "body": {
      "id": "msg_001",
      "type": "PRIVMSG",
      "from": "suspicious_user",
      "to": "#general",
      "message": "Hey everyone! Check out this amazing deal: http://bit.ly/free-money-now - you can make $5000 in just one day!",
      "channel": "#general",
      "server": "irc.example.com",
      "is_private": false,
      "user_info": {
        "nickname": "suspicious_user",
        "username": "user123",
        "hostname": "192.168.1.100",
        "real_name": "John Doe",
        "channels": ["#general", "#random"],
        "is_op": false,
        "is_voice": false,
        "idle_time": 300
      },
      "metadata": {
        "client": "HexChat",
        "version": "2.16.1"
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 2: Private message with phishing attempt
echo "Test 2: Private message with phishing attempt"
curl -X POST "$IRC_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "irc_test_002",
    "source": "irc_bot",
    "tag": "irc",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "IRC-Bot/1.0"
    },
    "body": {
      "id": "msg_002",
      "type": "PRIVMSG",
      "from": "admin_fake",
      "to": "target_user",
      "message": "URGENT: Your account will be suspended. Please verify your credentials at http://fake-bank-login.com/verify immediately!",
      "channel": "",
      "server": "irc.example.com",
      "is_private": true,
      "user_info": {
        "nickname": "admin_fake",
        "username": "admin",
        "hostname": "suspicious.domain.com",
        "real_name": "System Administrator",
        "channels": [],
        "is_op": false,
        "is_voice": false,
        "idle_time": 0
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 3: Trusted user message (should have lower threat threshold)
echo "Test 3: Trusted user message"
curl -X POST "$IRC_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "irc_test_003",
    "source": "irc_bot",
    "tag": "irc",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "IRC-Bot/1.0"
    },
    "body": {
      "id": "msg_003",
      "type": "PRIVMSG",
      "from": "trusted_admin",
      "to": "#general",
      "message": "Server maintenance scheduled for tonight at 2 AM EST. Please save your work.",
      "channel": "#general",
      "server": "irc.example.com",
      "is_private": false,
      "user_info": {
        "nickname": "trusted_admin",
        "username": "admin",
        "hostname": "company.com",
        "real_name": "System Administrator",
        "channels": ["#general", "#admin"],
        "is_op": true,
        "is_voice": true,
        "idle_time": 60
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 4: Bot spam message
echo "Test 4: Bot spam message"
curl -X POST "$IRC_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "irc_test_004",
    "source": "irc_bot",
    "tag": "irc",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "IRC-Bot/1.0"
    },
    "body": {
      "id": "msg_004",
      "type": "PRIVMSG",
      "from": "crypto_bot_2023",
      "to": "#trading",
      "message": "STARTINGSTARTINGSTARTING PUMP ALERT! STARTINGSTARTINGSTARTING $SCAMCOIN is going TO THE MOON! 🌙 Buy now before its too late! Only 24 hours left! 💰💰💰",
      "channel": "#trading",
      "server": "irc.example.com",
      "is_private": false,
      "user_info": {
        "nickname": "crypto_bot_2023",
        "username": "bot",
        "hostname": "automated.spam.net",
        "real_name": "Crypto Trading Bot",
        "channels": ["#trading", "#crypto", "#finance"],
        "is_op": false,
        "is_voice": false,
        "idle_time": 0
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 5: Malware distribution attempt
echo "Test 5: Malware distribution attempt"
curl -X POST "$IRC_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "irc_test_005",
    "source": "irc_bot",
    "tag": "irc",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "IRC-Bot/1.0"
    },
    "body": {
      "id": "msg_005",
      "type": "PRIVMSG",
      "from": "file_sharer",
      "to": "#downloads",
      "message": "New software crack available! Download: http://malware-site.ru/crack.exe - Works 100% guaranteed! No virus I promise!",
      "channel": "#downloads",
      "server": "irc.example.com",
      "is_private": false,
      "user_info": {
        "nickname": "file_sharer",
        "username": "sharer",
        "hostname": "proxy.tor.exit",
        "real_name": "File Sharer",
        "channels": ["#downloads", "#warez"],
        "is_op": false,
        "is_voice": false,
        "idle_time": 1200
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

echo "=== IRC Webhook Tests Complete ==="
