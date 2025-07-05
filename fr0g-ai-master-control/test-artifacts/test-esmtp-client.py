#!/usr/bin/env python3
"""
ESMTP Client Test Script for fr0g.ai Master Control Program
Tests the ESMTP threat vector interceptor functionality
"""

import smtplib
import time
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

def send_test_email(smtp_host="localhost", smtp_port=2525, 
                   from_addr="test@example.com", to_addr="admin@company.com",
                   subject="Test Email", body="Test message", 
                   test_name="Generic Test"):
    """Send a test email through the ESMTP interceptor"""
    
    print(f"\nTESTING {test_name}")
    print(f"   From: {from_addr}")
    print(f"   To: {to_addr}")
    print(f"   Subject: {subject}")
    
    try:
        # Create message
        msg = MIMEMultipart()
        msg['From'] = from_addr
        msg['To'] = to_addr
        msg['Subject'] = subject
        msg.attach(MIMEText(body, 'plain'))
        
        # Connect and send
        server = smtplib.SMTP(smtp_host, smtp_port)
        server.set_debuglevel(0)  # Set to 1 for verbose output
        
        # Send email
        text = msg.as_string()
        server.sendmail(from_addr, [to_addr], text)
        server.quit()
        
        print(f"   COMPLETED Email sent successfully")
        return True
        
    except Exception as e:
        print(f"   FAILED Failed to send email: {e}")
        return False

def main():
    print("SECURITY fr0g.ai ESMTP Threat Vector Interceptor Test Suite")
    print("=" * 60)
    
    # Test 1: Legitimate business email
    send_test_email(
        from_addr="ceo@company.com",
        to_addr="team@company.com", 
        subject="Q4 Strategy Meeting",
        body="Team, please join us for the Q4 strategy meeting on Friday at 2 PM.",
        test_name="Legitimate Business Email"
    )
    
    time.sleep(1)
    
    # Test 2: Newsletter email
    send_test_email(
        from_addr="newsletter@company.com",
        to_addr="subscriber@example.com",
        subject="Weekly Newsletter - AI Security Updates", 
        body="Welcome to our weekly newsletter featuring the latest in AI security research.",
        test_name="Newsletter Email"
    )
    
    time.sleep(1)
    
    # Test 3: Suspicious phishing attempt
    send_test_email(
        from_addr="security@bank-fake.com",
        to_addr="customer@example.com",
        subject="URGENT: Account Security Alert",
        body="Your account will be suspended in 24 hours unless you verify your credentials immediately. Click here: http://fake-bank-security.com/login",
        test_name="Suspicious Phishing Email"
    )
    
    time.sleep(1)
    
    # Test 4: Malware delivery attempt
    send_test_email(
        from_addr="admin@suspicious-domain.ru",
        to_addr="victim@company.com", 
        subject="Invoice Attached - Please Review",
        body="Please find attached invoice. Download and execute the file to view: http://malware-site.com/invoice.exe",
        test_name="Malware Delivery Attempt"
    )
    
    time.sleep(1)
    
    # Test 5: Social engineering attempt
    send_test_email(
        from_addr="it-support@company-fake.com",
        to_addr="employee@company.com",
        subject="IT Security Update Required",
        body="Your password expires today. Please update it immediately by clicking: http://fake-company-portal.com/update-password",
        test_name="Social Engineering Attempt"
    )
    
    print(f"\nCOMPLETED ESMTP Test Suite Completed!")
    print(f"\nTIP Tips:")
    print(f"   - Check MCP logs for threat analysis results")
    print(f"   - Monitor AI community reviews for each email")
    print(f"   - Watch for threat level classifications")
    print(f"   - Observe pattern recognition in email content")

if __name__ == "__main__":
    main()
