import { useState, useEffect, useCallback } from 'react';
import type { User, Election, Candidate, AuditLog, MerkleStatus, BallotReceipt } from './types';
import { Header } from './components/Header';
import { ArchitectureDiagram } from './components/ArchitectureDiagram';
import { VoterView } from './components/VoterView';
import { OfficerView } from './components/OfficerView';
import { CountingView } from './components/CountingView';
import { AdminView } from './components/AdminView';
import { FacultyWalkthroughModal } from './components/FacultyWalkthroughModal';
import { JiraTestMatrixModal } from './components/JiraTestMatrixModal';
import { ManualLoginModal } from './components/ManualLoginModal';

// Unified API base URL (Vite proxies /api to backend :8080)
const API_BASE = '/api';

export default function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('vote_token'));
  const [user, setUser] = useState<User | null>(null);
  const [activeScreen, setActiveScreen] = useState<string>('home');

  // Application Data States
  const [elections, setElections] = useState<Election[]>([]);
  const [selectedElection, setSelectedElection] = useState<Election | null>(null);
  const [candidates, setCandidates] = useState<Candidate[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [talliedResults, setTalliedResults] = useState<any[]>([]);
  const [merkleStatus, setMerkleStatus] = useState<MerkleStatus | null>(null);
  const [receipt, setReceipt] = useState<BallotReceipt | null>(null);

  // Modals & Panels
  const [walkthroughOpen, setWalkthroughOpen] = useState(false);
  const [jiraOpen, setJiraOpen] = useState(false);
  const [manualLoginOpen, setManualLoginOpen] = useState(false);
  const [showArchDiagram, setShowArchDiagram] = useState(true);

  // Status & Feedback
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const showNotification = useCallback((msg: string, isSuccess = false) => {
    if (isSuccess) {
      setSuccess(msg);
      setError(null);
    } else {
      setError(msg);
      setSuccess(null);
    }
    setTimeout(() => {
      setError(null);
      setSuccess(null);
    }, 6000);
  }, []);

  // Fetch Current User Profile
  const fetchMe = useCallback(async (jwtToken: string) => {
    try {
      const resp = await fetch(`${API_BASE}/auth/me`, {
        headers: { 'Authorization': `Bearer ${jwtToken}` }
      });
      if (resp.ok) {
        const userData: User = await resp.json();
        setUser(userData);
      } else {
        handleLogout();
      }
    } catch {
      showNotification("Backend server offline or unreachable. Ensure 'go run main.go' is running.");
    }
  }, [showNotification]);

  // Fetch Elections List
  const fetchElections = useCallback(async () => {
    try {
      const resp = await fetch(`${API_BASE}/elections`);
      if (resp.ok) {
        const data = await resp.json();
        setElections(data || []);
      }
    } catch {
      // Background retry
    }
  }, []);

  // Fetch Audit Logs
  const fetchAuditLogs = useCallback(async (jwtToken: string) => {
    try {
      const resp = await fetch(`${API_BASE}/audit/logs`, {
        headers: { 'Authorization': `Bearer ${jwtToken}` }
      });
      if (resp.ok) {
        const data = await resp.json();
        setAuditLogs(data || []);
      }
    } catch {
      // Continue
    }
  }, []);

  // Load candidate list for an election
  const loadCandidates = async (electionId: string): Promise<Candidate[]> => {
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}`);
      if (resp.ok) {
        const data = await resp.json();
        return data.candidates || [];
      }
    } catch {
      showNotification("Failed to fetch candidate list.");
    }
    return [];
  };

  // Instant Persona Switcher
  const handleQuickPersonaSwitch = useCallback(async (persona: 'voter1' | 'voter2' | 'officer1' | 'counting1' | 'admin1') => {
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: persona, password: 'password123' })
      });
      const data = await resp.json();
      if (resp.ok) {
        localStorage.setItem('vote_token', data.token);
        setToken(data.token);
        await fetchMe(data.token);
        await fetchElections();
        setActiveScreen('home');
        showNotification(`Switched persona to: ${persona.toUpperCase()}`, true);
      } else {
        showNotification(data.error || "Persona switch failed.");
      }
    } catch {
      showNotification("Failed to connect to backend server. Make sure Go backend is running on port 8080.");
    } finally {
      setLoading(false);
    }
  }, [fetchMe, fetchElections, showNotification]);

  useEffect(() => {
    if (token) {
      fetchMe(token);
    } else {
      // Default to auto-login as voter1 for instant preview
      handleQuickPersonaSwitch('voter1');
    }
    fetchElections();
  }, [token, fetchMe, fetchElections, handleQuickPersonaSwitch]);

  useEffect(() => {
    if (token && user?.role === 'Administrator') {
      fetchAuditLogs(token);
    }
  }, [token, user?.role, fetchAuditLogs]);

  // Manual Login Handler
  const handleManualLogin = async (userId: string, pass: string): Promise<boolean> => {
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, password: pass })
      });
      const data = await resp.json();
      if (resp.ok) {
        localStorage.setItem('vote_token', data.token);
        setToken(data.token);
        await fetchMe(data.token);
        await fetchElections();
        showNotification(`Welcome, ${userId}!`, true);
        return true;
      } else {
        showNotification(data.error || "Login failed.");
        return false;
      }
    } catch {
      showNotification("Authentication server offline.");
      return false;
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('vote_token');
    setToken(null);
    setUser(null);
    setReceipt(null);
    setMerkleStatus(null);
    setTalliedResults([]);
    setActiveScreen('home');
  };

  // 1-Click Reset Demo Environment
  const handleResetDemo = async () => {
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/demo/reset`, { method: 'POST' });
      const data = await resp.json();
      if (resp.ok) {
        showNotification(data.message || "Demo reset complete! Clean database loaded.", true);
        setReceipt(null);
        setMerkleStatus(null);
        setTalliedResults([]);
        await fetchElections();
        if (token) {
          await fetchMe(token);
        }
      } else {
        showNotification(data.error || "Failed to reset demo.");
      }
    } catch {
      showNotification("Backend connection error during demo reset.");
    } finally {
      setLoading(false);
    }
  };

  // --- Voter Functions ---
  const handleSelectElectionForVoting = async (election: Election) => {
    setSelectedElection(election);
    setReceipt(null);
    setLoading(true);
    const candList = await loadCandidates(election.election_id);
    setCandidates(candList);
    setActiveScreen('ballot');
    setLoading(false);
  };

  const handleCastVote = async (electionId: string, candidateId: string): Promise<boolean> => {
    if (!token) return false;
    setLoading(true);
    try {
      // 1. Request anonymous blind token from Auth Server
      const tokenResp = await fetch(`${API_BASE}/auth/token/request`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ election_id: electionId })
      });
      const tokenData = await tokenResp.json();
      if (!tokenResp.ok) {
        showNotification(tokenData.error || "Failed to acquire anonymous voting token.");
        setLoading(false);
        return false;
      }

      const anonymousToken = tokenData.token;

      // 2. Cast encrypted vote to Voting Server using anonymous token
      const voteResp = await fetch(`${API_BASE}/vote/cast`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          token: anonymousToken,
          candidate_id: candidateId,
          election_id: electionId
        })
      });
      const voteData = await voteResp.json();
      if (voteResp.ok) {
        setReceipt({
          receipt_h: voteData.receipt_h,
          root_hash: voteData.root_hash,
          timestamp: voteData.timestamp,
        });
        showNotification("Vote cast and added to append-only Merkle ledger!", true);

        // Update local user status
        if (user) {
          setUser({ ...user, has_voted: true });
        }
        setActiveScreen('home');
        return true;
      } else {
        showNotification(voteData.error || "Failed to cast vote.");
        return false;
      }
    } catch {
      showNotification("Voting service connection error.");
      return false;
    } finally {
      setLoading(false);
    }
  };

  // --- Officer Functions ---
  const handleCreateElection = async (formData: { title: string; desc: string; start: string; end: string; passphrase: string }) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          title: formData.title,
          description: formData.desc,
          start_time: new Date(formData.start).toISOString(),
          end_time: new Date(formData.end).toISOString(),
          passphrase: formData.passphrase
        })
      });
      const data = await resp.json();
      if (resp.ok) {
        showNotification("Election created & RSA-2048 keypair generated successfully!", true);
        await fetchElections();
      } else {
        showNotification(data.error || "Failed to create election.");
      }
    } catch {
      showNotification("Server error creating election.");
    } finally {
      setLoading(false);
    }
  };

  const handleOpenElection = async (electionId: string) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/open`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (resp.ok) {
        showNotification("Election is now ACTIVE for voting!", true);
        await fetchElections();
      } else {
        const data = await resp.json();
        showNotification(data.error || "Failed to open election.");
      }
    } catch {
      showNotification("Server error opening election.");
    } finally {
      setLoading(false);
    }
  };

  const handleCloseElection = async (electionId: string) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/close`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await resp.json();
      if (resp.ok) {
        showNotification("Election closed! Merkle root signed & archived.", true);
        await fetchElections();
      } else {
        showNotification(data.error || "Failed to close election.");
      }
    } catch {
      showNotification("Server error closing election.");
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyMerkle = async (electionId: string) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/verify-merkle`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data: MerkleStatus = await resp.json();
      setMerkleStatus(data);
      if (data.verification_pass) {
        showNotification(data.message, true);
      } else {
        showNotification(data.message, false);
      }
    } catch {
      showNotification("Server error verifying Merkle root.");
    } finally {
      setLoading(false);
    }
  };

  const handleAddCandidate = async (electionId: string, name: string, party: string) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/candidates`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ name, party })
      });
      if (resp.ok) {
        showNotification(`Candidate ${name} added!`, true);
      } else {
        const data = await resp.json();
        showNotification(data.error || "Failed to add candidate.");
      }
    } catch {
      showNotification("Server error adding candidate.");
    } finally {
      setLoading(false);
    }
  };

  // --- Counting Authority Functions ---
  const handleTallyVotes = async (electionId: string, authorityPassphrase: string) => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/tally`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ passphrase: authorityPassphrase })
      });
      const data = await resp.json();
      if (resp.ok) {
        setTalliedResults(data.results || []);
        showNotification("Tally complete! Ballots decrypted with AES key and counted.", true);
        await fetchElections();
      } else {
        showNotification(data.error || "Failed to tally election.");
      }
    } catch {
      showNotification("Server error tallying votes.");
    } finally {
      setLoading(false);
    }
  };

  const fetchTalliedResults = async (electionId: string) => {
    try {
      const resp = await fetch(`${API_BASE}/elections/${electionId}/results`);
      if (resp.ok) {
        const data = await resp.json();
        setTalliedResults(data.results || []);
      }
    } catch {
      // Continue
    }
  };

  // --- Admin Functions ---
  const handleTriggerBackup = async () => {
    if (!token) return;
    setLoading(true);
    try {
      const resp = await fetch(`${API_BASE}/audit/backup`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await resp.json();
      if (resp.ok) {
        showNotification(`Database backup created: ${data.file}`, true);
        await fetchAuditLogs(token);
      } else {
        showNotification(data.error || "Backup failed.");
      }
    } catch {
      showNotification("Audit backup server error.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-container">
      {/* HEADER WITH PERSONA BAR */}
      <Header
        user={user}
        onSwitchPersona={handleQuickPersonaSwitch}
        onOpenWalkthrough={() => setWalkthroughOpen(true)}
        onOpenJira={() => setJiraOpen(true)}
        onResetDemo={handleResetDemo}
        onOpenManualLogin={() => setManualLoginOpen(true)}
        onLogout={handleLogout}
        loading={loading}
      />

      <main className="main-content">
        {/* NOTIFICATION BANNERS */}
        {error && <div className="status-banner status-error animate-fade-in" data-testid="status-error-banner">{error}</div>}
        {success && <div className="status-banner status-success animate-fade-in" data-testid="status-success-banner">{success}</div>}

        {/* ARCHITECTURAL PRIVACY DIAGRAM */}
        <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '0.4rem' }}>
          <button
            onClick={() => setShowArchDiagram(!showArchDiagram)}
            className="btn btn-secondary"
            style={{ padding: '0.2rem 0.6rem', fontSize: '0.75rem', border: 'none', color: 'var(--text-muted)' }}
          >
            {showArchDiagram ? '▲ Hide Privacy Architecture Wall' : '▼ Show Privacy Architecture Wall'}
          </button>
        </div>
        {showArchDiagram && <ArchitectureDiagram />}

        {/* ROLE BASED CONTENT SCREEN */}
        {user ? (
          <>
            {user.role === 'ElectionOfficer' && (
              <OfficerView
                elections={elections}
                loading={loading}
                onRefreshElections={fetchElections}
                onCreateElection={handleCreateElection}
                onOpenElection={handleOpenElection}
                onCloseElection={handleCloseElection}
                onVerifyMerkle={handleVerifyMerkle}
                onAddCandidate={handleAddCandidate}
                onLoadCandidates={loadCandidates}
                merkleStatus={merkleStatus}
                token={token}
              />
            )}

            {user.role === 'CountingAuthority' && (
              <CountingView
                elections={elections}
                loading={loading}
                onRefreshElections={fetchElections}
                onTallyVotes={handleTallyVotes}
                talliedResults={talliedResults}
                onFetchResults={fetchTalliedResults}
              />
            )}

            {user.role === 'Administrator' && (
              <AdminView
                auditLogs={auditLogs}
                loading={loading}
                onRefreshLogs={() => token && fetchAuditLogs(token)}
                onTriggerBackup={handleTriggerBackup}
              />
            )}

            {user.role === 'Voter' && (
              <VoterView
                user={user}
                elections={elections}
                receipt={receipt}
                loading={loading}
                onSelectElectionForVoting={handleSelectElectionForVoting}
                onViewResults={(elect) => {
                  setSelectedElection(elect);
                  fetchTalliedResults(elect.election_id);
                  setActiveScreen('results');
                }}
                onCastVote={handleCastVote}
                candidates={candidates}
                selectedElection={selectedElection}
                activeScreen={activeScreen}
                onBackToElections={() => setActiveScreen('home')}
              />
            )}
          </>
        ) : (
          <div style={{ textAlign: 'center', padding: '4rem 1rem' }}>
            <h2 style={{ fontSize: '1.8rem', marginBottom: '0.8rem' }}>Authenticating Secure Portal...</h2>
            <p style={{ color: 'var(--text-secondary)' }}>Click any persona chip in the top bar to begin.</p>
          </div>
        )}
      </main>

      {/* FACULTY PRESENTATION STORYBOARD MODAL */}
      <FacultyWalkthroughModal
        isOpen={walkthroughOpen}
        onClose={() => setWalkthroughOpen(false)}
        onSwitchPersona={handleQuickPersonaSwitch}
        onResetDemo={handleResetDemo}
      />

      {/* JIRA / QA AUTOMATION TEST MATRIX MODAL */}
      <JiraTestMatrixModal
        isOpen={jiraOpen}
        onClose={() => setJiraOpen(false)}
      />

      {/* MANUAL LOGIN MODAL */}
      <ManualLoginModal
        isOpen={manualLoginOpen}
        onClose={() => setManualLoginOpen(false)}
        onLogin={handleManualLogin}
        loading={loading}
      />
    </div>
  );
}
