import { useCallback, useEffect, useRef, useState } from 'react';
import { getAuthParams } from '../hooks/useAuth';

const MAX_BUFFER_SIZE = 50 * 1024;

interface PeekPanelProps {
  aqueductName: string;
  onClose: () => void;
}

export function PeekPanel({ aqueductName, onClose }: PeekPanelProps) {
  const [output, setOutput] = useState<string>('');
  const [connected, setConnected] = useState(false);
  const terminalRef = useRef<HTMLPreElement>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const mountedRef = useRef(true);

  const appendOutput = useCallback((chunk: string) => {
    setOutput((prev) => {
      const next = prev + chunk;
      if (next.length > MAX_BUFFER_SIZE) {
        return next.slice(next.length - MAX_BUFFER_SIZE);
      }
      return next;
    });
  }, []);

  useEffect(() => {
    mountedRef.current = true;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const authParams = getAuthParams();
    const wsUrl = `${protocol}//${window.location.host}/ws/aqueducts/${encodeURIComponent(aqueductName)}/peek${authParams ? '?' + authParams : ''}`;
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => { if (mountedRef.current) setConnected(true); };
    ws.onmessage = (e) => {
      if (mountedRef.current) appendOutput(e.data as string);
    };
    ws.onclose = () => { if (mountedRef.current) setConnected(false); };
    ws.onerror = () => { if (mountedRef.current) setConnected(false); };

    return () => {
      mountedRef.current = false;
      ws.close();
      wsRef.current = null;
    };
  }, [aqueductName, appendOutput]);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output]);

  return (
    <div className="fixed inset-y-0 right-0 w-full md:w-[600px] bg-cistern-bg border-l border-cistern-border shadow-2xl z-50 flex flex-col">
      <div className="flex items-center justify-between px-4 py-3 border-b border-cistern-border">
        <div className="flex items-center gap-3">
          <h3 className="font-mono text-cistern-accent">{aqueductName}</h3>
          <span className="text-xs text-cistern-muted">Peek</span>
          <div className={`w-2 h-2 rounded-full ${connected ? 'bg-cistern-green' : 'bg-cistern-red'}`} />
        </div>
        <button
          onClick={onClose}
          className="text-cistern-muted hover:text-cistern-fg transition-colors text-lg leading-none"
        >
          ×
        </button>
      </div>
      <pre
        ref={terminalRef}
        className="flex-1 overflow-auto p-4 font-mono text-xs text-cistern-green bg-cistern-bg whitespace-pre-wrap break-all"
      >
        {output || 'Connecting…'}
      </pre>
    </div>
  );
}