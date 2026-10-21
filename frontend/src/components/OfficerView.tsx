import React, { useState } from 'react';
import type { Election, Candidate, MerkleStatus } from '../types';
import { MerkleTreeVisualizer } from './MerkleTreeVisualizer';

interface OfficerViewProps {
  elections: Election[];
  loading: boolean;
  onRefreshElections: () => void;
  onCreateElection: (formData: { title: string; desc: string; start: string; end: string; passphrase: string }) => Promise<void>;
  onOpenElection: (electionId: string) => Promise<void>;
  onCloseElection: (electionId: string) => Promise<void>;
  onVerifyMerkle: (electionId: string) => Promise<void>;
  onAddCandidate: (electionId: string, name: string, party: string) => Promise<void>;
  onLoadCandidates: (electionId: string) => Promise<Candidate[]>;
  merkleStatus: MerkleStatus | null;
  token: string | null;
}

export const OfficerView: React.FC<OfficerViewProps> = ({
  elections,
  loading,
  onRefreshElections,
  onCreateElection,
  onOpenElection,
  onCloseElection,
  onVerifyMerkle,
  onAddCandidate,
  onLoadCandidates,
  merkleStatus,
  token,
}) => {
  // Create Election Form States
  const [newTitle, setNewTitle] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [newStart, setNewStart] = useState('');
  const [newEnd, setNewEnd] = useState('');
  const [newPassphrase, setNewPassphrase] = useState('demo123');

  // Candidate Management State
  const [selectedElection, setSelectedElection] = useState<Election | null>(null);
  const [candidateList, setCandidateList] = useState<Candidate[]>([]);
  const [candidateModalOpen, setCandidateModalOpen] = useState(false);
  const [candName, setCandName] = useState('');
  const [candParty, setCandParty] = useState('');

  // Active Merkle Inspection Election
  const [inspectElectionId, setInspectElectionId] = useState<string | null>(null);

  const fillDemoValues = () => {
    const now = new Date();
    const later = new Date(Date.now() + 24 * 60 * 60 * 1000);
    setNewTitle('Faculty Council Presidential Election 2026');
    setNewDesc('University faculty and student council representative vote.');
    setNewStart(now.toISOString().slice(0, 16));
    setNewEnd(later.toISOString().slice(0, 16));
    setNewPassphrase('demo123');
  };

  const handleCreateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onCreateElection({
      title: newTitle,
      desc: newDesc,
      start: newStart,
      end: newEnd,
      passphrase: newPassphrase,
    });
    setNewTitle('');
    setNewDesc('');
    setNewStart('');
    setNewEnd('');
    setNewPassphrase('demo123');
  };

  const openCandidateModal = async (elect: Election) => {
    setSelectedElection(elect);
    setCandidateModalOpen(true);
    const list = await onLoadCandidates(elect.election_id);
    setCandidateList(list);
  };

  const handleAddCandidateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedElection) return;
    await onAddCandidate(selectedElection.election_id, candName, candParty);
    setCandName('');
    setCandParty('');
    const list = await onLoadCandidates(selectedElection.election_id);
    setCandidateList(list);
  };

  return (
    <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }} data-testid="officer-console-view">
      {/* HEADER BANNER */}
      <div className="glass-card" style={{ padding: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem' }}>
        <div>
          <span style={{ fontSize: '0.8rem', color: 'var(--success)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Officer Console (Lifecycle & Integrity Management)
          </span>
          <h2 style={{ fontSize: '1.7rem', marginTop: '0.2rem' }}>Election Officer Console</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', marginTop: '0.3rem' }}>
            Initialize elections, manage candidates, trigger digital signatures, and verify append-only Merkle ledger integrity.
          </p>
        </div>

        <button 
          onClick={onRefreshElections} 
          disabled={loading}
          className="btn btn-secondary"
          data-testid="btn-refresh-officer-elections"
        >
          🔄 Refresh Elections
        </button>
      </div>

      <div className="grid-2">
        {/* LEFT COLUMN: CREATE NEW ELECTION */}
        <div className="glass-card" style={{ padding: '2rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.2rem', borderBottom: '1px solid var(--glass-border)', paddingBottom: '0.8rem' }}>
            <h3 style={{ fontSize: '1.25rem' }}>Create New Election</h3>
            <button
              type="button"
              onClick={fillDemoValues}
              className="btn btn-secondary"
              style={{ padding: '0.3rem 0.6rem', fontSize: '0.75rem' }}
              data-testid="btn-fill-demo-election"
              title="Auto-fills sample values for quick demonstration"
            >
              ⚡ Fill Demo Data
            </button>
          </div>

          <form onSubmit={handleCreateSubmit} data-testid="form-create-election">
            <div className="form-group">
              <label className="form-label">Election Title</label>
              <input 
                type="text" 
                value={newTitle}
                onChange={(e) => setNewTitle(e.target.value)}
                className="input-field" 
                placeholder="e.g. Student Council Representative 2026"
                required
                data-testid="input-election-title"
              />
            </div>

            <div className="form-group">
              <label className="form-label">Description</label>
              <textarea 
                value={newDesc}
                onChange={(e) => setNewDesc(e.target.value)}
                className="input-field" 
                style={{ minHeight: '75px', resize: 'vertical' }}
                placeholder="Details and scope of this election..."
                required
                data-testid="input-election-desc"
              />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
              <div className="form-group">
                <label className="form-label">Start Date/Time</label>
                <input 
                  type="datetime-local" 
                  value={newStart}
                  onChange={(e) => setNewStart(e.target.value)}
                  className="input-field" 
                  required
                  data-testid="input-election-start"
                />
              </div>

              <div className="form-group">
                <label className="form-label">End Date/Time</label>
                <input 
                  type="datetime-local" 
                  value={newEnd}
                  onChange={(e) => setNewEnd(e.target.value)}
                  className="input-field" 
                  required
                  data-testid="input-election-end"
                />
              </div>
            </div>

            <div className="form-group">
              <label className="form-label">
                Authority Decryption Passphrase (AES-256 Key Derivation)
              </label>
              <input 
                type="password" 
                value={newPassphrase}
                onChange={(e) => setNewPassphrase(e.target.value)}
                className="input-field" 
                placeholder="Passphrase used by Counting Authority to decrypt"
                required
                data-testid="input-election-passphrase"
              />
              <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                Derives AES key to encrypt the RSA-2048 private key. Default: <code>demo123</code>
              </span>
            </div>

            <button 
              type="submit" 
              disabled={loading} 
              className="btn btn-primary" 
              style={{ width: '100%', marginTop: '0.8rem' }}
              data-testid="btn-submit-create-election"
            >
              {loading ? 'Generating RSA-2048 Keypair...' : 'Initialize Election & Generate Keys'}
            </button>
          </form>
        </div>

        {/* RIGHT COLUMN: SYSTEM ELECTIONS MANAGEMENT */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <h3 style={{ fontSize: '1.25rem', color: 'var(--text-secondary)' }}>
            System Elections ({elections.length})
          </h3>

          {elections.map(elect => (
            <div 
              key={elect.election_id} 
              className="glass-card" 
              style={{ padding: '1.5rem' }}
              data-testid={`officer-election-${elect.election_id}`}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                <span className="logo-text" style={{ fontSize: '1.15rem' }}>{elect.title}</span>
                <span className={`user-badge ${elect.status === 'Created' ? 'badge-voter' : elect.status === 'Active' ? 'badge-officer' : 'badge-admin'}`}>
                  {elect.status}
                </span>
              </div>

              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginBottom: '0.8rem', fontFamily: 'monospace' }}>
                ID: {elect.election_id}
              </p>

              {/* ACTION BUTTONS */}
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                {elect.status === 'Created' && (
                  <>
                    <button 
                      onClick={() => openCandidateModal(elect)} 
                      className="btn btn-secondary" 
                      style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                      data-testid={`btn-add-cand-${elect.election_id}`}
                    >
                      👥 Manage Candidates
                    </button>
                    <button 
                      onClick={() => onOpenElection(elect.election_id)} 
                      className="btn btn-primary" 
                      style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                      data-testid={`btn-open-elect-${elect.election_id}`}
                    >
                      ▶ Open Election
                    </button>
                  </>
                )}

                {elect.status === 'Active' && (
                  <>
                    <button 
                      onClick={() => onCloseElection(elect.election_id)} 
                      className="btn btn-danger" 
                      style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                      data-testid={`btn-close-elect-${elect.election_id}`}
                    >
                      🔒 Close & Sign Merkle Root
                    </button>
                    <button 
                      onClick={() => {
                        setInspectElectionId(elect.election_id);
                        onVerifyMerkle(elect.election_id);
                      }} 
                      className="btn btn-secondary" 
                      style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                      data-testid={`btn-verify-merkle-${elect.election_id}`}
                    >
                      🔍 Inspect Merkle Tree
                    </button>
                  </>
                )}

                {(elect.status === 'Closed' || elect.status === 'Archived') && (
                  <button 
                    onClick={() => {
                      setInspectElectionId(elect.election_id);
                      onVerifyMerkle(elect.election_id);
                    }} 
                    className="btn btn-secondary" 
                    style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                    data-testid={`btn-verify-merkle-${elect.election_id}`}
                  >
                    🔍 Verify Signed Merkle Root
                  </button>
                )}
              </div>
            </div>
          ))}

          {elections.length === 0 && (
            <div style={{ color: 'var(--text-muted)', fontStyle: 'italic', padding: '1rem' }}>
              No elections currently configured. Create one on the left or click "Reset Demo".
            </div>
          )}
        </div>
      </div>

      {/* MERKLE TREE VISUALIZER FOR INSPECTED ELECTION */}
      {inspectElectionId && (
        <div style={{ marginTop: '1rem' }}>
          <MerkleTreeVisualizer
            electionId={inspectElectionId}
            status={merkleStatus}
            onRefresh={() => onVerifyMerkle(inspectElectionId)}
            token={token}
          />
        </div>
      )}

      {/* MANAGE CANDIDATES MODAL */}
      {candidateModalOpen && selectedElection && (
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
          data-testid="modal-manage-candidates"
        >
          <div className="glass-card" style={{ maxWidth: '600px', width: '100%', padding: '2.5rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
              <div>
                <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>CANDIDATE REGISTRY</span>
                <h2 style={{ fontSize: '1.4rem' }}>{selectedElection.title}</h2>
              </div>
              <button 
                onClick={() => setCandidateModalOpen(false)} 
                className="btn btn-secondary" 
                style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                data-testid="btn-close-candidate-modal"
              >
                ✕ Close
              </button>
            </div>

            {/* ADD CANDIDATE FORM */}
            <form onSubmit={handleAddCandidateSubmit} style={{ display: 'flex', gap: '0.6rem', marginBottom: '2rem' }}>
              <input 
                type="text" 
                value={candName}
                onChange={(e) => setCandName(e.target.value)}
                className="input-field" 
                placeholder="Candidate Name (e.g. Alice Chen)"
                required
                data-testid="input-candidate-name"
              />
              <input 
                type="text" 
                value={candParty}
                onChange={(e) => setCandParty(e.target.value)}
                className="input-field" 
                placeholder="Party / Affiliation"
                required
                data-testid="input-candidate-party"
              />
              <button 
                type="submit" 
                disabled={loading} 
                className="btn btn-primary" 
                style={{ whiteSpace: 'nowrap' }}
                data-testid="btn-submit-add-candidate"
              >
                + Add
              </button>
            </form>

            {/* CANDIDATE LIST */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.6rem', maxHeight: '300px', overflowY: 'auto' }}>
              {candidateList.map((cand, idx) => (
                <div 
                  key={cand.candidate_id} 
                  style={{ background: 'rgba(0,0,0,0.2)', padding: '0.8rem 1.2rem', borderRadius: '8px', border: '1px solid var(--glass-border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
                  data-testid={`registered-candidate-${idx}`}
                >
                  <span style={{ fontWeight: 600 }}>{cand.candidate_name}</span>
                  <span style={{ color: 'var(--secondary)', fontSize: '0.85rem', fontWeight: 600 }}>{cand.party}</span>
                </div>
              ))}

              {candidateList.length === 0 && (
                <div style={{ textAlign: 'center', color: 'var(--text-muted)', fontStyle: 'italic', padding: '1.5rem' }}>
                  No candidates registered yet. Add at least two candidates before opening.
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
