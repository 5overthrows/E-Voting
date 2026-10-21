import React from 'react';
import type { User } from '../types';

interface HeaderProps {
  user: User | null;
  onSwitchPersona: (persona: 'voter1' | 'voter2' | 'officer1' | 'counting1' | 'admin1') => void;
  onOpenWalkthrough: () => void;
  onOpenJira: () => void;
  onResetDemo: () => void;
  onOpenManualLogin: () => void;
  onLogout: () => void;
  loading: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  user,
  onSwitchPersona,
  onOpenWalkthrough,
  onOpenJira,
  onResetDemo,
  onOpenManualLogin,
  onLogout,
  loading,
}) => {
  return (
    <header className="app-header" data-testid="app-header">
      {/* BRAND & BADGE */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '1.2rem' }}>
        <div className="logo-container" data-testid="logo-brand">
          <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="#8b5cf6" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="#06b6d4" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="#8b5cf6" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
          <div>
            <span className="logo-text">Aegis<span className="text-gradient-purple">Vote</span></span>
            <span style={{ display: 'block', fontSize: '0.65rem', color: 'var(--text-muted)', letterSpacing: '0.04em', fontWeight: 600 }}>
              IEEE SRS PRIVACY-PRESERVING LEDGER
            </span>
          </div>
        </div>

        {/* FACULTY / JIRA TOOLBAR BUTTONS */}
        <div style={{ display: 'flex', gap: '0.5rem', marginLeft: '0.5rem' }}>
          <button
            onClick={onOpenWalkthrough}
            className="btn btn-secondary"
            style={{ padding: '0.4rem 0.8rem', fontSize: '0.8rem', borderColor: 'rgba(139, 92, 246, 0.4)', background: 'rgba(139, 92, 246, 0.08)' }}
            data-testid="btn-open-walkthrough"
            title="Step-by-step presentation storyboard for faculty review"
          >
            🎓 Faculty Walkthrough
          </button>

          <button
            onClick={onOpenJira}
            className="btn btn-secondary"
            style={{ padding: '0.4rem 0.8rem', fontSize: '0.8rem', borderColor: 'rgba(6, 182, 212, 0.4)', background: 'rgba(6, 182, 212, 0.08)' }}
            data-testid="btn-open-jira"
            title="Live test matrix and Jira / Xray test runner"
          >
            🧪 Jira QA Suite
          </button>

          <button
            onClick={onResetDemo}
            disabled={loading}
            className="btn btn-secondary"
            style={{ padding: '0.4rem 0.7rem', fontSize: '0.8rem' }}
            data-testid="btn-reset-demo"
            title="Clean database reset for a fresh demonstration"
          >
            ⚡ Reset Demo
          </button>
        </div>
      </div>

      {/* PERSONA QUICK SWITCHER BAR */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.8rem' }}>
        <div 
          style={{
            display: 'flex',
            alignItems: 'center',
            background: 'rgba(0,0,0,0.35)',
            padding: '3px',
            borderRadius: '10px',
            border: '1px solid var(--glass-border)',
            gap: '3px',
          }}
          data-testid="persona-switcher-bar"
        >
          <span style={{ fontSize: '0.7rem', color: 'var(--text-muted)', padding: '0 0.5rem', fontWeight: 600 }}>
            PERSONA:
          </span>

          <button
            onClick={() => onSwitchPersona('voter1')}
            disabled={loading}
            style={{
              background: user?.user_id === 'voter1' ? 'rgba(6, 182, 212, 0.25)' : 'transparent',
              border: `1px solid ${user?.user_id === 'voter1' ? 'var(--secondary)' : 'transparent'}`,
              color: user?.user_id === 'voter1' ? '#22d3ee' : 'var(--text-secondary)',
              padding: '0.3rem 0.6rem',
              borderRadius: '6px',
              fontSize: '0.75rem',
              fontWeight: 600,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.3rem',
            }}
            data-testid="persona-switch-voter1"
            title="Alice - Eligible Voter"
          >
            <span>👤</span> Voter 1
          </button>

          <button
            onClick={() => onSwitchPersona('voter2')}
            disabled={loading}
            style={{
              background: user?.user_id === 'voter2' ? 'rgba(6, 182, 212, 0.25)' : 'transparent',
              border: `1px solid ${user?.user_id === 'voter2' ? 'var(--secondary)' : 'transparent'}`,
              color: user?.user_id === 'voter2' ? '#22d3ee' : 'var(--text-secondary)',
              padding: '0.3rem 0.6rem',
              borderRadius: '6px',
              fontSize: '0.75rem',
              fontWeight: 600,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.3rem',
            }}
            data-testid="persona-switch-voter2"
            title="Bob - Eligible Voter"
          >
            <span>👤</span> Voter 2
          </button>

          <button
            onClick={() => onSwitchPersona('officer1')}
            disabled={loading}
            style={{
              background: user?.role === 'ElectionOfficer' ? 'rgba(16, 185, 129, 0.25)' : 'transparent',
              border: `1px solid ${user?.role === 'ElectionOfficer' ? 'var(--success)' : 'transparent'}`,
              color: user?.role === 'ElectionOfficer' ? '#34d399' : 'var(--text-secondary)',
              padding: '0.3rem 0.6rem',
              borderRadius: '6px',
              fontSize: '0.75rem',
              fontWeight: 600,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.3rem',
            }}
            data-testid="persona-switch-officer1"
            title="Election Officer - Creates & Closes Elections"
          >
            <span>🛡️</span> Officer
          </button>

          <button
            onClick={() => onSwitchPersona('counting1')}
            disabled={loading}
            style={{
              background: user?.role === 'CountingAuthority' ? 'rgba(236, 72, 153, 0.25)' : 'transparent',
              border: `1px solid ${user?.role === 'CountingAuthority' ? 'var(--accent)' : 'transparent'}`,
              color: user?.role === 'CountingAuthority' ? '#f472b6' : 'var(--text-secondary)',
              padding: '0.3rem 0.6rem',
              borderRadius: '6px',
              fontSize: '0.75rem',
              fontWeight: 600,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.3rem',
            }}
            data-testid="persona-switch-counting1"
            title="Counting Authority - Holds Decryption Passphrase"
          >
            <span>⚖️</span> Authority
          </button>

          <button
            onClick={() => onSwitchPersona('admin1')}
            disabled={loading}
            style={{
              background: user?.role === 'Administrator' ? 'rgba(245, 158, 11, 0.25)' : 'transparent',
              border: `1px solid ${user?.role === 'Administrator' ? 'var(--warning)' : 'transparent'}`,
              color: user?.role === 'Administrator' ? '#fbbf24' : 'var(--text-secondary)',
              padding: '0.3rem 0.6rem',
              borderRadius: '6px',
              fontSize: '0.75rem',
              fontWeight: 600,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.3rem',
            }}
            data-testid="persona-switch-admin1"
            title="Administrator - Audit Logs & Backups"
          >
            <span>🔍</span> Admin
          </button>
        </div>

        {/* USER PROFILE & LOGOUT */}
        {user ? (
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
            <span className={`user-badge badge-${user.role.toLowerCase()}`} style={{ fontSize: '0.75rem' }} data-testid="current-user-role">
              {user.name} ({user.role})
            </span>
            <button
              onClick={onLogout}
              className="btn btn-secondary"
              style={{ padding: '0.35rem 0.7rem', fontSize: '0.8rem' }}
              data-testid="btn-logout"
            >
              Logout
            </button>
          </div>
        ) : (
          <button
            onClick={onOpenManualLogin}
            className="btn btn-primary"
            style={{ padding: '0.4rem 0.9rem', fontSize: '0.85rem' }}
            data-testid="btn-open-login"
          >
            Sign In
          </button>
        )}
      </div>
    </header>
  );
};
