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

  it('includes token query param when API key is stored', () => {
    localStorage.setItem('cistern_api_key', 'test-key-123');
    const openSpy = vi.spyOn(window, 'open').mockReturnValue(null);
    render(<ExportButton />);
    fireEvent.click(screen.getByText('Export'));
    fireEvent.click(screen.getByText('JSON'));
    expect(openSpy).toHaveBeenCalledTimes(1);
    const calledUrl = openSpy.mock.calls[0][0] as string;
    expect(calledUrl).toContain('token=test-key-123');
    expect(calledUrl).toContain('format=json');
  });

  it('omits token query param when no API key stored', () => {
    const openSpy = vi.spyOn(window, 'open').mockReturnValue(null);
    render(<ExportButton />);
    fireEvent.click(screen.getByText('Export'));
    fireEvent.click(screen.getByText('CSV'));
    expect(openSpy).toHaveBeenCalledTimes(1);
    const calledUrl = openSpy.mock.calls[0][0] as string;
    expect(calledUrl).not.toContain('token=');
    expect(calledUrl).toContain('format=csv');
  });
});