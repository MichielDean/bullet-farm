interface ComplexitySelectorProps {
  value: number;
  onChange: (value: number) => void;
  disabled?: boolean;
}

const COMPLEXITY_LEVELS: Record<number, { label: string; stages: string[] }> = {
  1: { label: 'Standard', stages: ['implement', 'delivery'] },
  2: { label: 'Full', stages: ['implement', 'review', 'qa', 'docs', 'delivery'] },
  3: { label: 'Critical', stages: ['implement', 'review', 'qa', 'security-review', 'docs', 'delivery'] },
};

export function ComplexitySelector({ value, onChange, disabled }: ComplexitySelectorProps) {
  return (
    <div className="space-y-2">
      <label className="block text-xs font-mono text-cistern-muted uppercase tracking-wider mb-1">Complexity</label>
      <div className="space-y-2">
        {[1, 2, 3].map((level) => {
          const config = COMPLEXITY_LEVELS[level];
          const isSelected = value === level;
          return (
            <label
              key={level}
              className={`flex items-start gap-3 p-2 rounded-lg border cursor-pointer transition-colors ${
                isSelected
                  ? 'border-cistern-accent bg-cistern-accent/10'
                  : 'border-cistern-border hover:border-cistern-muted'
              } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              <input
                type="radio"
                name="complexity"
                value={level}
                checked={isSelected}
                onChange={() => onChange(level)}
                disabled={disabled}
                className="mt-0.5 accent-cistern-accent"
              />
              <div className="flex-1 min-w-0">
                <div className="text-sm font-mono text-cistern-fg font-medium">
                  {config.label} ({level})
                </div>
                {isSelected && (
                  <div className="flex items-center gap-1 mt-1 flex-wrap">
                    {config.stages.map((stage, i) => (
                      <span key={stage} className="flex items-center">
                        <span className="text-[10px] font-mono px-1 py-0.5 rounded bg-cistern-accent/20 text-cistern-accent">
                          {stage}
                        </span>
                        {i < config.stages.length - 1 && (
                          <span className="text-cistern-muted text-[10px] mx-0.5">→</span>
                        )}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </label>
          );
        })}
      </div>
    </div>
  );
}