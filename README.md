# 📘 Payroll System - Golang (Gin, GORM, PostgreSQL)

This is a scalable backend payroll system designed to handle employee attendance, overtime, and reimbursement with JWT-based authentication.

---

## 📖 Table of Contents

1. [Background](#-background)
2. [Getting Started](#-getting-started)
3. [How-To Guides](#-how-to-guides)
4. [API Usage](#-api-usage)
5. [Software Architecture](#-software-architecture)
6. [Appendix](#-appendix)

---

## 📟 Background

### Context

In a company, employees are paid monthly based on an 8-hour workday (9AM-5PM), Monday to Friday. Payroll is prorated based on attendance. Overtime is paid at twice the prorated rate and reimbursements are added to the payslip.

### Objective

Build a backend system to manage payroll based on:

* Attendance
* Overtime
* Reimbursement

### Requirements

* 100 fake employees with salary, username, and password.
* 1 fake admin with credentials.
* Admin-defined attendance periods.
* Employee attendance, overtime, and reimbursement submission.
* Constraints:

  * Only one submission per day.
  * No weekend attendance.
  * Max 3 hours of overtime/day.
  * One payroll run per period.
* Payslip generation with breakdowns.
* Admin summary of all payslips.

---

## 🚀 Getting Started

### Requirements

* Golang 1.20+
* PostgreSQL
* GORM
* Gin

### Setup

```bash
# Clone the repo
$ git clone https://github.com/goraasep/payslip-generation-system.git
$ cd payroll-system

# Create and configure your .env file
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=payslip
ACCESS_SECRET=supersecretaccesskey123!@#
REFRESH_SECRET=refreshsecretkey456!@#

# Run the app
$ go run main.go
```

---

## 🛠️ How-To Guides

### 🔒 Default Credentials

* Admin: `admin@admin.com` with password `admin`
* User: `user1@user.com` with password `user`

### ✅ Register & Login

* POST `/register` to register new users
* POST `/login` to log in and receive access/refresh tokens

### 📋 View Current User

* GET `/api/me`

### 👥 Admin: Manage Users

* GET `/api/admin/users` to list all users

### 🗓️ Attendance Period

* Admin: POST `/api/admin/attendance-periods` to create new period
* Admin & User: GET `/api/attendance-periods` to list attendance periods

### 📅 User: Attendance Logs

* POST `/api/user/attendance-logs`
* GET `/api/attendance-logs`

### 🕒 User: Overtime Logs

* POST `/api/user/overtime-logs`
* GET `/api/overtime-logs`

### 💵 User: Reimbursement Logs

* POST `/api/user/reimburse-logs`
* GET `/api/reimburse-logs`

### 💰 Admin: Run Payroll

* POST `/api/admin/run-payroll`
* GET `/api/payrolls`

### 📈 User: Generate Payslip

* POST `/api/user/generate-payslip?pdf=true`

### 📊 Admin: Generate Payslip Summary

* POST `/api/admin/generate-payslip-summary?pdf=true`

---

## 🛡️ API Usage

### Auth Endpoints

```http
POST /register
POST /login
POST /refresh
GET  /api/me
```

### User

```http
GET  /api/profile/me
GET  /api/admin/users
```

### Attendance Period

```http
POST /api/admin/attendance-periods
GET  /api/attendance-periods?start=0&length=10&order=asc&field=id
```

### Attendance Log

```http
POST /api/user/attendance-logs
GET  /api/attendance-logs
```

### Overtime Log

```http
POST /api/user/overtime-logs
GET  /api/overtime-logs
```

### Reimbursement Log

```http
POST /api/user/reimburse-logs
GET  /api/reimburse-logs
```

### Payroll

```http
POST /api/admin/run-payroll
GET  /api/payrolls
```

### Payslip

```http
POST /api/user/generate-payslip?pdf=true
POST /api/admin/generate-payslip-summary?pdf=true
```

---

## 🧐 Software Architecture

### Tech Stack

* **Language**: Golang
* **Framework**: Gin
* **ORM**: GORM
* **Database**: PostgreSQL
* **Auth**: JWT

### Folder Structure

```
/config       # Configuration files
/controllers  # Request handlers / controllers
/dto          # Data transfer objects
/helpers      # Helper functions
/middleware   # Middleware (auth, logging, etc.)
/models       # GORM models
/routes       # Route definitions
/utils        # Utility functions
```

### Monitoring

* Statsviz monitoring tool: http://localhost:7000/debug/statsviz/


### Data Model Overview

| Model                    | Columns                                                                                                                                                                                                     |                |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| **AttendanceLog**        | id, attendance\_period\_id, user\_id, date, created\_at, updated\_at, deleted\_at                                                                                                                           |                |
| **AttendancePeriod**     | id, start\_date, end\_date, created\_at, updated\_at, deleted\_at                                                                                                                                           |                |
| **OvertimeLog**          | id, attendance\_period\_id, user\_id, date, hour, description, created\_at, updated\_at, deleted\_at                                                                                                        |                |
| **ReimburseLog**         | id, attendance\_period\_id, user\_id, date, amount, description, created\_at, updated\_at, deleted\_at                                                                                                      |                |
| **Payroll**              | id, attendance\_period\_id, processed\_at, created\_at, updated\_at, deleted\_at                                                                                                                            |                |
| **Payslip**              | id, payroll\_id, user\_id, base\_salary, prorated\_salary, overtime\_pay, overtime\_count, overtime\_hours, attendance\_count, reimbursement\_total, take\_home\_pay, created\_at, updated\_at, deleted\_at |                |
| **PayslipReimbursement** | id, payslip\_id, reimburse\_log\_id, description, amount, created\_at, updated\_at, deleted\_at                                                                                                             |                |
| **Role**                 | id, name, created\_at, updated\_at, deleted\_at                                                                                                                                                             |                |
| **User**                 | id, name, email, password, salary, created\_at, updated\_at, deleted\_at                                                                                                                                    |  |


### 📌 Appendix

### .env Sample

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=payslip
ACCESS_SECRET=supersecretaccesskey123!@#
REFRESH_SECRET=refreshsecretkey456!@#
```

### Notes

- Attendance on weekends is rejected
- Duplicate daily submissions are ignored
- Payroll is immutable once processed

