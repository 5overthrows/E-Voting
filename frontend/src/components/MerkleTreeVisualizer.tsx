import React, { useState } from 'react';
import type { MerkleStatus } from '../types';

interface MerkleTreeVisualizerProps {
  electionId: string;
  status: MerkleStatus | null;
  onRefresh: () => void | Promise<void>;
  token: string | null;
}

export const MerkleTreeVisualizer: React.FC<MerkleTreeVisualizerProps> = ({
  electionId,
  status,
  onRefresh,
}) => {
  const [tampering, setTampering] = useState(false);
  const [verifying, setVerifying] = useState(false);
  const [tamperActionMsg, setTamperActionMsg] = useState<string | null>(null);
  const [selectedLeaf, setSelectedLeaf] = useState<string | null>(null);

  const leaves = status?.leaves || [];

  const handleReVerify = async () => {
    setVerifying(true);
    try {
      await onRefresh();
    } finally {
      setVerifying(false);
    }
  };

  const handleSimulateTamper = async () => {
    setTampering(true);
    setTamperActionMsg(null);
    try {
      const resp = await fetch('/api/demo/tamper', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ election_id: electionId }),
      });
      const data = await resp.json();
      if (resp.ok) {
        setTamperActionMsg('⚠️ Malicious modification injected into Votes table! Re-verify to see cryptographic detection.');
        await onRefresh();
      } else {
        setTamperActionMsg(data.error || 'Failed to simulate tamper.');
      }
    } catch {
      setTamperActionMsg('Server connection error.');
    } finally {
      setTampering(false);
    }
  };

  const handleRestore = async () => {
    setTampering(true);
    setTamperActionMsg(null);
    try {
      const resp = await fetch('/api/demo/restore', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ election_id: electionId }),
      });
      const data = await resp.json();
      if (resp.ok) {
        setTamperActionMsg('✅ Database restored to authentic state! Merkle tree integrity is back to 100%.');
        await onRefresh();
      } else {
        setTamperActionMsg(data.error || 'Failed to restore ledger.');
      }
    } catch {
      setTamperActionMsg('Server connection error.');
    } finally {
      setTampering(false);
    }
  };

  const isTampered = status?.status === 'tampered' || status?.verification_pass === false;

  return (
    <div 
      className="glass-card" 
      style={{ 
        padding: '1.8rem', 
        borderColor: isTampered ? 'rgba(239, 68, 68, 0.4)' : 'var(--glass-border)',
        background: isTampered ? 'rgba(239, 68, 68, 0.04)' : 'var(--glass-bg)'
      }}
      data-testid="merkle-tree-visualizer"
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', flexWrap: 'wrap', gap: '1rem', marginBottom: '1.2rem' }}>
        <div>
          <span style={{ fontSize: '0.75rem', color: isTampered ? 'var(--error)' : 'var(--success)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Cryptographic Integrity Engine (REQ-MRK-1 to REQ-MRK-4)
          </span>
          <h3 style={{ fontSize: '1.3rem', marginTop: '0.2rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <span>🌲</span> Merkle Tree Ledger Visualizer
          </h3>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
            Every encrypted ballot forms a tamper-evident cryptographic leaf hash.
          </p>
        </div>

        {/* Status Pill */}
        <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
          {status && (
            <span 
              className={`user-badge ${isTampered ? 'status-error' : 'status-success'}`}
              style={{ fontSize: '0.85rem', padding: '0.4rem 0.9rem' }}
              data-testid="merkle-verification-badge"
            >
              {isTampered ? '🚨 TAMPER DETECTED / MISMATCH' : '✓ MERKLE INTEGRITY VERIFIED'}
            </span>
          )}
          <button 
            onClick={handleReVerify} 
            disabled={verifying}
            className="btn btn-secondary" 
            style={{ padding: '0.4rem 0.8rem', fontSize: '0.8rem', opacity: verifying ? 0.7 : 1 }}
            data-testid="btn-refresh-merkle"
          >
            {verifying ? '🔄 Verifying...' : 'Re-Verify'}
          </button>
        </div>
      </div>

      {/* Action Notification */}
      {tamperActionMsg && (
        <div className={`status-banner ${tamperActionMsg.includes('⚠️') ? 'status-error' : 'status-success'}`} style={{ fontSize: '0.85rem' }}>
          {tamperActionMsg}
        </div>
      )}

      {/* ROOT NODE CARD */}
      <div 
        style={{ 
          background: isTampered ? 'rgba(239, 68, 68, 0.15)' : 'rgba(16, 185, 129, 0.1)',
          border: `1px solid ${isTampered ? 'rgba(239, 68, 68, 0.4)' : 'rgba(16, 185, 129, 0.3)'}`,
          borderRadius: '12px',
          padding: '1.2rem',
          textAlign: 'center',
          marginBottom: '1.5rem'
        }}
        data-testid="merkle-root-node"
      >
        <div style={{ fontSize: '0.75rem', fontWeight: 700, textTransform: 'uppercase', color: isTampered ? 'var(--error)' : 'var(--success)' }}>
          System Merkle Root (Level 0)
        </div>
        <div style={{ fontFamily: 'monospace', fontSize: '0.95rem', fontWeight: 600, wordBreak: 'break-all', margin: '0.5rem 0' }}>
          {status?.computed_root || 'Computing Merkle Root...'}
        </div>
        <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
          Ballots in Ledger: <strong>{status?.ballots_count ?? leaves.length}</strong> | 
          Status: <strong style={{ color: isTampered ? 'var(--error)' : 'var(--success)' }}>{status?.status?.toUpperCase()}</strong>
        </div>
      </div>

      {/* TREE CONNECTOR LINE */}
      {leaves.length > 0 && (
        <div style={{ textAlign: 'center', margin: '-1rem 0 1rem 0', color: 'var(--text-muted)', fontSize: '1.2rem' }}>
          │<br />▼
        </div>
      )}

      {/* LEAF NODES (BALLOTS) */}
      <div>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.6rem' }}>
          <span style={{ fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-secondary)', textTransform: 'uppercase' }}>
            Ballot Leaves: H(EncryptedBallot) ({leaves.length} Total)
          </span>
          <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Click a leaf to inspect SHA-256 payload</span>
        </div>

        {leaves.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-muted)', fontStyle: 'italic', background: 'rgba(0,0,0,0.1)', borderRadius: '8px' }}>
            No votes have been added to this ledger yet. Cast a ballot to see the tree grow!
          </div>
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '0.8rem' }}>
            {leaves.map((leaf, idx) => {
              const isSelected = selectedLeaf === leaf;
              const isLeafTampered = leaf.startsWith('0000tampered');
              return (
                <div
                  key={idx}
                  onClick={() => setSelectedLeaf(isSelected ? null : leaf)}
                  style={{
                    background: isLeafTampered 
                      ? 'rgba(239, 68, 68, 0.25)' 
                      : isSelected ? 'rgba(139, 92, 246, 0.15)' : 'rgba(0,0,0,0.2)',
                    border: `1px solid ${isLeafTampered ? 'var(--error)' : isSelected ? 'var(--primary)' : 'var(--glass-border)'}`,
                    borderRadius: '8px',
                    padding: '0.8rem',
                    cursor: 'pointer',
                    transition: 'all 0.2s ease',
                  }}
                  data-testid={`merkle-leaf-${idx}`}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.3rem' }}>
                    <span style={{ fontSize: '0.75rem', fontWeight: 700, color: isLeafTampered ? 'var(--error)' : 'var(--secondary)' }}>
                      Leaf #{idx + 1}
                    </span>
                    <span style={{ 
                      fontSize: '0.65rem', 
                      background: isLeafTampered ? 'var(--error)' : 'rgba(255,255,255,0.05)', 
                      color: isLeafTampered ? '#ffffff' : 'inherit',
                      padding: '2px 5px', 
                      borderRadius: '4px',
                      fontWeight: isLeafTampered ? 700 : 400
                    }}>
                      {isLeafTampered ? 'TAMPERED' : 'SHA-256'}
                    </span>
                  </div>
                  <div style={{ fontFamily: 'monospace', fontSize: '0.75rem', color: isLeafTampered ? '#fca5a5' : 'var(--text-primary)', wordBreak: 'break-all' }}>
                    {leaf.slice(0, 16)}...{leaf.slice(-8)}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* SELECTED LEAF DETAILS */}
        {selectedLeaf && (
          <div className="glass-card" style={{ padding: '1rem', marginTop: '1rem', background: 'rgba(0,0,0,0.3)', fontSize: '0.85rem' }}>
            <div style={{ fontWeight: 600, color: 'var(--primary)', marginBottom: '0.3rem' }}>Selected Merkle Leaf Hash (Full SHA-256):</div>
            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all', background: '#0a0815', padding: '0.5rem', borderRadius: '6px' }}>
              {selectedLeaf}
            </div>
            <div style={{ marginTop: '0.4rem', color: 'var(--text-muted)', fontSize: '0.75rem' }}>
              This leaf is cryptographically bound into the parent branches and signed by the officer upon closing.
            </div>
          </div>
        )}
      </div>

      {/* FACULTY INTERACTIVE TAMPER EXPERIMENT BUTTONS */}
      <div 
        style={{ 
          marginTop: '1.5rem', 
          paddingTop: '1.2rem', 
          borderTop: '1px solid var(--glass-border)',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: '1rem'
        }}
      >
        <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
          <strong>Faculty Demo Controls:</strong> Demonstrate mathematical tamper-evidence to the review committee.
        </div>
        <div style={{ display: 'flex', gap: '0.6rem' }}>
          <button
            onClick={handleSimulateTamper}
            disabled={tampering || leaves.length === 0}
            className="btn btn-danger"
            style={{ padding: '0.5rem 1rem', fontSize: '0.85rem' }}
            data-testid="btn-simulate-tamper"
            title="Simulates modifying a stored ballot inside the database"
          >
            {tampering ? 'Simulating...' : '💥 Simulate DB Tamper'}
          </button>
          <button
            onClick={handleRestore}
            disabled={tampering}
            className="btn btn-secondary"
            style={{ padding: '0.5rem 1rem', fontSize: '0.85rem' }}
            data-testid="btn-restore-ledger"
            title="Restores the original untampered ballot payload"
          >
            🔄 Restore Ledger
          </button>
        </div>
      </div>
    </div>
  );
};
