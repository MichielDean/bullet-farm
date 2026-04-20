import { useState, useRef, useEffect } from 'react';

interface RenameInputProps {
  value: string;
  onSave: (newTitle: string) => Promise<void>;
  className?: string;
}

export function RenameInput({ value, onSave, className = '' }: RenameInputProps) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(value);
  const [saving, setSaving] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (editing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [editing]);

  const handleSave = async () => {
    const trimmed = draft.trim();
    if (!trimmed || trimmed === value) {
      setDraft(value);
      setEditing(false);
      return;
    }
    setSaving(true);
    try {
      await onSave(trimmed);
      setEditing(false);
    } catch {
      setDraft(value);
    } finally {
      setSaving(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleSave();
    } else if (e.key === 'Escape') {
      setDraft(value);
      setEditing(false);
    }
  };

  if (!editing) {
    return (
      <h1
        className={`text-xl font-mono font-bold text-cistern-fg truncate cursor-pointer hover:text-cistern-accent transition-colors ${className}`}
        onClick={() => { setDraft(value); setEditing(true); }}
        title="Click to rename"
        role="button"
        tabIndex={0}
        onKeyDown={(e) => { if (e.key === 'Enter') { setDraft(value); setEditing(true); } }}
      >
        {value}
      </h1>
    );
  }

  return (
    <input
      ref={inputRef}
      type="text"
      value={draft}
      onChange={(e) => setDraft(e.target.value)}
      onBlur={handleSave}
      onKeyDown={handleKeyDown}
      disabled={saving}
      className={`text-xl font-mono font-bold text-cistern-fg bg-cistern-bg border border-cistern-accent rounded px-2 py-1 ${className}`}
      aria-label="Rename droplet"
    />
  );
}