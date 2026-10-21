# 🎓 AegisVote: Faculty Presentation & Defense Guide
**IEEE SRS-Compliant Privacy-Preserving Electronic Voting System**

---

## 🚀 1. Quick Start Commands

### Backend (Go Unified Microservices Gateway)
```bash
cd backend
go run main.go
```
*Listens on `http://localhost:8080`. Automatically mounts all 5 microservice services (Auth, Token, Voting, Counting, Audit) with zero port conflicts.*

### Frontend (React 19 + TypeScript + Vite)
```bash
cd frontend
npm run dev
```
*Runs on `http://localhost:5173`. Proxies `/api` requests directly to backend on `:8080`.*

---

## 🎯 2. How the Presentation UI Works

| Feature | Where to find it | What to demonstrate to Faculty |
| :--- | :--- | :--- |
| **Persona Quick Switcher** | Top Header Bar | Click between **Voter 1 (Alice)**, **Voter 2 (Bob)**, **Officer**, **Authority**, and **Admin** in 1 click without typing passwords. |
| **Reset Demo State** | Top Header Bar (`⚡ Reset Demo`) | Resets all test users and seeds a fresh demo election ready for presentation. |
| **Privacy Wall Visualizer** | Top Banner | Demonstrates the cryptographic barrier between **Identity (Auth DB)** and **Ballot (Votes DB)**. |
| **Merkle Tree Inspector** | Officer / Visualizer | Real-time tree showing SHA-256 leaves and system Merkle Root. |
| **Live Tamper Demonstration** | Merkle Tree (`💥 Simulate DB Tamper`) | Simulates an insider attacking the database. Verification immediately turns **RED (TAMPER DETECTED)**. Restores to **GREEN** with 1 click. |
| **Jira QA Test Suite** | Top Header (`🧪 Jira QA Suite`) | Runs real automated test cases (`AEGIS-TC01` to `AEGIS-TC08`) live against the backend microservices with test execution trace and JSON/CSV export. |
| **Faculty Walkthrough** | Top Header (`🎓 Faculty Walkthrough`) | 8-step interactive storyboard with viva defense talking points for each phase. |

---

## 📋 3. Step-by-Step Viva Presentation Script (5-Minute Walkthrough)

### Step 1: System Overview & Architecture
1. Show the **System Architecture & Cryptographic Privacy Wall** at the top of the screen.
2. **Talking Point**: *"Our system implements architectural anonymity. Unlike traditional systems that store voter IDs alongside votes, our architecture completely separates authentication from vote casting using single-use blind tokens."*

### Step 2: Officer & Key Generation Ceremony
1. Click **"Officer"** persona in the top bar.
2. Show the active election: **"Faculty Review Demo: University Senate 2026"**.
3. **Talking Point**: *"When an election is initialized, the system generates an asymmetric RSA-2048 keypair. The private key is encrypted with AES-256 derived from the Counting Authority's secret passphrase. Plaintext ballots are NEVER stored."*

### Step 3: Voter 1 Casts Secret Ballot
1. Click **"Voter 1"** persona.
2. Click **"Cast Secret Ballot"**, select candidate **"Dr. Elena Vance"**, and click **"Confirm & Cast Ballot"**.
3. Point out the **Cryptographic Ballot Receipt**:
   - Merkle Leaf Hash: `SHA256(EncryptedBallot)`
   - System Merkle Root: Root hash of the append-only ledger.
4. **Talking Point**: *"The voter receives a cryptographic receipt that allows individual verifiability without compromising ballot secrecy."*

### Step 4: Sybil & Double-Voting Attack Defense
1. While still as **Voter 1**, show that the button now displays **"Already Voted (Locked)"**.
2. **Talking Point**: *"The system enforces strict single-use token invalidation and verifies the HasVoted flag. Any subsequent attempt to request a token returns HTTP 403 Forbidden."*

### Step 5: Voter 2 (Bob) Casts Vote & Tree Growth
1. Click **"Voter 2"** persona in the top bar.
2. Click **"Cast Secret Ballot"** and select **"Marcus Thorne"**.
3. Cast vote and show that the Merkle Root dynamically updates.

### Step 6: Live Database Tamper Simulation (The Showstopper!)
1. Click **"Officer"** persona in the top bar.
2. Under active elections, click **"Inspect Merkle Tree"**.
3. Show the green **"MERKLE INTEGRITY VERIFIED"** badge.
4. Click **"💥 Simulate DB Tamper"** to simulate an insider changing a vote in the database.
5. Click **"Re-Verify"** and show the badge immediately switch to:
   - **"🚨 TAMPER DETECTED / MISMATCH"** (Status: `tampered`, `verification_pass: false`).
6. Click **"🔄 Restore Ledger"** to restore authentic cryptographic integrity.
7. **Talking Point**: *"This proves mathematical tamper-evidence. If an attacker gains raw database access and alters even a single bit of any stored ballot, the recomputed Merkle root fails against the officer's digital signature."*

### Step 7: Counting Authority Decryption & Results
1. Click **"Authority"** persona in the top bar.
2. Choose the closed election, keep default passphrase `demo123`, and click **"🔓 Decrypt & Compute Final Tally"**.
3. Show the live vote tally with percentage progress bars.

### Step 8: Jira QA Automation Test Matrix
1. Click **"🧪 Jira QA Suite"** in the top bar.
2. Click **"▶ Run All Jira Tests"**.
3. Watch all 8 test cases execute in real time against the backend microservices (`AEGIS-TC01` to `AEGIS-TC08`).
4. Click on any test case row to view the expandable assertion execution trace.
5. Click **"📥 Export Xray JSON"** or **"📄 Export CSV"** to demonstrate integration with Jira, Xray, and Zephyr test management software.

---

## 🔒 4. SRS Requirements Traceability (IEEE v1.0)

| Requirement | Description | Jira Test Key | Verification Method |
| :--- | :--- | :--- | :--- |
| **REQ-AUTH-1** | Voter authentication via Voter ID & hash | `AEGIS-TC01` | Automated JWT auth validation |
| **REQ-TOK-1/2** | Cryptographically random blind token with Zero PII | `AEGIS-TC02` | Token schema inspection |
| **REQ-AUTH-4** | Rejection of double-voting (`HasVoted == true`) | `AEGIS-TC03` | Repeat request returns HTTP 403 |
| **REQ-TOK-4** | Single-use token invalidation | `AEGIS-TC03` | Token reuse returns HTTP 401 |
| **REQ-ENC-1/2** | RSA-2048 ballot encryption & Zero Plaintext DB | `AEGIS-TC04` | Database ciphertext assertion |
| **REQ-MRK-1/2** | Incremental SHA-256 Merkle Tree ledger | `AEGIS-TC05` | Tree root recalculation |
| **REQ-MRK-4** | Cryptographic tamper detection | `AEGIS-TC05` | Live tamper simulation |
| **REQ-DB-3** | Lifecycle state machine enforcement | `AEGIS-TC06` | Status transition assertion |
| **REQ-ENC-3** | Counting Authority decryption passphrase | `AEGIS-TC07` | 401 on wrong pass, 200 on correct |
| **REQ-DB-4** | Immutable audit trail & backup | `AEGIS-TC08` | Audit query & file backup export |
