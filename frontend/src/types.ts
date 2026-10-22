export type Persona = 'voter1' | 'voter2' | 'voter3' | 'voter4' | 'officer1' | 'counting1' | 'admin1';

export interface User {
  user_id: string;
  name: string;
  role: 'Voter' | 'ElectionOfficer' | 'CountingAuthority' | 'Administrator';
  has_voted?: boolean;
}

export interface Candidate {
  candidate_id: string;
  election_id: string;
  candidate_name: string;
  party: string;
}

export interface Election {
  election_id: string;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  status: 'Created' | 'Active' | 'Closed' | 'Archived';
  public_key: string;
  signed_merkle_root?: string;
}

export interface AuditLog {
  log_id: string;
  action: string;
  timestamp: string;
  user_role: string;
  ip_address: string;
}

export interface MerkleStatus {
  status: 'in_progress' | 'verified' | 'tampered' | 'error';
  ballots_count: number;
  computed_root: string;
  signed_merkle_root?: string;
  leaves?: string[];
  verification_pass: boolean;
  message: string;
}

export interface BallotReceipt {
  receipt_h: string;
  root_hash: string;
  timestamp?: string;
  candidate_name?: string;
  party?: string;
}

export interface JiraTestCase {
  key: string;
  summary: string;
  requirement: string;
  category: 'AUTH' | 'TOKEN' | 'CRYPTO' | 'LEDGER' | 'LIFECYCLE' | 'AUDIT';
  status: 'IDLE' | 'RUNNING' | 'PASS' | 'FAIL';
  executionTimeMs?: number;
  requestPayload?: any;
  responsePayload?: any;
  assertionLog?: string[];
}
