import React, { useState } from 'react';
import type { AuditLog } from '../types';

interface AdminViewProps {
  auditLogs: AuditLog[];
  loading: boolean;
  onRefreshLogs: () => void;
  onTriggerBackup: () => Promise<void>;
}

export const AdminView: React.FC<AdminViewProps> = ({
  auditLogs,
  loading,
  onRefreshLogs,
  onTriggerBackup,
}) => {
  const [searchTerm, setSearchTerm] = useState('');

  const filteredLogs = auditLogs.filter(log => 
    log.action.toLowerCase().includes(searchTerm.toLowerCase()) ||
    log.user_role.toLowerCase().includes(searchTerm.toLowerCase()) ||
    log.ip_address.includes(searchTerm) ||
    log.log_id.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }} data-testid="admin-console-view">
      {/* HEADER BANNER */}
      <div className="glass-card" style={{ padding: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem' }}>
        <div>
          <span style={{ fontSize: '0.8rem', color: 'var(--warning)', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            System Administrator & Auditor Portal (REQ-DB-4)
          </span>
          <h2 style={{ fontSize: '1.7rem', marginTop: '0.2rem' }}>Security Audit & Ledger Backup</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', marginTop: '0.3rem' }}>
            Review immutable append-only audit events and execute cryptographic database backups.
          </p>
        </div>

        <div style={{ display: 'flex', gap: '0.8rem' }}>
          <button 
            onClick={onTriggerBackup} 
            disabled={loading}
            className="btn btn-primary"
            data-testid="btn-trigger-backup"
          >
            💾 {loading ? 'Creating Backup...' : 'Trigger Database Backup'}
          </button>
          <button 
            onClick={onRefreshLogs} 
            disabled={loading}
            className="btn btn-secondary"
            data-testid="btn-refresh-audit-logs"
          >
            🔄 Refresh Logs
          </button>
        </div>
      </div>

      {/* AUDIT LOG TABLE CARD */}
      <div className="glass-card" style={{ padding: '2rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem', marginBottom: '1.5rem' }}>
          <div>
            <h3 style={{ fontSize: '1.25rem' }}>Immutable Security Audit Trail</h3>
            <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
              Total Events Logged: <strong>{auditLogs.length}</strong>
            </span>
          </div>

          {/* SEARCH INPUT */}
          <div style={{ width: '280px' }}>
            <input 
              type="text" 
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="input-field" 
              placeholder="🔍 Filter audit events..."
              data-testid="input-filter-audit"
            />
          </div>
        </div>

        <div style={{ overflowX: 'auto' }}>
          <table className="custom-table">
            <thead>
              <tr>
                <th style={{ width: '100px' }}>Log ID</th>
                <th>Security Action Event</th>
                <th style={{ width: '180px' }}>Timestamp</th>
                <th style={{ width: '140px' }}>Operator Role</th>
                <th style={{ width: '120px' }}>IP Address</th>
              </tr>
            </thead>
            <tbody>
              {filteredLogs.map(log => (
                <tr key={log.log_id} data-testid={`audit-row-${log.log_id}`}>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                    {log.log_id.slice(0, 8)}...
                  </td>
                  <td style={{ fontWeight: 500, color: log.action.includes('ALERT') || log.action.includes('WARNING') || log.action.includes('Failed') ? 'var(--error)' : 'var(--text-primary)' }}>
                    {log.action}
                  </td>
                  <td style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                    {new Date(log.timestamp).toLocaleString()}
                  </td>
                  <td>
                    <span 
                      className={`user-badge badge-${log.user_role.toLowerCase()}`} 
                      style={{ fontSize: '0.7rem', padding: '0.2rem 0.6rem' }}
                    >
                      {log.user_role}
                    </span>
                  </td>
                  <td style={{ color: 'var(--text-muted)', fontSize: '0.85rem', fontFamily: 'monospace' }}>
                    {log.ip_address}
                  </td>
                </tr>
              ))}

              {filteredLogs.length === 0 && (
                <tr>
                  <td colSpan={5} style={{ textAlign: 'center', color: 'var(--text-muted)', fontStyle: 'italic', padding: '2rem' }}>
                    {searchTerm ? 'No audit events match your search query.' : 'No audit events recorded yet.'}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
