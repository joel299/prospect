import { useEffect, useMemo, useState } from 'react';
import { Pulse, ArrowsClockwise, Bell, ChatCircle, Gear, Lightning, MagnifyingGlass, PaperPlaneRight, Pause, Robot, SignOut, UsersThree, WarningCircle } from '@phosphor-icons/react';
import './styles.css';

type Message = { id: string; conversation_id: string; direction: 'inbound' | 'outbound'; content: string; created_at: string; status?: string };
type Conversation = { id: string; lead_id?: string; channel?: string; agent_enabled?: boolean; updated_at?: string };
const API_BASE = (import.meta.env.VITE_API_URL || 'https://chat-api.iainfinito.com.br').replace(/\/$/, '');
const nav = [{ label: 'Inbox', icon: ChatCircle }, { label: 'Leads', icon: UsersThree }, { label: 'Follow-ups', icon: ArrowsClockwise }, { label: 'Agent', icon: Robot }, { label: 'Automations', icon: Lightning }, { label: 'Integrations', icon: Gear }, { label: 'Test Lab', icon: Pulse }, { label: 'Observability', icon: WarningCircle }, { label: 'Settings', icon: Gear }];

function time(value?: string) { return value ? new Date(value).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' }) : ''; }

export function App() {
  const [active, setActive] = useState('Inbox');
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [items, setItems] = useState<Message[]>([]);
  const [draft, setDraft] = useState('');
  const [processing, setProcessing] = useState(false);
  const [connected, setConnected] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [query, setQuery] = useState('');

  const loadConversations = async () => {
    setLoading(true); setError('');
    try {
      const response = await fetch(`${API_BASE}/api/v1/conversations`);
      if (!response.ok) throw new Error(`Falha ao carregar conversas (${response.status})`);
      const data = await response.json() as { conversations?: Conversation[] };
      const next = data.conversations || [];
      setConversations(next);
      setSelectedId((current) => next.some((item) => item.id === current) ? current : next[0]?.id || '');
    } catch (reason) { setError(reason instanceof Error ? reason.message : 'Falha ao carregar conversas'); }
    finally { setLoading(false); }
  };

  useEffect(() => { void loadConversations(); const wsBase = API_BASE.replace(/^http/, 'ws'); const ws = new WebSocket(`${wsBase}/ws/v1`); ws.onopen = () => setConnected(true); ws.onclose = () => setConnected(false); ws.onerror = () => setConnected(false); ws.onmessage = () => { if (selectedId) void loadMessages(selectedId); }; return () => ws.close(); }, []);
  useEffect(() => { if (selectedId) void loadMessages(selectedId); else setItems([]); }, [selectedId]);

  const loadMessages = async (id: string) => {
    try { const response = await fetch(`${API_BASE}/api/v1/conversations/${encodeURIComponent(id)}/messages`); if (!response.ok) throw new Error(`Falha ao carregar mensagens (${response.status})`); const data = await response.json() as { messages?: Message[] }; setItems(data.messages || []); }
    catch (reason) { setError(reason instanceof Error ? reason.message : 'Falha ao carregar mensagens'); }
  };

  const selected = conversations.find((item) => item.id === selectedId);
  const filtered = useMemo(() => conversations.filter((item) => `${item.id} ${item.lead_id || ''}`.toLowerCase().includes(query.toLowerCase())), [conversations, query]);
  const send = async () => {
    const content = draft.trim(); if (!content || !selected || processing) return;
    setProcessing(true); setError('');
    try {
      const response = await fetch(`${API_BASE}/api/v1/conversations/${encodeURIComponent(selected.id)}/send`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ lead_id: selected.lead_id, content }) });
      if (!response.ok) { const body = await response.json().catch(() => null) as { error?: { message?: string } } | null; throw new Error(body?.error?.message || `Falha ao enviar (${response.status})`); }
      setDraft(''); await loadMessages(selected.id);
    } catch (reason) { setError(reason instanceof Error ? reason.message : 'Falha ao enviar mensagem'); }
    finally { setProcessing(false); }
  };

  return <div className="shell"><aside className="sidebar"><div className="brand"><b>H</b><span><strong>HERMES</strong><small>OPERATION CONSOLE</small></span></div><p className="workspace">WORKSPACE <em>PROSPECT</em></p><nav aria-label="Navegação principal">{nav.map(({ label, icon: Icon }) => <button className={active === label ? 'nav active' : 'nav'} key={label} onClick={() => setActive(label)}><Icon size={19}/><span>{label}</span></button>)}</nav><div className="side-bottom"><p><i className={connected ? 'dot live' : 'dot'} />{connected ? 'API conectada' : 'API offline'}</p><button className="nav"><SignOut size={19}/><span>Sair</span></button></div></aside><main><header><div><small>OPERAÇÃO / {active.toUpperCase()}</small><h1>{active}</h1></div><div className="actions"><span className="environment"><i className="dot live"/> PRODUÇÃO CONTROLADA</span><button className="icon" aria-label="Notificações"><Bell size={20}/></button><span className="avatar">JQ</span></div></header>{active === 'Inbox' ? <Inbox conversations={filtered} selected={selected} selectedId={selectedId} setSelectedId={setSelectedId} items={items} draft={draft} setDraft={setDraft} send={send} processing={processing} loading={loading} error={error} query={query} setQuery={setQuery} refresh={loadConversations}/> : <div className="empty-page"><WarningCircle size={32}/><small>BACKEND PENDENTE</small><h2>{active}</h2><p>Esta aba ainda não possui endpoint real conectado. Nenhum dado fictício é exibido.</p><code>Contrato necessário antes da implementação</code></div>}</main></div>;
}

function Inbox({ conversations, selected, selectedId, setSelectedId, items, draft, setDraft, send, processing, loading, error, query, setQuery, refresh }: { conversations: Conversation[]; selected?: Conversation; selectedId: string; setSelectedId: (id: string) => void; items: Message[]; draft: string; setDraft: (value: string) => void; send: () => void; processing: boolean; loading: boolean; error: string; query: string; setQuery: (value: string) => void; refresh: () => void }) {
  return <section className="inbox"><aside className="conversations"><div className="panel-title"><div><small>ATENDIMENTO</small><h2>Conversas <sup>{conversations.length}</sup></h2></div><button className="icon" aria-label="Atualizar conversas" onClick={refresh}><ArrowsClockwise size={19}/></button></div><label className="search"><MagnifyingGlass size={17}/><input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Buscar conversa" aria-label="Buscar conversa"/></label>{loading && <p className="empty-state">Carregando conversas...</p>}{!loading && conversations.length === 0 && <p className="empty-state">Nenhuma conversa real recebida.</p>}{conversations.map((conversation) => <button className={selectedId === conversation.id ? 'conversation selected' : 'conversation'} key={conversation.id} onClick={() => setSelectedId(conversation.id)}><span className="contact">{(conversation.lead_id || conversation.id).slice(0, 2).toUpperCase()}</span><span><b>{conversation.lead_id || conversation.id}</b><small>{conversation.channel || 'Ryze'}</small></span><time>{time(conversation.updated_at)}</time></button>)}</aside><div className="chat"><div className="chat-head"><span className="contact">{(selected?.lead_id || selected?.id || '--').slice(0, 2).toUpperCase()}</span><div><b>{selected?.lead_id || selected?.id || 'Nenhuma conversa'}</b><small><i className="dot live"/> {selected ? 'Conversa real' : 'Aguardando inbound'}</small></div><button className="pause" disabled={!selected}><Pause size={16}/>Pausar agente</button></div>{error && <p className="error-banner">{error}</p>}<div className="stream" aria-live="polite">{items.length === 0 && selected && <p className="empty-state">Nenhuma mensagem nesta conversa.</p>}{items.map((item) => <div className={item.direction === 'outbound' ? 'message outbound' : 'message'} key={item.id}><div><span>{item.content}</span><small>{time(item.created_at)} {item.direction === 'outbound' && <em>{item.status || 'accepted'}</em>}</small></div></div>)}{processing && <p className="activity"><i className="pulse"/> Enviando pelo Ryze...</p>}</div><div className="composer"><textarea rows={1} value={draft} disabled={!selected || processing} onChange={(event) => setDraft(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter' && !event.shiftKey) { event.preventDefault(); send(); } }} placeholder={selected ? 'Escreva uma mensagem...' : 'Aguardando conversa real...'} aria-label="Mensagem"/><footer><span><Robot size={16}/>Agente server-side</span><button className="send" disabled={!selected || !draft.trim() || processing} onClick={send} aria-label="Enviar mensagem"><PaperPlaneRight size={18} weight="fill"/></button></footer></div></div><aside className="inspector"><div className="panel-title"><small>CONVERSA</small></div>{selected ? <><div className="profile"><span className="contact big">{(selected.lead_id || selected.id).slice(0, 2).toUpperCase()}</span><h2>{selected.lead_id || selected.id}</h2><small>{selected.channel || 'Ryze'}</small></div><section><label>STATUS</label><p><i className="dot live"/>{selected.agent_enabled ? 'Agente ativo' : 'Agente pausado'}</p></section></> : <p className="empty-state">Selecione uma conversa para ver detalhes.</p>}</aside></section>;
}
