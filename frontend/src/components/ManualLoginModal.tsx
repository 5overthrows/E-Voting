import React, { useState } from 'react';

interface ManualLoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLogin: (userId: string, pass: string) => Promise<boolean>;
  loading: boolean;
}

export const ManualLoginModal: React.FC<ManualLoginModalProps> = ({
  isOpen,
  onClose,
  onLogin,
  loading,
}) => {
  const [userId, setUserId] = useState('');
  const [password, setPassword] = useState('');
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);
    const ok = await onLogin(userId, password);
    if (ok) {
      onClose();
    } else {
      setErrorMsg('Authentication failed. Check credentials or lockout status.');
    }
  };

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
      data-testid="manual-login-modal"
    >
      <div className="glass-card" style={{ maxWidth: '450px', width: '100%', padding: '2.5rem', border: '1px solid var(--primary-glow)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
          <h2 style={{ fontSize: '1.5rem' }}>Authenticate User</h2>
          <button onClick={onClose} className="btn btn-secondary" style={{ padding: '0.3rem 0.6rem', fontSize: '0.8rem' }}>
            ✕
          </button>
        </div>

        {errorMsg && (
          <div className="status-banner status-error" style={{ fontSize: '0.85rem' }}>
            {errorMsg}
          </div>
        )}

        <form onSubmit={handleSubmit} data-testid="form-manual-login">
          <div className="form-group">
            <label className="form-label">User ID / Voter ID</label>
            <input 
              type="text" 
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              className="input-field" 
              placeholder="e.g. voter1, officer1, admin1"
              required
              data-testid="manual-input-userid"
            />
          </div>

          <div className="form-group">
            <label className="form-label">Password</label>
            <input 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="input-field" 
              placeholder="••••••••"
              required
              data-testid="manual-input-password"
            />
          </div>

          <button 
            type="submit" 
            disabled={loading} 
            className="btn btn-primary" 
            style={{ width: '100%', marginTop: '1rem' }}
            data-testid="manual-btn-login"
          >
            {loading ? 'Authenticating...' : 'Sign In'}
          </button>
        </form>

        <div style={{ marginTop: '1.5rem', fontSize: '0.8rem', color: 'var(--text-muted)', textAlign: 'center', borderTop: '1px solid var(--glass-border)', paddingTop: '1rem' }}>
          Default Passwords: <code>password123</code>
        </div>
      </div>
    </div>
  );
};
