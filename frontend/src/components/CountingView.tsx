import React, { useState } from 'react';
import type { Election } from '../types';

interface CountingViewProps {
  elections: Election[];
  loading: boolean;
  onRefreshElections: () => void;
  onTallyVotes: (electionId: string, passphrase: string) => Promise<void>;
  talliedResults: any[];
  onFetchResults: (electionId: string) => Promise<void>;
}

export const CountingView: React.FC<CountingViewProps> = ({
  elections,
  loading,
  onRefreshElections,
  onTallyVotes,
  talliedResults,
  onFetchResults,
}) => {
  const [selectedElectionId, setSelectedElectionId] = useState('');
  const [passphrase, setPassphrase] = useState('demo123');

  const closedElections = elections.filter(e => e.status === 'Closed' || e.status === 'Archived');
  const selectedElection = elections.find(e => e.election_id === selectedElectionId);

  const handleTallySubmit = async () => {
    if (!selectedElectionId || !passphrase) return;
    await onTallyVotes(selectedElectionId, passphrase);
  };

  const totalVotes = (talliedResults || []).reduce((acc: number, r: any) => acc + (r.votes || 0), 0);

  return (
    <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }} data-testid="counting-console-view">
      {/* HEADER BANNER */}
      <div className="glass-card" style={{ padding: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem' }}>
        <div>
          <span style={{ fontSize: '0.8rem', color: 'var(--accent)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Counting Authority Portal (Decryption & Tally)
          </span>
          <h2 style={{ fontSize: '1.7rem', marginTop: '0.2rem' }}>Counting Authority Console</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', marginTop: '0.3rem' }}>
            Decrypt closed election ballots, verify Merkle Tree signatures, and publish final authenticated tallies.
          </p>
        </div>

        <button 
          onClick={onRefreshElections} 
          disabled={loading}
          className="btn btn-secondary"
          data-testid="btn-refresh-counting-elections"
        >
          🔄 Refresh Elections
        </button>
      </div>

      <div className="grid-2">
        {/* LEFT COLUMN: SELECTION & PASSPHRASE ENTRY */}
        <div className="glass-card" style={{ padding: '2rem' }}>
          <h3 style={{ fontSize: '1.25rem', marginBottom: '1.5rem', borderBottom: '1px solid var(--glass-border)', paddingBottom: '0.6rem' }}>
            Decrypt Ledger & Run Tally
          </h3>

          <div className="form-group">
            <label className="form-label">Select Closed Election</label>
            <select 
              value={selectedElectionId}
              onChange={(e) => {
                setSelectedElectionId(e.target.value);
                if (e.target.value) {
                  onFetchResults(e.target.value);
                }
              }}
              className="input-field"
              style={{ background: '#0a0815' }}
              data-testid="select-closed-election"
            >
              <option value="">-- Choose Election to Tally --</option>
              {closedElections.map(el => (
                <option key={el.election_id} value={el.election_id}>
                  {el.title} ({el.status})
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Authority Private Key Passphrase</label>
            <input 
              type="password" 
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
              className="input-field" 
              placeholder="Passphrase to derive AES key and decrypt RSA private key"
              data-testid="input-authority-passphrase"
            />
            <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
              Uses PBKDF2/SHA-256 + Salt to derive AES-256 key. Default demo passphrase: <code>demo123</code>
            </span>
          </div>

          <button 
            onClick={handleTallySubmit} 
            disabled={!selectedElectionId || !passphrase || loading} 
            className="btn btn-primary" 
            style={{ width: '100%', marginTop: '1rem' }}
            data-testid="btn-submit-tally"
          >
            {loading ? 'Decrypting Ballots & Verifying...' : '🔓 Decrypt & Compute Final Tally'}
          </button>

          {/* PRIVACY EXPLANATION */}
          <div style={{ marginTop: '1.5rem', padding: '1rem', background: 'rgba(0,0,0,0.2)', borderRadius: '8px', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Cryptographic Guarantee (REQ-ENC-3):</strong>
            <p style={{ marginTop: '0.3rem' }}>
              Ballots cannot be decrypted one-by-one during the election. Decryption occurs only in memory after the election closes and requires the authority's secret passphrase.
            </p>
          </div>
        </div>

        {/* RIGHT COLUMN: RESULTS VISUALIZATION */}
        <div className="glass-card" style={{ padding: '2rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.2rem', borderBottom: '1px solid var(--glass-border)', paddingBottom: '0.6rem' }}>
            <h3 style={{ fontSize: '1.25rem' }}>Election Results Tally</h3>
            {talliedResults.length > 0 && (
              <span className="user-badge status-success" data-testid="tally-count-badge">
                {totalVotes} Total Ballots Counted
              </span>
            )}
          </div>

          {selectedElection?.signed_merkle_root && (
            <div style={{ marginBottom: '1.2rem', fontSize: '0.8rem' }}>
              <span style={{ color: 'var(--text-muted)', fontWeight: 600 }}>ARCHIVED SIGNED MERKLE ROOT:</span>
              <div style={{ background: 'rgba(0,0,0,0.3)', padding: '0.5rem', borderRadius: '4px', fontFamily: 'monospace', overflowX: 'auto', marginTop: '0.2rem' }}>
                {selectedElection.signed_merkle_root}
              </div>
            </div>
          )}

          {talliedResults.length > 0 ? (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              {talliedResults.map((res: any, idx: number) => {
                const percentage = totalVotes > 0 ? Math.round((res.votes / totalVotes) * 100) : 0;
                return (
                  <div 
                    key={res.candidate_id || idx} 
                    style={{ background: 'rgba(255,255,255,0.02)', padding: '1rem', borderRadius: '8px', border: '1px solid var(--glass-border)' }}
                    data-testid={`tally-result-row-${idx}`}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                      <div>
                        <span style={{ fontWeight: 700, fontSize: '1.1rem' }}>{res.candidate_name}</span>
                        <span style={{ color: 'var(--secondary)', fontSize: '0.85rem', marginLeft: '0.6rem' }}>({res.party})</span>
                      </div>
                      <div style={{ fontWeight: 800, fontSize: '1.2rem', color: 'var(--primary)' }}>
                        {res.votes} <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>({percentage}%)</span>
                      </div>
                    </div>

                    {/* PROGRESS BAR */}
                    <div style={{ width: '100%', height: '8px', background: 'rgba(0,0,0,0.3)', borderRadius: '9999px', overflow: 'hidden' }}>
                      <div 
                        style={{ 
                          width: `${percentage}%`, 
                          height: '100%', 
                          background: 'linear-gradient(90deg, var(--primary), var(--secondary))',
                          transition: 'width 0.5s ease',
                        }} 
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div style={{ textAlign: 'center', padding: '3rem 1rem', color: 'var(--text-muted)', fontStyle: 'italic' }}>
              Select a closed election and click "Decrypt & Compute Final Tally" to view verified results.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
