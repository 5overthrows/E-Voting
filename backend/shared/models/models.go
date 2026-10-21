package models

import "time"

type User struct {
	UserID             string     `json:"user_id"`
	Name               string     `json:"name"`
	PasswordHash       string     `json:"-"`
	HasVoted           bool       `json:"has_voted"`
	Role               string     `json:"role"` // Voter, ElectionOfficer, Administrator, CountingAuthority
	FailedLoginAttempts int        `json:"-"`
	LockedUntil        *time.Time `json:"-"`
}

type Election struct {
	ElectionID          string     `json:"election_id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	StartTime           time.Time  `json:"start_time"`
	EndTime             time.Time  `json:"end_time"`
	Status              string     `json:"status"` // Created, Active, Closed, Archived
	PublicKey           string     `json:"public_key"`
	EncryptedPrivateKey string     `json:"-"`
	PrivateKeySalt      string     `json:"-"`
	SignedMerkleRoot    string     `json:"signed_merkle_root,omitempty"`
}

type Candidate struct {
	CandidateID   string `json:"candidate_id"`
	ElectionID    string `json:"election_id"`
	CandidateName string `json:"candidate_name"`
	Party         string `json:"party"`
}

type Token struct {
	TokenID    string    `json:"token_id"`
	TokenHash  string    `json:"token_hash"`
	IssuedTime time.Time `json:"issued_time"`
	ExpiryTime time.Time `json:"expiry_time"`
	Used       bool      `json:"used"`
}

type Vote struct {
	VoteID          string    `json:"vote_id"`
	EncryptedBallot string    `json:"encrypted_ballot"`
	Timestamp       time.Time `json:"timestamp"`
	MerkleLeafHash  string    `json:"merkle_leaf_hash"`
}

type AuditLog struct {
	LogID     string    `json:"log_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	UserRole  string    `json:"user_role"`
	IPAddress string    `json:"ip_address"`
}
