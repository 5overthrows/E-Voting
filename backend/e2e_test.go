package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"evoting/services/audit"
	"evoting/services/auth"
	"evoting/services/counting"
	"evoting/services/demo"
	"evoting/services/token"
	"evoting/services/voting"
	"evoting/shared/db"

	"github.com/gin-gonic/gin"
)

type ginTransport struct {
	engine *gin.Engine
}

func (t *ginTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.engine.ServeHTTP(w, req)
	return w.Result(), nil
}

func TestE2EVotingPipeline(t *testing.T) {
	// Clean up any existing test db
	os.Remove("test_evoting.db")
	os.Setenv("SQLITE_PATH", "test_evoting.db")
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Remove("test_evoting.db")

	testBase := "http://aegisvote.local"
	os.Setenv("TOKEN_SERVICE_URL", testBase+"/api/token/generate")
	os.Setenv("TOKEN_VALIDATE_URL", testBase+"/api/token/validate")
	os.Setenv("TOKEN_INVALIDATE_URL", testBase+"/api/token/invalidate")
	os.Setenv("AUTH_CONFIRM_URL", testBase+"/api/auth/confirm-vote")

	// 1. Initialize Database
	db.Connect()

	// Configure in-memory test engine
	gin.SetMode(gin.TestMode)
	r := gin.New()
	auth.RegisterRoutes(r)
	token.RegisterRoutes(r)
	voting.RegisterRoutes(r)
	counting.RegisterRoutes(r)
	audit.RegisterRoutes(r)
	demo.RegisterRoutes(r)

	// Direct all HTTP client calls (including internal microservice calls) in-memory
	transport := &ginTransport{engine: r}
	http.DefaultTransport = transport
	client := &http.Client{Transport: transport}

	// 2. Login as Election Officer
	t.Log("Step 1: Logging in as Election Officer...")
	officerLoginBody, _ := json.Marshal(map[string]string{
		"user_id":  "officer1",
		"password": "password123",
	})
	resp, err := client.Post(testBase+"/api/auth/login", "application/json", bytes.NewBuffer(officerLoginBody))
	if err != nil {
		t.Fatalf("Officer login failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Officer login returned status %d: %s", resp.StatusCode, string(b))
	}
	var officerAuth struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&officerAuth)
	t.Log("Officer authenticated. JWT acquired.")

	// 3. Create Election
	t.Log("Step 2: Creating new election...")
	passphrase := "super_secret_authority_passphrase"
	createElectBody, _ := json.Marshal(map[string]interface{}{
		"title":       "University Student Union Election 2026",
		"description": "General election for student representatives",
		"start_time":  time.Now().Add(-1 * time.Hour),
		"end_time":    time.Now().Add(24 * time.Hour),
		"passphrase":  passphrase,
	})
	req, _ := http.NewRequest("POST", testBase+"/api/elections/create", bytes.NewBuffer(createElectBody))
	req.Header.Set("Authorization", "Bearer "+officerAuth.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create election failed with status %d: %s", resp.StatusCode, string(b))
	}
	var electResp struct {
		ElectionID string `json:"election_id"`
	}
	json.NewDecoder(resp.Body).Decode(&electResp)
	electionID := electResp.ElectionID
	t.Logf("Election created with ID: %s", electionID)

	// 4. Add Candidates
	t.Log("Step 3: Registering candidates...")
	candidates := []struct {
		Name  string
		Party string
	}{
		{"Alice Chen", "Progressive Alliance"},
		{"Bob Smith", "Independent Students"},
	}
	var candIDs []string
	for _, c := range candidates {
		candBody, _ := json.Marshal(map[string]string{
			"name":  c.Name,
			"party": c.Party,
		})
		cReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/elections/%s/candidates", testBase, electionID), bytes.NewBuffer(candBody))
		cReq.Header.Set("Authorization", "Bearer "+officerAuth.Token)
		cReq.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(cReq)
		if err != nil || resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("Add candidate %s failed: %s", c.Name, string(b))
		}
		var cResp struct {
			CandidateID string `json:"candidate_id"`
		}
		json.NewDecoder(resp.Body).Decode(&cResp)
		candIDs = append(candIDs, cResp.CandidateID)
		t.Logf("Added candidate %s with ID %s", c.Name, cResp.CandidateID)
	}

	// 5. Open Election
	t.Log("Step 4: Opening election...")
	openReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/elections/%s/open", testBase, electionID), nil)
	openReq.Header.Set("Authorization", "Bearer "+officerAuth.Token)
	resp, err = client.Do(openReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Open election failed: %s", string(b))
	}
	t.Log("Election is now Active.")

	// 6. Voter 1 votes for Alice Chen (candIDs[0])
	t.Log("Step 5: Voter 1 login & casting vote...")
	voterLoginBody, _ := json.Marshal(map[string]string{
		"user_id":  "voter1",
		"password": "password123",
	})
	resp, err = client.Post(testBase+"/api/auth/login", "application/json", bytes.NewBuffer(voterLoginBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Voter 1 login failed")
	}
	var voter1Auth struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&voter1Auth)

	// Request Token
	tokReqBody, _ := json.Marshal(map[string]string{
		"election_id": electionID,
	})
	tokReq, _ := http.NewRequest("POST", testBase+"/api/auth/token/request", bytes.NewBuffer(tokReqBody))
	tokReq.Header.Set("Authorization", "Bearer "+voter1Auth.Token)
	tokReq.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(tokReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Voter 1 token request failed: %s", string(b))
	}
	var tokResp1 struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokResp1)
	anonymousToken1 := tokResp1.Token
	t.Logf("Voter 1 received anonymous single-use token: %s", anonymousToken1[:12]+"...")

	// Cast Vote
	voteCastBody, _ := json.Marshal(map[string]string{
		"token":        anonymousToken1,
		"candidate_id": candIDs[0],
		"election_id":  electionID,
	})
	resp, err = client.Post(testBase+"/api/vote/cast", "application/json", bytes.NewBuffer(voteCastBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Voter 1 cast vote failed: %s", string(b))
	}
	var castResult1 struct {
		ReceiptH string `json:"receipt_h"`
		RootHash string `json:"root_hash"`
	}
	json.NewDecoder(resp.Body).Decode(&castResult1)
	t.Logf("Voter 1 vote confirmed! Leaf Hash: %s, Current Merkle Root: %s", castResult1.ReceiptH, castResult1.RootHash)

	// 7. Test Double Voting Prevention (Voter 1 tries to request token again)
	t.Log("Step 6: Verifying double-voting prevention...")
	tokReqAgain, _ := http.NewRequest("POST", testBase+"/api/auth/token/request", bytes.NewBuffer(tokReqBody))
	tokReqAgain.Header.Set("Authorization", "Bearer "+voter1Auth.Token)
	tokReqAgain.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(tokReqAgain)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden on second token request, got %d", resp.StatusCode)
	}
	t.Log("Double voting prevented: voter1 has_voted flag verified.")

	// Test Token Reuse (try voting with anonymousToken1 again)
	resp, err = client.Post(testBase+"/api/vote/cast", "application/json", bytes.NewBuffer(voteCastBody))
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized on reused token, got %d", resp.StatusCode)
	}
	t.Log("Token reuse prevented: single-use token invalidated.")

	// 8. Voter 2 votes for Alice Chen as well
	t.Log("Step 7: Voter 2 login & casting vote...")
	voter2LoginBody, _ := json.Marshal(map[string]string{
		"user_id":  "voter2",
		"password": "password123",
	})
	resp, _ = client.Post(testBase+"/api/auth/login", "application/json", bytes.NewBuffer(voter2LoginBody))
	var voter2Auth struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&voter2Auth)

	tokReq2, _ := http.NewRequest("POST", testBase+"/api/auth/token/request", bytes.NewBuffer(tokReqBody))
	tokReq2.Header.Set("Authorization", "Bearer "+voter2Auth.Token)
	tokReq2.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(tokReq2)
	var tokResp2 struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokResp2)

	voteCastBody2, _ := json.Marshal(map[string]string{
		"token":        tokResp2.Token,
		"candidate_id": candIDs[0], // Also votes for Alice
		"election_id":  electionID,
	})
	resp, _ = client.Post(testBase+"/api/vote/cast", "application/json", bytes.NewBuffer(voteCastBody2))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Voter 2 vote failed")
	}
	t.Log("Voter 2 vote successfully cast.")

	// 9. Officer closes election
	t.Log("Step 8: Officer closes election...")
	closeReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/elections/%s/close", testBase, electionID), nil)
	closeReq.Header.Set("Authorization", "Bearer "+officerAuth.Token)
	resp, err = client.Do(closeReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Close election failed: %s", string(b))
	}
	var closeResp struct {
		MerkleRoot       string `json:"merkle_root"`
		SignedMerkleRoot string `json:"signed_merkle_root"`
	}
	json.NewDecoder(resp.Body).Decode(&closeResp)
	t.Logf("Election Closed. Signed Merkle Root: %s", closeResp.SignedMerkleRoot[:30]+"...")

	// 10. Verify Merkle Tree Integrity
	t.Log("Step 9: Verifying Merkle Tree Ledger integrity...")
	verifyReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/elections/%s/verify-merkle", testBase, electionID), nil)
	verifyReq.Header.Set("Authorization", "Bearer "+officerAuth.Token)
	resp, _ = client.Do(verifyReq)
	var verifyResp struct {
		Status           string `json:"status"`
		VerificationPass bool   `json:"verification_pass"`
	}
	json.NewDecoder(resp.Body).Decode(&verifyResp)
	if !verifyResp.VerificationPass {
		t.Fatalf("Merkle Tree verification failed!")
	}
	t.Log("Merkle Tree verified intact! 2 ballots in ledger.")

	// 11. Counting Authority logs in and tallies
	t.Log("Step 10: Counting Authority tallies decrypted votes...")
	countingLoginBody, _ := json.Marshal(map[string]string{
		"user_id":  "counting1",
		"password": "password123",
	})
	resp, _ = client.Post(testBase+"/api/auth/login", "application/json", bytes.NewBuffer(countingLoginBody))
	var countingAuth struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&countingAuth)

	// Wrong passphrase test
	tallyReqWrong, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/elections/%s/tally", testBase, electionID), bytes.NewBuffer([]byte(`{"passphrase":"wrong_pass"}`)))
	tallyReqWrong.Header.Set("Authorization", "Bearer "+countingAuth.Token)
	tallyReqWrong.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(tallyReqWrong)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 with wrong passphrase, got %d", resp.StatusCode)
	}
	t.Log("Passphrase protection verified: incorrect passphrase rejected.")

	// Correct passphrase tally
	tallyReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/elections/%s/tally", testBase, electionID), bytes.NewBuffer([]byte(fmt.Sprintf(`{"passphrase":"%s"}`, passphrase))))
	tallyReq.Header.Set("Authorization", "Bearer "+countingAuth.Token)
	tallyReq.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(tallyReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Tally failed with status %d: %s", resp.StatusCode, string(b))
	}
	var tallyResult struct {
		TotalVotes int `json:"total_votes"`
		Results    []struct {
			CandidateName string `json:"candidate_name"`
			Votes         int    `json:"votes"`
		} `json:"results"`
	}
	json.NewDecoder(resp.Body).Decode(&tallyResult)
	t.Logf("Tally completed! Total votes: %d", tallyResult.TotalVotes)
	for _, r := range tallyResult.Results {
		t.Logf("Candidate: %s - %d votes", r.CandidateName, r.Votes)
	}
	if tallyResult.TotalVotes != 2 {
		t.Fatalf("Expected 2 total votes, got %d", tallyResult.TotalVotes)
	}

	// 12. Check Results Public Endpoint
	t.Log("Step 11: Checking public results endpoint...")
	resp, err = client.Get(fmt.Sprintf("%s/api/elections/%s/results", testBase, electionID))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Public results fetch failed")
	}
	t.Log("Public results endpoint verified.")

	// 13. Administrator Audit Logs & Backup
	t.Log("Step 12: Testing Administrator Audit Logs & Backup...")
	adminLoginBody, _ := json.Marshal(map[string]string{
		"user_id":  "admin1",
		"password": "password123",
	})
	resp, _ = client.Post(testBase+"/api/auth/login", "application/json", bytes.NewBuffer(adminLoginBody))
	var adminAuth struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&adminAuth)

	// Fetch logs
	logsReq, _ := http.NewRequest("GET", testBase+"/api/audit/logs", nil)
	logsReq.Header.Set("Authorization", "Bearer "+adminAuth.Token)
	resp, _ = client.Do(logsReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Audit logs fetch failed")
	}
	var logs []interface{}
	json.NewDecoder(resp.Body).Decode(&logs)
	t.Logf("Administrator verified: %d audit events logged.", len(logs))

	// Trigger Backup
	backupReq, _ := http.NewRequest("POST", testBase+"/api/audit/backup", nil)
	backupReq.Header.Set("Authorization", "Bearer "+adminAuth.Token)
	resp, _ = client.Do(backupReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Backup failed")
	}
	t.Log("Database backup completed successfully.")

	// 14. Demo Endpoints Test (Reset & Tamper Simulation)
	t.Log("Step 13: Testing Demo & Tamper Simulation endpoints...")
	tamperBody, _ := json.Marshal(map[string]string{"election_id": electionID})
	tamperResp, err := client.Post(testBase+"/api/demo/tamper", "application/json", bytes.NewBuffer(tamperBody))
	if err != nil || tamperResp.StatusCode != http.StatusOK {
		t.Fatalf("Demo tamper failed: %v", err)
	}

	// Now check Merkle verification: should fail with tampered status
	verifyReqTamper, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/elections/%s/verify-merkle", testBase, electionID), nil)
	verifyReqTamper.Header.Set("Authorization", "Bearer "+officerAuth.Token)
	vResp, _ := client.Do(verifyReqTamper)
	var tamperVerify struct {
		Status           string `json:"status"`
		VerificationPass bool   `json:"verification_pass"`
	}
	json.NewDecoder(vResp.Body).Decode(&tamperVerify)
	if tamperVerify.VerificationPass || tamperVerify.Status != "tampered" {
		t.Fatalf("Expected tamper verification to fail with status 'tampered', got pass=%v status=%s", tamperVerify.VerificationPass, tamperVerify.Status)
	}
	t.Log("Tamper simulation successfully caught by Merkle Tree verification!")

	// Restore
	restoreResp, _ := client.Post(testBase+"/api/demo/restore", "application/json", bytes.NewBuffer(tamperBody))
	if restoreResp.StatusCode != http.StatusOK {
		t.Fatalf("Restore failed")
	}
	t.Log("Tampered ballot successfully restored to authentic cryptographic state.")

	t.Log("==========================================================")
	t.Log("ALL E2E PRIVACY-PRESERVING VOTING CHECKS PASSED PERFECTLY!")
	t.Log("==========================================================")
}
