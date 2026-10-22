# SOFTWARE ENGINEERING LABORATORY RECORD / EXPERIMENT DOSSIER
**Department of Computer Science and Engineering**  
**School of Computer Science and Engineering**  
**Course:** Software Engineering Laboratory (SWE Lab)  
**Academic Year:** 2025 – 2026  
**Student Name:** VISHAL V  
**Register Number:** 24BPS1089  

---

## Laboratory Experiments Index

| Exp No. | Experiment Title | Category / Phase | Modelling & Testing Tools |
| :--- | :--- | :--- | :--- |
| **Exp 2** | Project Scheduling & Work Breakdown Structures (WBS) | Project Planning & Scheduling | GanttProject, MS Project, StarUML |
| **Exp 3** | Requirement Modelling: Context (L-0), L-1 & L-2 DFDs | Structured Analysis & DFD | StarUML / DFD Engine |
| **Exp 4** | Requirement Modelling: Entity-Relationship (ER) Diagram | Database & Schema Design | Relational Modeller, StarUML |
| **Exp 5** | Behavioral Modelling: State Transition Diagram | Dynamic Behavioral Modelling | StarUML State Machine Modeller |
| **Exp 6** | Object-Oriented Design: Use Case and Class Diagrams | OOAD & Structural Modelling | StarUML OOAD Engine |
| **Exp 7** | Interaction Modelling: Sequence & Communication Diagrams | Dynamic Interaction Modelling | StarUML Interaction Engine |
| **Exp 8** | OO Design: Package, Component & Deployment Models | System Architecture & Deployment | StarUML Architectural Modeller |
| **Exp 9** | Software Testing: Design and Demonstration of Test Cases | Quality Assurance & Testing | Selenium IDE, Apache JMeter |

---

## EXPERIMENT 2: Project Scheduling & Work Breakdown Structures (WBS)

### Aim
To perform software project scheduling, decomposition, duration estimation, task usage, resource allocation, and Work Breakdown Structure (WBS) modelling for the Privacy-Preserving Electronic Voting System using StarUML, GanttProject, and Microsoft Project.

### 1. Project Overview & Scope
The **Privacy-Preserving Electronic Voting System (AegisVote)** is a secure, web-based electronic voting platform designed to guarantee ballot secrecy and cryptographic tamper-evidence without relying on blockchain consensus. The system solves the primary challenge of electronic voting: decoupling voter authentication from vote casting so that no single server or malicious database administrator can ever associate an authenticated voter's identity with their submitted candidate choice.

#### Core Architectural Innovations:
- **Decoupled Authentication:** Separates voter identity verification from ballot persistence using an independent, time-bound Anonymous Voting Token Service.
- **Single-Use Anonymous Tokens:** Each registered voter receives a single-use, unlinked cryptographic token authorizing exactly one vote per election.
- **Cryptographic Integrity:** Every vote payload is encrypted with AES-256 before disk storage; an append-only binary SHA-256 Merkle Tree ledger tracks ballot hashes in real time to instantly detect tampering or unauthorized ballot deletion.
- **Role-Based Governance:** Strict role boundaries for Voters, Election Officers, Administrators, and Counting Authorities with automated audit logging.

![GanttProject Timeline and Master Project Schedule](Report%20files/extracted_figures/Lab2_Gantt_Overview.png)

### 2. Work Breakdown Structure (WBS) Modelling
To ensure methodical execution across the software engineering lifecycle, the project was decomposed into hierarchical Work Breakdown Structures using four orthogonal decomposition criteria:
1. **Process-Based Decomposition:** Decomposes the project along the chronological software engineering lifecycle: Requirements Engineering → Architectural Modelling → Backend Cryptographic Core → Frontend Development → Integration & Penetration Testing → Deployment.
2. **Product-Based Decomposition:** Decomposes the project according to tangible architectural deliverables: Voter Portal UI, Authentication Gateway, Token Issuance Engine, AES-256 Ballot Encryption Service, PostgreSQL Merkle Ledger, and Election Counting Authority Dashboard.
3. **Geography-Based Decomposition:** Decomposes tasks across physical and cloud deployment zones: Regional Voter Client Terminals, Central Cloud Application Cluster (US-East), and Air-Gapped Verification Server for final cryptographic tallying.
4. **Role-Based Decomposition:** Decomposes work according to team specializations: Project Lead, Cryptography Specialist, Full-Stack Frontend Developer, Backend Go Developer, QA / Test Automation Engineer, and Database Administrator.

![WBS Decomposition Hierarchy (StarUML)](Report%20files/extracted_figures/Lab2_Planning_image16.png)
![Product-Based Work Breakdown Structure](Report%20files/extracted_figures/Lab2_Product_WBS.png)
![Geography-Based Distributed Architecture WBS](Report%20files/extracted_figures/Lab2_Geography_WBS.png)
![Role-Based Resource & Responsibility Decomposition](Report%20files/extracted_figures/Lab2_Role_WBS.png)

### 3. Activity Network, Task Scheduling & Duration Allocation
A Program Evaluation and Review Technique (PERT) / Critical Path Method (CPM) Activity Network was synthesized to model task precedence constraints, finish-to-start (FS) dependencies, and critical path milestones. Task usage, work effort (man-days), and duration allocations were formulated in GanttProject / Microsoft Project to ensure resource leveling and eliminate delivery bottlenecks.

![PERT / CPM Activity Network Diagram](Report%20files/extracted_figures/Lab2_Network_Diagram.png)
![Task Usage and Schedule Duration Allocation Matrix](Report%20files/extracted_figures/Lab2_Planning_image23.png)
![Resource Sheet and Personnel Allocation Matrix](Report%20files/extracted_figures/Lab2_Planning_image24.png)
![Resource Allocation Graph and Team Workload Distribution](Report%20files/extracted_figures/Lab2_Resource_Usage.png)

### 4. Result & Conclusion
The project scheduling, multi-criteria WBS decomposition, task network, and resource allocation models were successfully completed using StarUML, GanttProject, and Microsoft Project. The resulting schedule establishes a robust foundation for agile sprint tracking, milestone verification, and resource leveling across the 12-week development lifecycle.

---

## EXPERIMENT 3: Requirement Modelling: Context (Level-0), Level-1, and Level-2 DFDs

### Aim
To perform structured requirements analysis and construct Context Diagram (Level-0 DFD), Level-1 DFD, and Level-2 DFD blueprints for the Privacy-Preserving Electronic Voting System.

### 1. System Overview & Structured Analysis
Data Flow Diagrams (DFDs) provide a hierarchical, process-driven perspective of the voting system, mapping how data enters, transitions through processing units, and commits to persistent storage without exposing voter identities.
- **Data Separation & Privacy Wall:** Guarantees that no single process or data store ever contains both the voter's real-world identity and their candidate selection simultaneously.
- **Tamper-Evident Storage:** Maintains append-only ballot commitment verified through cryptographic hashing.
- **Strong Authentication:** Multi-factor authentication (Credentials + TOTP) verifying eligibility while keeping the ballot anonymous.

### 2. External Entities & Data Stores
- **External Entities:** Voter, Election Officer, Administrator, Counting Authority, and OTP Provider.
- **Data Stores:** D1: Users, D2: Elections, D3: Candidates, D4: Tokens, D5: Encrypted Ballots, D6: Audit Logs.

### 3. Context Diagram (Level-0 DFD) Blueprint
The Level-0 Context Diagram represents the entire voting platform as a single centralized process interacting with external actors.

![Context Diagram (Level-0 DFD) Blueprint](Report%20files/extracted_figures/DFD_Level0_Context.png)

### 4. Level-1 DFD: Functional Process Decomposition
The Level-1 DFD decomposes the system into five primary functional processes:
- **Process 1.0 (Voter Authentication & Eligibility Check):** Authenticates voter via credentials and TOTP, verifies active election eligibility, and checks the per-voter HasVoted flag.
- **Process 2.0 (Anonymous Token Issuance):** Issues a cryptographically random, single-use, time-bound anonymous token and flags the voter as having obtained a token.
- **Process 3.0 (Process & Encrypt Vote):** Validates token unicity, encrypts ballot payload using AES-256, appends to database, updates the Merkle tree root, and invalidates the token.
- **Process 4.0 (Vote Tallying & Integrity Verification):** Initiated post-election closure to verify Merkle root consistency, decrypt ballots, and compute final results.
- **Process 5.0 (Election Administration & Audit Logging):** Facilitates election setup, candidate registration, and system event logging.

![Level-1 DFD: Primary System Process Decomposition](Report%20files/extracted_figures/DFD_Level1_System.png)

### 5. Level-2 DFD: Process 3.0 (Process & Encrypt Vote Decomposition)
Level-2 DFD zooms directly into the high-security vote ingestion pipeline (Process 3.0):
- **Sub-Process 3.1 (Token Validation & Invalidation):** Validates the raw token against D4: Tokens and ensures `Used == False`.
- **Sub-Process 3.2 (Ballot Encryption):** Encrypts the candidate choice using AES-256-GCM symmetric encryption.
- **Sub-Process 3.3 (Merkle Tree Leaf Generation):** Computes SHA-256 leaf hash and recalculates the root of the Merkle Tree.
- **Sub-Process 3.4 (Append-Only Persistence):** Appends the encrypted ciphertext and leaf hash to D5: Encrypted Ballots.

![Level-2 DFD: Process 3.0 Decomposition](Report%20files/extracted_figures/DFD_Level2_Process3.png)

### 6. Result & Conclusion
The Context Diagram (Level-0), Level-1, and Level-2 Data Flow Diagrams were successfully developed. The models rigorously prove that voter identification information never flows into the ballot encryption or persistent storage pipelines, establishing complete architectural anonymity.

---

## EXPERIMENT 4: Requirement Modelling: Entity-Relationship (ER) Diagram & Relational Schema

### Aim
To design the Entity-Relationship (ER) Diagram, resolve Many-to-Many associations using bridge entities, formulate the complete data dictionary, and establish relational schema constraints enforcing zero voter-to-ballot linkage.

### 1. Problem Statement & Zero-Linkage Constraint
The objective is to design a relational database architecture that strictly preserves voter privacy:
- **Crucial Security Constraint 1:** The Vote entity intentionally contains NO UserID, VoterID, or foreign key referencing the User table.
- **Crucial Security Constraint 2:** The Token entity stores only a one-way cryptographic SHA-256 hash (TokenHash) with no direct link to the User record in persistent storage.
- **Crucial Security Constraint 3:** The Many-to-Many (M:N) relationship between Users and Elections is resolved via an associative bridge entity (`User_Election_Eligibility`) to track voter eligibility and the `HasVoted` flag independently for every election.

### 2. Comprehensive Data Dictionary Table

| Entity Name | Attribute | Attribute Type | Key / Integrity Constraint |
| :--- | :--- | :--- | :--- |
| **User** | UserID | UUID / Simple | Primary Key (PK) |
| User | Name | VARCHAR(100) | Not Null |
| User | PasswordHash | VARCHAR(255) | Bcrypt Hash, Not Null |
| User | OTPSecret | VARCHAR(64) | TOTP Shared Secret |
| User | Role | VARCHAR(20) | Enum ('VOTER', 'ADMIN', 'OFFICER', 'AUDITOR') |
| **Election** | ElectionID | UUID / Simple | Primary Key (PK) |
| Election | Title | VARCHAR(150) | Not Null |
| Election | Description | TEXT | Optional |
| Election | Status | VARCHAR(20) | Enum ('DRAFT', 'ACTIVE', 'CLOSED', 'TALLIED') |
| Election | StartTime | TIMESTAMP | Not Null |
| Election | EndTime | TIMESTAMP | Not Null |
| Election | TotalVotes | INTEGER | Derived Attribute (Computed from Vote count) |
| **Candidate** | CandidateID | UUID / Simple | Primary Key (PK) |
| Candidate | ElectionID | UUID / Simple | Foreign Key (FK References Election.ElectionID) |
| Candidate | CandidateName | VARCHAR(100) | Not Null |
| Candidate | Party | VARCHAR(100) | Not Null |
| **Token** | TokenID | UUID / Simple | Primary Key (PK) |
| Token | TokenHash | VARCHAR(64) | SHA-256 Hash of raw token, Unique, Not Null |
| Token | IssuedTime | TIMESTAMP | Not Null |
| Token | ExpiryTime | TIMESTAMP | Not Null (Time-bound 15-minute window) |
| Token | Used | BOOLEAN | Default False, Flipped to True upon vote |
| **Vote** | VoteID | UUID / Simple | Primary Key (PK) |
| Vote | EncryptedBallot | BYTEA / TEXT | AES-256-GCM Ciphertext Payload |
| Vote | Timestamp | TIMESTAMP | Commitment Timestamp |
| Vote | MerkleLeafHash | VARCHAR(64) | SHA-256 Hash of EncryptedBallot |
| **AuditLog** | LogID | UUID / Simple | Primary Key (PK) |
| AuditLog | Action | VARCHAR(100) | Event Description (e.g., 'TOKEN_ISSUED') |
| AuditLog | Timestamp | TIMESTAMP | Audit Event Timestamp |
| AuditLog | UserRole | VARCHAR(20) | Role of initiating identity |
| AuditLog | IPAddress | VARCHAR(45) | IPv4 / IPv6 Client Address |
| **User_Election_Eligibility** | UserID | UUID / Simple | Composite PK, FK (References User.UserID) |
| User_Election_Eligibility | ElectionID | UUID / Simple | Composite PK, FK (References Election.ElectionID) |
| User_Election_Eligibility | HasVoted | BOOLEAN | Default False, Tracks participation per election |

### 3. Cardinality & Relationship Constraints
- **Election to Candidate (1:M):** One Election contains multiple Candidates. Partial participation for Election; total participation for Candidate.
- **User to Election (Resolved M:N):** Users are eligible for multiple Elections, and an Election includes multiple eligible Users. Resolved through `User_Election_Eligibility`.
- **Token to Vote (Conceptual 1:1):** A single-use token authorizes exactly one vote. To preserve privacy, no foreign key connects Token and Vote in physical storage.
- **System to AuditLog (1:M):** System operations generate sequential, append-only security logs.

### 4. Entity-Relationship Diagram

![Final Relational ER Diagram with Zero-PII Voting Architecture](Report%20files/extracted_figures/ER_diagram_image1.png)

### 5. Result & Conclusion
The Entity-Relationship Diagram and relational schema were successfully formulated and normalized. The architecture guarantees mathematical and structural impossibility of linking individual voter identities to decrypted ballots, ensuring unconditional ballot confidentiality.

---

## EXPERIMENT 5: Behavioral Modelling: State Transition Diagram

### Aim
To model the dynamic behavioral lifecycle, states, events, guard conditions, and cryptographic transitions of the Privacy-Preserving Electronic Voting System using UML State Transition Diagrams.

### 1. Behavioral Modelling Concept
UML State Transition Diagrams capture the finite-state machine (FSM) governing a system entity over its operational lifetime. In AegisVote, behavioral modelling ensures that a voter's session and voting token progress through strictly validated, irreversible stages, preventing replay attacks and double voting.

### 2. State Transition Diagram Blueprint

![UML State Transition Diagram: Voter Session and Cryptographic Lifecycle States](Report%20files/extracted_figures/State_Transition_image1.png)

### 3. Detailed State Lifecycle & Transition Analysis
The dynamic machine transitions across nine discrete states governed by rigorous guard conditions:
1. **Idle / Unauthenticated:** The initial state before user interaction.
2. **Dual-Factor Authenticating:** Triggered by `submitCredentials()`. System verifies Bcrypt password hash and prompts for 2FA TOTP.
3. **Eligibility Verification:** Entered upon successful OTP validation. System evaluates guard condition `[HasVoted == False && ElectionStatus == 'ACTIVE']`. If `HasVoted == True`, transitions to Access Denied.
4. **Token Generated & Active:** Generated upon successful eligibility confirmation. TokenService issues a cryptographically random 256-bit token; `User_Election_Eligibility.HasVoted` is immediately updated to True.
5. **Ballot Selection:** Voter UI displays the candidate ballot. Governed by a strict 15-minute countdown guard `[ExpiryTime > CurrentTime]`.
6. **Ballot Encrypted:** Triggered by `castVote(token, candidateId)`. VotingServer invokes AES-256 encryption engine to encrypt candidate selection in memory.
7. **Merkle Tree Commitment:** SHA-256 leaf hash is computed and inserted into the append-only Merkle Tree engine, updating the global Merkle root.
8. **Token Invalidation:** The Token is permanently invalidated (`Used = True`), and the raw token is purged from memory.
9. **Receipt Dispatched & Terminated:** Voter receives a cryptographic receipt containing the Merkle leaf hash and timestamp; session terminates securely.

### 4. Result & Conclusion
The State Transition Diagram was successfully designed and validated in StarUML. The formal model proves that token consumption and ballot commitment occur within an atomic, irreversible sequence, eliminating token replay attacks and ballot tampering.

---

## EXPERIMENT 6: Object-Oriented Design: Use Case and Class Diagrams

### Aim
To perform Object-Oriented Analysis and Design (OOAD) by developing the Use Case Diagram and Class Diagram for the Privacy-Preserving Electronic Voting System using StarUML.

### 1. Identified Actors & Roles
- **Voter:** Primary user who authenticates with multi-factor credentials, requests an anonymous token, casts one vote, and reviews election results.
- **Election Officer:** Administrative authority responsible for creating elections, registering candidates, opening/closing polls, and certifying results.
- **Administrator:** Maintains system uptime, provisions accounts, reviews audit logs, and initiates database backups without access to ballot data.
- **Counting Authority:** Trusted entity executing cryptographic Merkle root audits, post-election ballot decryption, and tally calculation.
- **OTP Provider:** External SaaS provider delivering out-of-band TOTP codes via SMS or Authenticator protocols.

### 2. Use Case Model & Relationships
- **Voter Use Cases:** Authenticate Voter, Authenticate with OTP, Obtain Anonymous Voting Token, Cast Vote, Validate Anonymous Token, Encrypt Ballot, Store Encrypted Ballot, Check Election Status, View Results.
- **Election Officer Use Cases:** Create Election, Manage Candidates, Open Election, Close Election, Publish Results.
- **Administrator Use Cases:** Manage User Accounts, Backup Database, View Audit Logs.
- **Counting Authority Use Cases:** Verify Merkle Integrity, Decrypt Ballots, Tally Votes, Generate Final Results Report.
- **Key Dependencies:** `Authenticate Voter` <<includes>> `Authenticate with OTP`; `Cast Vote` <<includes>> `Validate Anonymous Token`, `Encrypt Ballot`, and `Store Encrypted Ballot`; `Tally Votes` <<includes>> `Verify Merkle Integrity` and `Decrypt Ballots`; `Account Lockout` <<extends>> `Authenticate Voter` (triggered upon repeated authentication failures exceeding threshold).

![StarUML Use Case Diagram: Actors, Use Cases, Include & Extend Dependencies](Report%20files/extracted_figures/Use_Case_and_Class_image1.png)

### 3. Class Diagram & Domain Model
The Class Diagram models the static structural architecture, encapsulating attributes, operations, and relationships across 14 domain classes:
- **User & Role Classes:** `User` (base class), `Voter`, `ElectionOfficer`, `Administrator`, `CountingAuthority`. Generalization hierarchy where specific roles inherit authentication properties.
- **Election Entities:** `Election` (composed of Candidates and Ballots) and `Candidate`.
- **Security & Voting Entities:** `AnonymousVotingToken` (handles token hashing and unicity) and `Ballot` (stores encrypted payload and Merkle leaf hash without any voter ID).
- **Service Classes:** `AuthenticationServer`, `TokenService`, `VotingServer`, `VoteEncryption`, `MerkleTreeIntegrity`, `CountingService`, `AuditLog`.

![StarUML Class Diagram: Entity Relationships, Attributes, Operations & Multiplicities](Report%20files/extracted_figures/Use_Case_and_Class_image2.png)

### 4. Object-Oriented Design Principles Applied
- **Encapsulation:** Internal representation of sensitive attributes (password hashes, raw tokens, private keys) is strictly hidden.
- **Abstraction:** High-level services expose abstract interfaces, decoupling business logic from underlying cryptographic libraries.
- **Inheritance:** Voter, ElectionOfficer, Administrator, and CountingAuthority share common credentials from the User superclass.
- **Single Responsibility Principle (SRP):** Each service maintains a single clear purpose (TokenService manages tokens; MerkleTreeIntegrity manages cryptographic ledgers).

### 5. Result & Conclusion
The Use Case Diagram and Class Diagram were successfully constructed in StarUML. The object-oriented model cleanly separates responsibilities, decouples voting operations from identity management, and provides a scalable design for implementation.

---

## EXPERIMENT 7: Interaction Modelling: Sequence Diagram and Communication Diagram

### Aim
To model dynamic object interactions, temporal execution ordering, and structural collaboration networks of the Privacy-Preserving Electronic Voting System using UML Sequence and Communication Diagrams.

### 1. Sequence Diagram: Temporal Interaction
The Sequence Diagram models the chronological execution timeline across six interacting lifelines: `Voter`, `:VoterClient`, `:VotingServer`, `:TokenService`, `:EncryptionModule`, and `:ElectionDatabase`.
- **Phase 1 (User Intent & Authorization):** Voter selects candidate; UI presents confirmation dialog. Upon confirmation, `:VoterClient` transmits `submitBallot(token, candidateId)` to `:VotingServer`. `:VotingServer` queries `:TokenService` with `validateToken(tokenHash)`.
- **Phase 2 (Ballot Encryption & Ledger Commitment):** `:VotingServer` passes candidateId to `:EncryptionModule` in memory to generate an AES-256 ciphertext. `:VotingServer` invokes `appendEncryptedBallot()` on `:ElectionDatabase` to store the ciphertext and trigger Merkle tree root recalculation.
- **Phase 3 (Sanitization & Receipt Delivery):** `:ElectionDatabase` returns commitment confirmation; `:VotingServer` calls `invalidateToken()` on `:TokenService` and delivers a cryptographic receipt to `:VoterClient`.

![UML Sequence Diagram: Chronological Vote Casting and Cryptographic Commitment](Report%20files/extracted_figures/Sequence_and_Comm_image1.png)

### 2. Communication Diagram: Structural Collaboration
While the Sequence Diagram emphasizes the time axis, the Communication Diagram models object topology and message pathways using numbered messaging discipline (messages 1 through 12).
- **Mediator Architecture:** The Communication Diagram demonstrates that `:VotingServer` acts as a Mediator, decoupling client requests from cryptographic engines and persistent databases.
- **Numbered Messaging Discipline:** Synchronous forward invocations (e.g., `5: validateToken`, `7: encryptBallot`, `9: appendEncryptedBallot`) are paired with corresponding return messages (`6: tokenValid`, `8: encryptedBallot`, `10: commitmentSuccess`) tracing back along the collaboration links.

![UML Communication Diagram: Object Collaborations and Numbered Message Discipline](Report%20files/extracted_figures/Sequence_and_Comm_image2.png)

### 3. Result & Conclusion
The Sequence and Communication Diagrams were successfully designed in StarUML. The dual perspectives rigorously confirm both the chronological execution flow and structural decoupling of the voting pipeline.

---

## EXPERIMENT 8: OO Design: Package, Component, and Deployment Models

### Aim
To design the Package Diagram, Component Diagram, and Deployment Diagram for the Privacy-Preserving Electronic Voting System using StarUML.

### 1. Package Diagram: Logical Module Organization
The Package Diagram groups related model elements into five cohesive, loosely coupled logical packages:
- **Authentication and Privacy Package:** Encapsulates the Authentication Server, TOTP validation module, and Anonymous Token Service.
- **Voting and Encryption Package:** Houses the Voting Server and AES-256 Vote Encryption Engine.
- **Election Data and Integrity Package:** Contains the PostgreSQL database repository and SHA-256 Merkle Tree Integrity Module.
- **Counting and Results Package:** Encapsulates post-election decryption and tallying services.
- **Audit and Administration Package:** Contains the audit logging engine and administrative interfaces.

![StarUML Package Diagram: Logical Organization and Module Dependencies](Report%20files/extracted_figures/Package_Component_Deploy_image1.png)

### 2. Component Diagram: Service Architecture & Interfaces
- **Components:** Web Client, Authentication Server, Anonymous Token Service, Voting Server, Vote Encryption Module, Election Database, Merkle Tree Integrity Module, Counting Authority, Audit Logging Module, and External OTP Provider.
- **Interfaces:** Components communicate via lightweight RESTful HTTPS interfaces. The Voting Server consumes interfaces from Token Service and Encryption Module but operates completely detached from Authentication Server.

![StarUML Component Diagram: Decoupled Microservice Components and REST Interfaces](Report%20files/extracted_figures/Package_Component_Deploy_image2.png)

### 3. Deployment Diagram: Physical Infrastructure & Execution Environments
- **Voter Device Node:** Executes modern Web Browsers (Chrome, Edge, Firefox) accessing the system over HTTPS/TLS 1.3.
- **Application Server Node:** Hosts the Go backend runtime, running containerized microservices behind a secure NGINX reverse proxy.
- **Database Server Node:** Hosts the PostgreSQL relational database engine and cryptographic append-only Merkle ledger.
- **Counting Authority Console:** Air-gapped secure hardware console utilized solely post-closure for ballot decryption and tally computation.

![StarUML Deployment Diagram: Physical Client, Gateway Server, and Database Nodes](Report%20files/extracted_figures/Package_Component_Deploy_image3.png)

### 4. Result & Conclusion
The Package, Component, and Deployment Diagrams were successfully constructed. The models represent the complete logical organization, microservice boundaries, physical infrastructure, and network communication paths of the system.

---

## EXPERIMENT 9: Software Testing: Design and Demonstration of Test Cases

### Aim
To design, demonstrate, and execute functional automated test cases using Selenium IDE and non-functional performance / load test suites using Apache JMeter on the target web application (BlazeDemo).

### 1. Testing Strategy & Target Application
Quality Assurance (QA) verifies that software functions reliably under expected operations and scales under concurrent user loads. The testing suite was executed on **BlazeDemo** (https://blazedemo.com), a web-based flight reservation simulator:
- **Functional Testing:** Selenium IDE browser automation validating departure/destination routing, flight selection, passenger data entry, and confirmation receipt generation.
- **Non-Functional Testing:** Apache JMeter thread group simulation evaluating server response times, latency, throughput, and error rates under concurrent user load.

### 2. Master Test Case Suite Table

| Test ID | Test Type | Scenario Description | Steps to Execute | Expected Result | Actual Result | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **TC_01** | Functional | Verify Flight Search Functionality | 1. Navigate to blazedemo.com<br>2. Select Departure City (Paris)<br>3. Select Destination City (Rome)<br>4. Click 'Find Flights' | User is redirected to flights list page displaying available flights for route. | Flights list page rendered matching search parameters. | **PASS** |
| **TC_02** | Functional | Verify End-to-End Flight Booking | 1. From flights list, click 'Choose This Flight'<br>2. Fill passenger details (Name, Address, CC)<br>3. Click 'Purchase Flight' | User is redirected to confirmation page displaying 'Thank you for your purchase today!' with unique booking ID. | Confirmation page displayed with valid booking receipt ID. | **PASS** |
| **TC_03** | Non-Functional | Basic Load Test (10 Concurrent Users) | 1. Configure JMeter Thread Group: 10 threads, 5s ramp-up period, 1 loop<br>2. Send HTTP GET requests to blazedemo.com<br>3. Capture Aggregate Report & Response Graphs | Server responds without crashes or timeouts; 0.00% error rate; average response time < 1000 ms. | All 10 samples completed successfully; 0.00% error rate; average latency 890 ms. | **PASS** |

![Master Test Case Suite Table and Verification Status](Report%20files/extracted_figures/Lab9_Testing_img1.png)

### 3. Functional Testing Demonstration (Selenium IDE)
Functional test cases TC_01 and TC_02 were recorded and executed in Selenium IDE. The test runner executed all locator assertions (ID, XPath, CSS selectors), command events, and page transitions with zero failures.

#### TC_01: Flight Search Automation Walkthrough
The automated script navigated to the portal, selected 'Paris' as origin and 'Rome' as destination, and clicked the search button. The application returned valid flight options, confirming the routing filter logic.

![TC_01: Selenium IDE Flight Search Execution Log](Report%20files/extracted_figures/Lab9_Testing_img2.png)

#### TC_02: Flight Booking & Purchase Confirmation Walkthrough
The automation script selected the first available carrier, populated form fields with passenger credentials, and dispatched the payment request. The confirmation screen verified successful transaction completion.

![TC_02: Selenium IDE Flight Purchase Confirmation](Report%20files/extracted_figures/Lab9_Testing_img3.png)

### 4. Non-Functional Performance & Load Testing (Apache JMeter)
Test case TC_03 evaluated the target application's responsiveness under concurrent traffic. A JMeter Thread Group was instantiated with 10 simulated users ramping up over a 5-second interval.

![TC_03: Apache JMeter Thread Group Load Configuration](Report%20files/extracted_figures/Lab9_Testing_img4.png)

#### JMeter Aggregate Report & Throughput Analysis
- **Total Samples:** 10 HTTP GET requests dispatched across the 5-second ramp-up window.
- **Response Time Metrics:** Average response time observed was 890 ms, with median latency at 860 ms and 90% Line at 980 ms.
- **Error Rate:** 0.00% error rate across all threads; zero HTTP 5xx server exceptions or 4xx client errors.
- **Throughput:** Sustainable throughput of 1.9 transactions/second with continuous stream delivery.

![TC_03: Apache JMeter Aggregate Report](Report%20files/extracted_figures/Lab9_Testing_img5.png)
![TC_03: Apache JMeter Response Time Graph & Latency Distribution Curve](Report%20files/extracted_figures/Lab9_Testing_img6.png)

### 5. Result & Conclusion
Both functional test cases (TC_01, TC_02) and the non-functional performance test case (TC_03) were executed and validated. The target application demonstrated 100% test pass rate in Selenium IDE and maintained 0.00% error rate under concurrent load in Apache JMeter.