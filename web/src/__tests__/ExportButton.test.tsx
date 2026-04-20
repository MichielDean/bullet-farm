import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ExportButton } from '../components/ExportButton';

describe('ExportButton', () => {
  beforeEach(() => { localStorage.clear(); });
  afterEach(() => { vi.restoreAllMocks(); });

  it('renders export button', () => {
    render(<ExportButton />);
    expect(screen.getByText('Export')).toBeInTheDocument();
  });

  it('toggles dropdown on click', () => {
    render(<ExportButton />);
    const btn = screen.getByText('Export');
    fireEvent.click(btn);
    expect(screen.getByText('JSON')).toBeInTheDocument();
    expect(screen.getByText('CSV')).toBeInTheDocument();
  });
});