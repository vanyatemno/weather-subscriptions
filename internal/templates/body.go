package templates

const weatherEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weather Update</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            background-color: #3498db;
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            margin: -20px -20px 20px -20px;
        }
        .weather-info {
            background-color: #ecf0f1;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .weather-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 0;
            border-bottom: 1px solid #bdc3c7;
        }
        .weather-item:last-child {
            border-bottom: none;
        }
        .weather-label {
            font-weight: bold;
            color: #2c3e50;
        }
        .weather-value {
            color: #34495e;
            font-size: 18px;
        }
        .temperature {
            font-size: 24px;
            font-weight: bold;
            color: #e74c3c;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #7f8c8d;
            font-size: 14px;
        }
        @media only screen and (max-width: 600px) {
            .container {
                padding: 15px;
            }
            .weather-item {
                flex-direction: column;
                text-align: center;
            }
            .weather-label {
                margin-bottom: 5px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌤️ Weather Update</h1>
            <p>Current weather conditions for your location</p>
        </div>
        
        <div class="weather-info">
            <div class="weather-item">
                <span class="weather-label">🌡️ Temperature: </span>
                <span class="weather-value temperature">%s</span>
            </div>
            
            <div class="weather-item">
                <span class="weather-label">💧 Humidity: </span>
                <span class="weather-value">%s</span>
            </div>
            
            <div class="weather-item">
                <span class="weather-label">☁️ Conditions: </span>
                <span class="weather-value">%s</span>
            </div>
        </div>
        
        <div class="footer">
            <p>This is an automated weather notification.</p>
            <p>Stay safe and have a great day!</p>
            <p>Code to unsubscribe from mailer: %s</p>						
        </div>
    </div>
</body>
</html>`

const verificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            background-color: #2ecc71;
            color: white;
            padding: 20px;
            border-radius: 10px 10px 0 0;
            margin: -20px -20px 20px -20px;
        }
        .content {
            text-align: center;
            padding: 20px 0;
        }
        .verification-code {
            background-color: #f8f9fa;
            border: 2px solid #2ecc71;
            border-radius: 8px;
            padding: 30px 20px;
            margin: 30px 0;
            display: inline-block;
        }
        .code {
            font-size: 36px;
            font-weight: bold;
            color: #2c3e50;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
        }
        .instructions {
            background-color: #ecf0f1;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            text-align: left;
        }
        .instructions h3 {
            color: #2c3e50;
            margin-top: 0;
        }
        .instructions ul {
            color: #34495e;
            padding-left: 20px;
        }
        .warning {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            text-align: center;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #7f8c8d;
            font-size: 14px;
            border-top: 1px solid #ecf0f1;
            padding-top: 20px;
        }
        @media only screen and (max-width: 600px) {
            .container {
                padding: 15px;
            }
            .code {
                font-size: 28px;
                letter-spacing: 4px;
            }
            .verification-code {
                padding: 20px 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Email Verification</h1>
            <p>Please verify your email address</p>
        </div>
        
        <div class="content">
            <h2>Welcome!</h2>
            <p>To complete your registration, please enter the verification code below:</p>
            
            <div class="verification-code">
                <div class="code">%s</div>
            </div>
            
            <div class="instructions">
                <h3>📋 Instructions:</h3>
                <ul>
                    <li>Enter this 6-digit code in the verification field</li>
                    <li>The code is valid for 10 minutes</li>
                    <li>If you didn't request this code, please ignore this email</li>
                </ul>
            </div>
            
            <div class="warning">
                ⚠️ <strong>Security Notice:</strong> Never share this code with anyone. Our team will never ask for your verification code.
            </div>
        </div>
        
        <div class="footer">
            <p>This is an automated message. Please do not reply to this email.</p>
            <p>If you need assistance, contact our support team.</p>
        </div>
    </div>
</body>
</html>`
