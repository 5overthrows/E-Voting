import React, { useState } from 'react';
import type { Persona } from '../types';

interface FacultyWalkthroughModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSwitchPersona: (persona: Persona) => void;
  onResetDemo: () => void;
}

interface Step {
  id: number;
  title: string;
  badge: string;
  persona: 'officer1' | 'voter1' | 'voter2' | 'counting1' | 'admin1';
  personaLabel: string;
  talkingPoints: string[];
  actionDescription: string;
  actionButtonLabel: string;
}

const STEPS: Step[] = [
  {
    id: 1,
    title: 'Election Setup & Key Ceremony',
    badge: 'Officer Tier',
    persona: 'officer1',
    personaLabel: 'Election Officer',
    talkingPoints: [
      'The Election Officer initializes an election and generates an RSA-2048 asymmetric keypair.',
      'The private decryption key is encrypted using AES-256 with a PBKDF2/SHA-256 key derived from a secret passphrase.',
      'Neither the Officer nor the server can decrypt ballots before the election closes.'
    ],
    actionDescription: 'Switch to Election Officer to view active elections or create new election.',
    actionButtonLabel: 'Switch to Election Officer',
  },
  {
    id: 2,
    title: 'Voter Authentication & Anonymous Token Issuance',
    badge: 'Auth Tier',
    persona: 'voter1',
    personaLabel: 'Voter 1 (Alice)',
    talkingPoints: [
      'Voter authenticates against the Authentication Server using their credentials.',
      'The Auth Server verifies eligibility (HasVoted == false).',
      'The Token Server generates a single-use cryptographically random token with ZERO Personally Identifiable Information (PII).',
      'The Auth server registers the token hash in-memory but NEVER sees which candidate the voter chooses.'
    ],
    actionDescription: 'Switch to Voter 1 to view ballot and cast anonymous vote.',
    actionButtonLabel: 'Switch to Voter 1',
  },
  {
    id: 3,
    title: 'Ballot Encryption & Merkle Leaf Calculation',
    badge: 'Voting Tier',
    persona: 'voter1',
    personaLabel: 'Voter 1 (Alice)',
    talkingPoints: [
      'The voter casts their vote to the Voting Server (port 8083) using ONLY the anonymous token.',
      'The candidate selection is encrypted using the election RSA-2048 public key.',
      'A SHA-256 leaf hash is computed over the ciphertext and stored in the append-only ledger.',
      'The voter receives a cryptographically verifiable ballot receipt (Merkle Leaf Hash).'
    ],
    actionDescription: 'Review ballot candidate cards and confirm selection.',
    actionButtonLabel: 'Go to Ballot Screen',
  },
  {
    id: 4,
    title: 'Double-Voting Prevention (Sybil Defense)',
    badge: 'Security',
    persona: 'voter1',
    personaLabel: 'Voter 1 (Alice)',
    talkingPoints: [
      'The system immediately marks the voter HasVoted = true and irreversibly invalidates the anonymous token.',
      'Any second token request returns HTTP 403 Forbidden.',
      'Any reuse of the same voting token returns HTTP 401 Unauthorized.',
      'Try clicking "Cast Ballot" again as Voter 1 to demonstrate this live to the faculty!'
    ],
    actionDescription: 'Observe that Voter 1 is marked as "Already Voted" and blocked from voting again.',
    actionButtonLabel: 'Verify Voter Status',
  },
  {
    id: 5,
    title: 'Second Voter (Bob) Casts Vote & Tree Grows',
    badge: 'Voter Tier',
    persona: 'voter2',
    personaLabel: 'Voter 2 (Bob)',
    talkingPoints: [
      'A second voter (Bob) casts their vote.',
      'The Merkle Tree recalculates the root hash by combining Leaf 1 and Leaf 2: Root = SHA256(H1 + H2).',
      'Demonstrates that the tree grows incrementally without leaking voter identities.'
    ],
    actionDescription: 'Switch to Voter 2 to cast a second vote and watch the Merkle Tree update.',
    actionButtonLabel: 'Switch to Voter 2',
  },
  {
    id: 6,
    title: 'Election Closure & Digital Signature of Merkle Root',
    badge: 'Officer Tier',
    persona: 'officer1',
    personaLabel: 'Election Officer',
    talkingPoints: [
      'When the voting window expires, the Election Officer closes the election.',
      'The system freezes the ledger, calculates the final Merkle Root, and digitally signs it with the system private key.',
      'No more votes can be appended after this step.'
    ],
    actionDescription: 'Switch to Officer to close election and sign the root.',
    actionButtonLabel: 'Switch to Election Officer',
  },
  {
    id: 7,
    title: 'Counting Authority Decryption & Results Publication',
    badge: 'Authority Tier',
    persona: 'counting1',
    personaLabel: 'Counting Authority',
    talkingPoints: [
      'The Counting Authority inputs their secret passphrase to derive the AES key.',
      'The RSA private key is decrypted in memory and used to decrypt only the closed ballots.',
      'Votes are tallied and published with the archived signed Merkle root for public verification.'
    ],
    actionDescription: 'Switch to Counting Authority to enter passphrase (demo123) and decrypt the tally.',
    actionButtonLabel: 'Switch to Counting Authority',
  },
  {
    id: 8,
    title: 'Live Mathematical Tamper Simulation & Immutable Audit',
    badge: 'Auditor Tier',
    persona: 'admin1',
    personaLabel: 'System Administrator',
    talkingPoints: [
      'Demonstrate the ultimate defense: Click "Simulate DB Tamper" in the Merkle Tree Visualizer!',
      'Show that when an attacker modifies a vote directly in SQLite/Postgres, the Merkle verification instantly fails (turns RED).',
      'Click "Restore Ledger" and show the cryptographic integrity return to 100% verified (turns GREEN).',
      'Review immutable audit events and trigger a database backup.'
    ],
    actionDescription: 'Switch to Administrator to view immutable audit logs and backup database.',
    actionButtonLabel: 'Switch to Administrator',
  },
];

export const FacultyWalkthroughModal: React.FC<FacultyWalkthroughModalProps> = ({
  isOpen,
  onClose,
  onSwitchPersona,
  onResetDemo,
}) => {
  const [activeStepIndex, setActiveStepIndex] = useState(0);

  if (!isOpen) return null;

  const currentStep = STEPS[activeStepIndex];

  return (
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
      data-testid="faculty-walkthrough-modal"
    >
      <div 
        className="glass-card" 
        style={{
          maxWidth: '900px',
          width: '100%',
          maxHeight: '90vh',
          display: 'flex',
          flexDirection: 'column',
          padding: '2rem',
          border: '1px solid var(--primary-glow)',
        }}
      >
        {/* HEADER */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', borderBottom: '1px solid var(--glass-border)', paddingBottom: '1rem' }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
              <span style={{ fontSize: '1.5rem' }}>🎓</span>
              <h2 style={{ fontSize: '1.4rem' }}>Faculty Presentation & Defense Storyboard</h2>
              <span className="user-badge" style={{ background: 'rgba(16, 185, 129, 0.2)', color: 'var(--success)' }}>
                Step {currentStep.id} of {STEPS.length}
              </span>
            </div>
            <p style={{ color: 'var(--text-secondary)', fontSize: '0.85rem', marginTop: '0.2rem' }}>
              Structured, error-free demonstration guide for presenting to professors and review panel.
            </p>
          </div>

          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button 
              onClick={() => {
                onResetDemo();
                setActiveStepIndex(0);
              }} 
              className="btn btn-secondary" 
              style={{ padding: '0.4rem 0.8rem', fontSize: '0.8rem' }}
              title="Resets demo database to clean initial state"
              data-testid="btn-reset-demo-walkthrough"
            >
              ⚡ Reset Demo State
            </button>
            <button 
              onClick={onClose} 
              className="btn btn-secondary" 
              style={{ padding: '0.4rem 0.8rem', fontSize: '0.8rem' }}
              data-testid="btn-close-walkthrough"
            >
              ✕ Close
            </button>
          </div>
        </div>

        {/* STEP PROGRESS TRACKER */}
        <div style={{ display: 'flex', gap: '0.3rem', marginBottom: '1.5rem', overflowX: 'auto', paddingBottom: '0.5rem' }}>
          {STEPS.map((step, idx) => (
            <button
              key={step.id}
              onClick={() => setActiveStepIndex(idx)}
              style={{
                background: idx === activeStepIndex ? 'var(--primary)' : 'rgba(255,255,255,0.04)',
                border: `1px solid ${idx === activeStepIndex ? 'var(--primary)' : 'var(--glass-border)'}`,
                color: idx === activeStepIndex ? '#fff' : 'var(--text-muted)',
                padding: '0.4rem 0.8rem',
                borderRadius: '6px',
                fontSize: '0.75rem',
                fontWeight: 600,
                cursor: 'pointer',
                whiteSpace: 'nowrap',
                transition: 'all 0.2s ease',
              }}
              data-testid={`step-tab-${step.id}`}
            >
              {step.id}. {step.badge}
            </button>
          ))}
        </div>

        {/* STEP CONTENT BODY */}
        <div style={{ flex: 1, overflowY: 'auto', background: 'rgba(0,0,0,0.2)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--glass-border)', marginBottom: '1.5rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.8rem' }}>
            <span className="user-badge" style={{ background: 'rgba(139, 92, 246, 0.2)', color: 'var(--primary)', fontSize: '0.75rem' }}>
              {currentStep.badge}
            </span>
            <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
              Target Persona: <strong style={{ color: 'var(--text-primary)' }}>{currentStep.personaLabel}</strong>
            </span>
          </div>

          <h3 style={{ fontSize: '1.3rem', marginBottom: '1rem', color: 'var(--text-primary)' }}>
            {currentStep.title}
          </h3>

          <div style={{ marginBottom: '1.2rem' }}>
            <div style={{ fontSize: '0.8rem', fontWeight: 700, color: 'var(--secondary)', textTransform: 'uppercase', marginBottom: '0.5rem' }}>
              Faculty Viva Talking Points:
            </div>
            <ul style={{ paddingLeft: '1.2rem', display: 'flex', flexDirection: 'column', gap: '0.5rem', fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
              {currentStep.talkingPoints.map((pt, pIdx) => (
                <li key={pIdx}>{pt}</li>
              ))}
            </ul>
          </div>

          <div style={{ borderTop: '1px solid var(--glass-border)', paddingTop: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '0.8rem' }}>
            <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
              {currentStep.actionDescription}
            </div>
            <button
              onClick={() => {
                onSwitchPersona(currentStep.persona);
                onClose();
              }}
              className="btn btn-primary"
              style={{ padding: '0.5rem 1.2rem', fontSize: '0.85rem' }}
              data-testid="btn-walkthrough-execute-action"
            >
              ▶ {currentStep.actionButtonLabel}
            </button>
          </div>
        </div>

        {/* FOOTER NAVIGATION */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <button
            onClick={() => setActiveStepIndex(Math.max(0, activeStepIndex - 1))}
            disabled={activeStepIndex === 0}
            className="btn btn-secondary"
            style={{ padding: '0.5rem 1.2rem', fontSize: '0.85rem' }}
            data-testid="btn-walkthrough-prev"
          >
            ← Previous Step
          </button>

          <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
            Step {activeStepIndex + 1} of {STEPS.length}
          </span>

          <button
            onClick={() => setActiveStepIndex(Math.min(STEPS.length - 1, activeStepIndex + 1))}
            disabled={activeStepIndex === STEPS.length - 1}
            className="btn btn-primary"
            style={{ padding: '0.5rem 1.2rem', fontSize: '0.85rem' }}
            data-testid="btn-walkthrough-next"
          >
            Next Step →
          </button>
        </div>
      </div>
    </div>
  );
};
