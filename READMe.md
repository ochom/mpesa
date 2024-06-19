# Mpesa Broker

## Overview

## Running the APP

## Configuration

### Env Vars

```sh
# Basic credentials
DATABASE_DRIVER=0 # 0=>sqlite, 1=>postgres, 2=>mysql
DATABASE_URL=""
BASE_URL=""
MPESA_API_URL=""

# MpesaC2B credentials
MPESA_C2B_SHORT_CODE=""
MPESA_C2B_PASSKEY=""
MPESA_C2B_CONSUMER_KEY=""
MPESA_C2B_CONSUMER_SECRET=""
CLIENT_DEPOSIT_VALIDATION_URL=""
CLIENT_DEPOSIT_CONFIRMATION_URL=""

# MpesaB2C credentials
B2C_ALLOWED_ORIGINS=""
MPESA_B2C_SHORT_CODE=""
MPESA_B2C_PASSKEY=""
MPESA_B2C_CONSUMER_KEY=""
MPESA_B2C_CONSUMER_SECRET=""
MPESA_B2C_INITIATOR_NAME=""
MPESA_B2C_INITIATOR_PASSWORD=""
MPESA_B2C_CERTIFICATE_PATH=""
MPESA_B2C_PAYMENT_COMMENT=""

# Tax Remittance credentials
TAX_REMITTANCE_ALLOWED_ORIGINS=""
MPESA_TAX_CONSUMER_SECRET=""
MPESA_TAX_CONSUMER_KEY=""
MPESA_TAX_INITIATOR_PASSWORD=""
MPESA_TAX_CERTIFICATE_PATH=""
```