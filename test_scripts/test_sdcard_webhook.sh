#!/bin/bash

# SD Card Webhook Test Script
# Tests the SD Card processor with various data scenarios

BASE_URL="http://localhost:8080"
SDCARD_ENDPOINT="$BASE_URL/webhook/sdcard"

echo "=== SD Card Webhook Processor Tests ==="
echo

# Test 1: Suspicious executable files
echo "Test 1: SD Card with suspicious executable files"
curl -X POST "$SDCARD_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sdcard_test_001",
    "source": "udev_monitor",
    "tag": "sdcard",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "SDCard-Monitor/1.0"
    },
    "body": {
      "id": "sdcard_001",
      "device_path": "/dev/sdb1",
      "mount_point": "/media/usb0",
      "file_system": "FAT32",
      "total_size": 8589934592,
      "used_size": 1073741824,
      "files": [
        {
          "path": "/media/usb0/autorun.inf",
          "name": "autorun.inf",
          "extension": ".inf",
          "size": 256,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "[autorun]\nopen=setup.exe\nicon=icon.ico\nlabel=My USB Drive",
          "mime_type": "text/plain"
        },
        {
          "path": "/media/usb0/setup.exe",
          "name": "setup.exe",
          "extension": ".exe",
          "size": 2048576,
          "is_executable": true,
          "is_hidden": false,
          "permissions": "-rwxr-xr-x",
          "hash": "d41d8cd98f00b204e9800998ecf8427e",
          "mime_type": "application/x-executable"
        },
        {
          "path": "/media/usb0/document.pdf.exe",
          "name": "document.pdf.exe",
          "extension": ".exe",
          "size": 1024000,
          "is_executable": true,
          "is_hidden": false,
          "permissions": "-rwxr-xr-x",
          "hash": "5d41402abc4b2a76b9719d911017c592",
          "mime_type": "application/x-executable"
        },
        {
          "path": "/media/usb0/.hidden_file",
          "name": ".hidden_file",
          "extension": "",
          "size": 512,
          "is_executable": false,
          "is_hidden": true,
          "permissions": "-rw-------",
          "content": "secret data here",
          "mime_type": "text/plain"
        }
      ]
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 2: Data exfiltration scenario
echo "Test 2: SD Card with potential data exfiltration"
curl -X POST "$SDCARD_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sdcard_test_002",
    "source": "udev_monitor",
    "tag": "sdcard",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "SDCard-Monitor/1.0"
    },
    "body": {
      "id": "sdcard_002",
      "device_path": "/dev/sdc1",
      "mount_point": "/media/usb1",
      "file_system": "NTFS",
      "total_size": 16106127360,
      "used_size": 5368709120,
      "files": [
        {
          "path": "/media/usb1/employee_database_backup.sql",
          "name": "employee_database_backup.sql",
          "extension": ".sql",
          "size": 52428800,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "-- Employee Database Backup\nCREATE TABLE employees (\n  id INT PRIMARY KEY,\n  name VARCHAR(100),\n  ssn VARCHAR(11),\n  salary DECIMAL(10,2)\n);",
          "mime_type": "text/plain"
        },
        {
          "path": "/media/usb1/passwords.txt",
          "name": "passwords.txt",
          "extension": ".txt",
          "size": 1024,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "admin:password123\nroot:admin\nuser1:qwerty\ndbadmin:database_pass_2023",
          "mime_type": "text/plain"
        },
        {
          "path": "/media/usb1/financial_data.xlsx",
          "name": "financial_data.xlsx",
          "extension": ".xlsx",
          "size": 10485760,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "mime_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        },
        {
          "path": "/media/usb1/customer_list.csv",
          "name": "customer_list.csv",
          "extension": ".csv",
          "size": 2097152,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "name,email,phone,ssn\nJohn Doe,john@email.com,555-1234,123-45-6789\nJane Smith,jane@email.com,555-5678,987-65-4321",
          "mime_type": "text/csv"
        }
      ]
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 3: Malware distribution scenario
echo "Test 3: SD Card with malware distribution"
curl -X POST "$SDCARD_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sdcard_test_003",
    "source": "udev_monitor",
    "tag": "sdcard",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "SDCard-Monitor/1.0"
    },
    "body": {
      "id": "sdcard_003",
      "device_path": "/dev/sdd1",
      "mount_point": "/media/usb2",
      "file_system": "FAT32",
      "total_size": 4294967296,
      "used_size": 1073741824,
      "files": [
        {
          "path": "/media/usb2/vacation_photos.zip",
          "name": "vacation_photos.zip",
          "extension": ".zip",
          "size": 104857600,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "hash": "098f6bcd4621d373cade4e832627b4f6",
          "mime_type": "application/zip"
        },
        {
          "path": "/media/usb2/codec_pack.exe",
          "name": "codec_pack.exe",
          "extension": ".exe",
          "size": 20971520,
          "is_executable": true,
          "is_hidden": false,
          "permissions": "-rwxr-xr-x",
          "hash": "e99a18c428cb38d5f260853678922e03",
          "mime_type": "application/x-executable"
        },
        {
          "path": "/media/usb2/keygen.exe",
          "name": "keygen.exe",
          "extension": ".exe",
          "size": 512000,
          "is_executable": true,
          "is_hidden": false,
          "permissions": "-rwxr-xr-x",
          "hash": "5d41402abc4b2a76b9719d911017c592",
          "mime_type": "application/x-executable"
        },
        {
          "path": "/media/usb2/readme.txt",
          "name": "readme.txt",
          "extension": ".txt",
          "size": 256,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "Run codec_pack.exe to install required codecs for viewing the photos. Use keygen.exe to generate license key.",
          "mime_type": "text/plain"
        }
      ]
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 4: Clean SD card (should have low threat level)
echo "Test 4: Clean SD Card with legitimate files"
curl -X POST "$SDCARD_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sdcard_test_004",
    "source": "udev_monitor",
    "tag": "sdcard",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "SDCard-Monitor/1.0"
    },
    "body": {
      "id": "sdcard_004",
      "device_path": "/dev/sde1",
      "mount_point": "/media/usb3",
      "file_system": "exFAT",
      "total_size": 32212254720,
      "used_size": 2147483648,
      "files": [
        {
          "path": "/media/usb3/IMG_001.jpg",
          "name": "IMG_001.jpg",
          "extension": ".jpg",
          "size": 2097152,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "hash": "a1b2c3d4e5f6789012345678901234567890abcd",
          "mime_type": "image/jpeg"
        },
        {
          "path": "/media/usb3/document.pdf",
          "name": "document.pdf",
          "extension": ".pdf",
          "size": 1048576,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "hash": "fedcba0987654321098765432109876543210fedcb",
          "mime_type": "application/pdf"
        },
        {
          "path": "/media/usb3/notes.txt",
          "name": "notes.txt",
          "extension": ".txt",
          "size": 512,
          "is_executable": false,
          "is_hidden": false,
          "permissions": "-rw-r--r--",
          "content": "Meeting notes from today:\n- Discussed project timeline\n- Reviewed budget\n- Next meeting scheduled for Friday",
          "mime_type": "text/plain"
        }
      ]
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

# Test 5: Live scan request (mount point only)
echo "Test 5: Live SD Card scan request"
curl -X POST "$SDCARD_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sdcard_test_005",
    "source": "udev_monitor",
    "tag": "sdcard",
    "timestamp": "'$(date -Iseconds)'",
    "headers": {
      "User-Agent": "SDCard-Monitor/1.0"
    },
    "body": {
      "mount_point": "/tmp"
    }
  }' | python3 -m json.tool 2>/dev/null || cat
echo -e "\n"

echo "=== SD Card Webhook Tests Complete ==="
