import React, { useState } from 'react';
import type { User, Election, Candidate, BallotReceipt } from '../types';

interface VoterViewProps {
  user: User;
  elections: Election[];
  receipt: BallotReceipt | null;
  loading: boolean;
  onSelectElectionForVoting: (election: Election) => void;
  onViewResults: (election: Election) => void;
  onCastVote: (electionId: string, candidateId: string) => Promise<boolean>;
  candidates: Candidate[];
  selectedElection: Election | null;
  activeScreen: string;
  onBackToElections: () => void;
}

export const VoterView: React.FC<VoterViewProps> = ({
  user,
  elections,
  receipt,
  loading,
  onSelectElectionForVoting,
  onViewResults,
  onCastVote,
  candidates,
  selectedElection,
  activeScreen,
  onBackToElections,
}) => {
  const [selectedCandidate, setSelectedCandidate] = useState<Candidate | null>(null);
  const [confirmModalOpen, setConfirmModalOpen] = useState(false);
  const [copiedReceipt, setCopiedReceipt] = useState(false);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedReceipt(true);
    setTimeout(() => setCopiedReceipt(false), 3000);
  };

  const handleConfirmAndCast = async () => {
    if (!selectedElection || !selectedCandidate) return;
    const ok = await onCastVote(selectedElection.election_id, selectedCandidate.candidate_id);
    if (ok) {
      setConfirmModalOpen(false);
      setSelectedCandidate(null);
    }
  };

  return (
    <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      {/* VOTER HERO CARD */}
      <div 
        className="glass-card" 
        style={{ padding: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem' }}
        data-testid="voter-welcome-card"
      >
        <div>
          <span style={{ fontSize: '0.8rem', color: 'var(--secondary)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Voter Portal (Anonymous Token Architecture)
          </span>
          <h2 style={{ fontSize: '1.7rem', marginTop: '0.2rem' }}>Welcome, {user.name}!</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', marginTop: '0.3rem' }}>
            {user.has_voted 
              ? '✅ Your ballot has been cryptographically recorded into the append-only ledger.' 
              : 'You are authenticated and eligible to cast an anonymous ballot in active elections below.'}
          </p>
        </div>

        <div style={{ display: 'flex', gap: '0.8rem', alignItems: 'center' }}>
          <span 
            className={`user-badge ${user.has_voted ? 'badge-officer' : 'badge-voter'}`}
            style={{ fontSize: '0.85rem', padding: '0.4rem 0.9rem' }}
            data-testid="voter-status-badge"
          >
            {user.has_voted ? 'BALLOT CAST (HAS_VOTED = TRUE)' : 'ELIGIBLE TO VOTE'}
          </span>
        </div>
      </div>

      {/* BALLOT RECEIPT CARD */}
      {receipt && (
        <div 
          className="glass-card" 
          style={{ padding: '2rem', border: '1px solid rgba(16, 185, 129, 0.4)', background: 'rgba(16, 185, 129, 0.04)' }}
          data-testid="ballot-receipt-card"
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
            <h3 style={{ color: 'var(--success)', display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '1.2rem' }}>
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              Cryptographic Ballot Receipt Issued
            </h3>
            <button
              onClick={() => copyToClipboard(`Merkle Leaf Hash: ${receipt.receipt_h}\nRoot Hash: ${receipt.root_hash}`)}
              className="btn btn-secondary"
              style={{ padding: '0.35rem 0.8rem', fontSize: '0.8rem' }}
              data-testid="btn-copy-receipt"
            >
              {copiedReceipt ? '✓ Copied' : '📋 Copy Receipt'}
            </button>
          </div>

          <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.2rem' }}>
            Your ballot was encrypted with RSA-2048 and stored as a leaf in the append-only Merkle ledger.
            Your Voter ID is permanently disconnected from this cryptographic record.
          </p>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.8rem', fontFamily: 'monospace', fontSize: '0.85rem' }}>
            <div>
              <span style={{ color: 'var(--text-muted)' }}>MERKLE LEAF HASH (SHA-256):</span>
              <div 
                style={{ background: 'rgba(0,0,0,0.3)', padding: '0.6rem', borderRadius: '6px', marginTop: '0.2rem', overflowX: 'auto', border: '1px solid var(--glass-border)' }}
                data-testid="receipt-leaf-hash"
              >
                {receipt.receipt_h}
              </div>
            </div>
            <div>
              <span style={{ color: 'var(--text-muted)' }}>SYSTEM MERKLE ROOT AT CAST TIME:</span>
              <div 
                style={{ background: 'rgba(0,0,0,0.3)', padding: '0.6rem', borderRadius: '6px', marginTop: '0.2rem', overflowX: 'auto', border: '1px solid var(--glass-border)' }}
                data-testid="receipt-root-hash"
              >
                {receipt.root_hash}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* SCREEN: BALLOT CANDIDATE SELECTION */}
      {activeScreen === 'ballot' && selectedElection ? (
        <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
              <button onClick={onBackToElections} className="btn btn-secondary" style={{ padding: '0.5rem 1rem' }} data-testid="btn-back-to-elections">
                ← Back
              </button>
              <div>
                <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>OFFICIAL BALLOT</span>
                <h2 style={{ fontSize: '1.5rem' }}>{selectedElection.title}</h2>
              </div>
            </div>
          </div>

          <div className="glass-card" style={{ padding: '2rem' }}>
            <p style={{ color: 'var(--text-secondary)', marginBottom: '1.5rem' }}>
              Select your candidate below. When you confirm, the Voting Server encrypts your choice with RSA-2048 and issues your Merkle receipt.
            </p>

            <div className="grid-2">
              {candidates.map(cand => {
                const isSelected = selectedCandidate?.candidate_id === cand.candidate_id;
                return (
                  <div 
                    key={cand.candidate_id} 
                    className="glass-card" 
                    onClick={() => setSelectedCandidate(cand)}
                    style={{ 
                      padding: '1.5rem', 
                      cursor: 'pointer',
                      borderColor: isSelected ? 'var(--primary)' : 'var(--glass-border)',
                      background: isSelected ? 'rgba(139, 92, 246, 0.08)' : 'var(--glass-bg)',
                      boxShadow: isSelected ? '0 0 0 2px var(--primary-glow)' : 'none',
                    }}
                    data-testid={`candidate-card-${cand.candidate_id}`}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                      <span style={{ fontSize: '0.85rem', color: 'var(--secondary)', fontWeight: 700 }}>
                        {cand.party}
                      </span>
                      {isSelected && (
                        <span style={{ color: 'var(--primary)', fontWeight: 'bold', fontSize: '0.9rem' }}>
                          ✓ Selected
                        </span>
                      )}
                    </div>
                    <h4 style={{ fontSize: '1.3rem', marginBottom: '0.3rem' }}>{cand.candidate_name}</h4>
                    <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>ID: {cand.candidate_id.slice(0, 8)}...</p>
                  </div>
                );
              })}
            </div>

            <div style={{ marginTop: '2rem', display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
              <button 
                onClick={onBackToElections} 
                className="btn btn-secondary"
                data-testid="btn-cancel-ballot"
              >
                Cancel
              </button>
              <button 
                disabled={!selectedCandidate || loading} 
                onClick={() => setConfirmModalOpen(true)} 
                className="btn btn-primary"
                data-testid="btn-confirm-selection"
              >
                Proceed to Confirmation →
              </button>
            </div>
          </div>
        </div>
      ) : (
        /* SCREEN: ELECTIONS LIST */
        <div>
          <h3 style={{ marginBottom: '1rem', fontSize: '1.3rem', color: 'var(--text-secondary)' }}>
            Available Elections
          </h3>
          <div className="grid-2">
            {elections.map(elect => (
              <div 
                key={elect.election_id} 
                className="glass-card" 
                style={{ padding: '1.8rem', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}
                data-testid={`election-card-${elect.election_id}`}
              >
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.8rem' }}>
                    <span className={`user-badge ${elect.status === 'Active' ? 'badge-officer' : 'badge-voter'}`}>
                      {elect.status}
                    </span>
                    <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                      Ends: {new Date(elect.end_time).toLocaleDateString()}
                    </span>
                  </div>
                  <h4 style={{ fontSize: '1.3rem', marginBottom: '0.6rem' }}>{elect.title}</h4>
                  <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', marginBottom: '1.5rem', lineHeight: '1.4' }}>
                    {elect.description}
                  </p>
                </div>

                <div style={{ marginTop: 'auto', display: 'flex', gap: '0.6rem' }}>
                  {elect.status === 'Active' && !user.has_voted && (
                    <button 
                      onClick={() => onSelectElectionForVoting(elect)} 
                      className="btn btn-primary" 
                      style={{ flex: 1 }}
                      data-testid={`btn-cast-ballot-${elect.election_id}`}
                    >
                      🗳️ Cast Secret Ballot
                    </button>
                  )}
                  {elect.status === 'Active' && user.has_voted && (
                    <button disabled className="btn btn-secondary" style={{ flex: 1 }} data-testid="btn-already-voted">
                      ✓ Already Voted (Locked)
                    </button>
                  )}
                  {(elect.status === 'Archived' || elect.status === 'Closed') && (
                    <button 
                      onClick={() => onViewResults(elect)} 
                      className="btn btn-secondary" 
                      style={{ flex: 1 }}
                      data-testid={`btn-view-results-${elect.election_id}`}
                    >
                      📊 View Tally Results
                    </button>
                  )}
                  {elect.status === 'Created' && (
                    <button disabled className="btn btn-secondary" style={{ flex: 1 }}>
                      Pending Opening
                    </button>
                  )}
                </div>
              </div>
            ))}

            {elections.length === 0 && (
              <div style={{ color: 'var(--text-muted)', fontStyle: 'italic', padding: '2rem' }}>
                No active elections found.
              </div>
            )}
          </div>
        </div>
      )}

      {/* CONFIRMATION POPUP MODAL */}
      {confirmModalOpen && selectedElection && selectedCandidate && (
        <div 
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'rgba(5, 3, 15, 0.85)',
            backdropFilter: 'blur(16px)',
            zIndex: 1000,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '1.5rem',
          }}
          data-testid="confirm-vote-modal"
        >
          <div className="glass-card" style={{ maxWidth: '550px', width: '100%', padding: '2.5rem', border: '1px solid var(--primary-glow)' }}>
            <h2 style={{ fontSize: '1.6rem', marginBottom: '1.2rem', textAlign: 'center' }}>
              Confirm Your Ballot Submission
            </h2>

            <div style={{ background: 'rgba(0,0,0,0.3)', padding: '1.5rem', borderRadius: '12px', marginBottom: '1.5rem' }}>
              <div style={{ marginBottom: '1rem' }}>
                <span className="form-label" style={{ fontSize: '0.75rem' }}>ELECTION:</span>
                <div style={{ fontSize: '1.1rem', fontWeight: 600, marginTop: '0.2rem' }}>
                  {selectedElection.title}
                </div>
              </div>

              <div>
                <span className="form-label" style={{ fontSize: '0.75rem' }}>CHOSEN CANDIDATE:</span>
                <div style={{ fontSize: '1.4rem', fontWeight: 800, color: 'var(--primary)', marginTop: '0.2rem' }}>
                  {selectedCandidate.candidate_name}
                </div>
                <div style={{ fontSize: '0.9rem', color: 'var(--secondary)', fontWeight: 600 }}>
                  {selectedCandidate.party}
                </div>
              </div>
            </div>

            <div style={{ borderLeft: '3px solid var(--warning)', paddingLeft: '1rem', marginBottom: '1.8rem', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
              <p style={{ fontWeight: 600, color: 'var(--text-primary)', marginBottom: '0.3rem' }}>
                Zero-Knowledge Privacy Guarantee (REQ-AUTH-5):
              </p>
              <p>
                When you click "Confirm & Cast", your vote is encrypted via RSA-2048 before transmission.
                The voting server authenticates your blind token but has zero access to your Voter ID.
              </p>
            </div>

            <div style={{ display: 'flex', gap: '1rem' }}>
              <button 
                onClick={() => setConfirmModalOpen(false)} 
                className="btn btn-secondary" 
                style={{ flex: 1 }}
                data-testid="btn-cancel-confirm"
              >
                Go Back
              </button>
              <button 
                onClick={handleConfirmAndCast} 
                disabled={loading} 
                className="btn btn-primary" 
                style={{ flex: 1 }}
                data-testid="btn-final-cast-vote"
              >
                {loading ? 'Encrypting & Casting...' : 'Confirm & Cast Ballot'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
