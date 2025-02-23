# 🏦 Golang Banking Ledger

A simple banking ledger application built with **Go**, **PostgreSQL**, **MongoDB**, and **RabbitMQ**. This project demonstrates a microservices-style architecture with event-driven communication.

---

## 🌟 Features
- **Account Management**: Create, update, and delete accounts.
- **Transaction Handling**: Deposit and Withdraw money using transactions with detailed logs.
- **Event-Driven Architecture**: Uses RabbitMQ for asynchronous event processing.
- **Multi-Database Support**: PostgreSQL for accounts and MongoDB for transactions.


## 🚀 Getting Started

### **1. Clone the Repository**
Run the following commands:
`git clone https://github.com/your-username/golang-banking-ledger.git`  
`cd golang-banking-ledger`

### **2. Start the Dependencies with Docker**
Run:
`docker compose up -d`

### **3. Run the Go Application Locally**
Run:
`go run cmd/main.go`

### **4. Access the Application at localhost:8080**
### **5.NOTE ON CREDENTIAL**
 For local development, the .env file is included in the repository with example values. However, in production, you should never push the .env file to version control. Instead, use environment variables or a secure secrets management tool.
