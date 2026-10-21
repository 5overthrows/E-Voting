# 🗳️ AegisVote: Privacy-Preserving Electronic Voting System
> **Compliant with IEEE SRS v1.0 specifications for Anonymous Voting Tokens and Append-Only Merkle Tree Ledgers.**

AegisVote is a secure, privacy-preserving electronic voting prototype designed to demonstrate:
1. **Architectural Anonymity**: Complete decoupling of voter identification from vote casting via single-use blind tokens.
2. **Cryptographic Tamper-Evidence**: Append-only Merkle Tree ledger where any modification to stored ballots invalidates the cryptographic root.
3. **Faculty Defense & QA Testing Suite**: Interactive presentation UI with 1-click persona switching, live tamper simulation, and Jira/Xray automated test execution.

---

## 🛠️ Technology Stack (100% Intact)

- **Backend**: Go 1.22+ (Gin Web Framework, SQLite3 / PostgreSQL)
  - Cryptography: RSA-2048 OAEP, AES-256 GCM, SHA-256 Merkle Tree, PBKDF2 Key Derivation, JWT (HS256)
  - Microservice Architecture: Auth, Token, Voting, Counting, Audit, and Demo Services
- **Frontend**: React 19, TypeScript, Vite
  - Visual Architecture Diagram & Privacy Wall
  - Real-time Merkle Tree Visualizer & Interactive Tamper Engine
  - Jira / Xray / Zephyr Automated QA Test Execution Matrix

---

## 🚀 Running the Project

### 1. Start Backend Gateway
```bash
cd backend
go run main.go
```
*Listens on `http://localhost:8080` (or individual microservices via flags).*

### 2. Start Frontend UI
```bash
cd frontend
npm run dev
```
*Opens on `http://localhost:5173`.*

---

## 🧪 Testing & QA Software (Jira / Xray / Zephyr)

- Every interactive UI element is annotated with standardized `data-testid` attributes.
- Click **"🧪 Jira QA Suite"** in the top navigation bar to execute all 8 SRS test cases (`AEGIS-TC01` through `AEGIS-TC08`) live against the backend microservices.
- Export test results in standard Jira Xray JSON or CSV format.

Refer to [FACULTY_PRESENTATION_GUIDE.md](./FACULTY_PRESENTATION_GUIDE.md) for the complete 5-minute presentation script.
