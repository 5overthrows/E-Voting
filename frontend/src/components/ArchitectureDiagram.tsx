import React from 'react';

export const ArchitectureDiagram: React.FC = () => {
  return (
    <div 
      className="glass-card" 
      style={{ padding: '1.5rem', marginBottom: '1.5rem', borderLeft: '4px solid var(--secondary)' }}
      data-testid="architecture-privacy-wall"
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <div>
          <span style={{ fontSize: '0.75rem', color: 'var(--secondary)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            System Architecture & Cryptographic Privacy Wall (IEEE SRS Model)
          </span>
          <h3 style={{ fontSize: '1.1rem', marginTop: '0.2rem' }}>
            Zero-Knowledge Voter Separation & Append-Only Ledger
          </h3>
        </div>
        <span className="user-badge" style={{ background: 'rgba(6, 182, 212, 0.1)', color: 'var(--secondary)' }}>
          REQ-DB-1: No UserID in Votes DB
        </span>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr auto 1fr auto 1fr', gap: '0.75rem', alignItems: 'center', fontSize: '0.85rem' }}>
        {/* Tier 1: Voter Identity */}
        <div style={{ background: 'rgba(139, 92, 246, 0.08)', border: '1px solid rgba(139, 92, 246, 0.25)', borderRadius: '10px', padding: '1rem' }}>
          <div style={{ fontWeight: 700, color: 'var(--primary)', marginBottom: '0.3rem', display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <span>👤</span> 1. Authentication Tier
          </div>
          <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <div>• Verifies Voter ID + Password</div>
            <div>• Validates <code>HasVoted == false</code></div>
            <div>• Issues single-use Blind Token</div>
            <div style={{ color: 'var(--text-muted)', marginTop: '0.4rem', fontSize: '0.75rem' }}>
              ⚠️ Knows WHO you are, but NEVER sees WHAT you vote.
            </div>
          </div>
        </div>

        {/* Separator / Privacy Barrier */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', color: 'var(--warning)' }}>
          <span style={{ fontSize: '1.2rem' }}>⚡</span>
          <div style={{ height: '35px', width: '2px', background: 'var(--warning)', margin: '0.2rem 0' }}></div>
          <span style={{ fontSize: '0.65rem', fontWeight: 800, textTransform: 'uppercase', background: 'rgba(245, 158, 11, 0.15)', padding: '2px 4px', borderRadius: '4px' }}>
            Privacy Wall
          </span>
        </div>

        {/* Tier 2: Anonymous Voting */}
        <div style={{ background: 'rgba(6, 182, 212, 0.08)', border: '1px solid rgba(6, 182, 212, 0.25)', borderRadius: '10px', padding: '1rem' }}>
          <div style={{ fontWeight: 700, color: 'var(--secondary)', marginBottom: '0.3rem', display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <span>🗳️</span> 2. Anonymous Voting Tier
          </div>
          <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <div>• Validates anonymous token hash</div>
            <div>• Encrypts ballot using RSA-2048</div>
            <div>• Immediately invalidates token</div>
            <div style={{ color: 'var(--text-muted)', marginTop: '0.4rem', fontSize: '0.75rem' }}>
              ⚠️ Sees WHAT is voted, but has ZERO knowledge of WHO voted.
            </div>
          </div>
        </div>

        {/* Arrow to Ledger */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', color: 'var(--success)' }}>
          <span style={{ fontSize: '1.2rem' }}>➔</span>
          <span style={{ fontSize: '0.65rem', fontWeight: 800, textTransform: 'uppercase', color: 'var(--success)' }}>
            Append Only
          </span>
        </div>

        {/* Tier 3: Merkle Tree Ledger */}
        <div style={{ background: 'rgba(16, 185, 129, 0.08)', border: '1px solid rgba(16, 185, 129, 0.25)', borderRadius: '10px', padding: '1rem' }}>
          <div style={{ fontWeight: 700, color: 'var(--success)', marginBottom: '0.3rem', display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
            <span>🌲</span> 3. Merkle Tree Ledger
          </div>
          <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <div>• Leaf = SHA256(EncryptedBallot)</div>
            <div>• Signed Root archived on close</div>
            <div>• Mathematically tamper-evident</div>
            <div style={{ color: 'var(--text-muted)', marginTop: '0.4rem', fontSize: '0.75rem' }}>
              🔒 Any modification to DB invalidates Merkle proof.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
