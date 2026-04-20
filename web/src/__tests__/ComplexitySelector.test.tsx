import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ComplexitySelector } from '../components/ComplexitySelector';

describe('ComplexitySelector', () => {
  it('renders all three complexity levels', () => {
    render(<ComplexitySelector value={1} onChange={vi.fn()} />);
    expect(screen.getByText('Standard (1)')).toBeInTheDocument();
    expect(screen.getByText('Full (2)')).toBeInTheDocument();
    expect(screen.getByText('Critical (3)')).toBeInTheDocument();
  });

  it('shows the selected value as checked', () => {
    render(<ComplexitySelector value={2} onChange={vi.fn()} />);
    const radios = screen.getAllByRole('radio');
    expect(radios[0]).not.toBeChecked();
    expect(radios[1]).toBeChecked();
    expect(radios[2]).not.toBeChecked();
  });

  it('calls onChange when a level is selected', () => {
    const onChange = vi.fn();
    render(<ComplexitySelector value={1} onChange={onChange} />);
    const criticalRadio = screen.getByRole('radio', { name: 'Critical (3)' });
    fireEvent.click(criticalRadio);
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it('shows pipeline stages for selected complexity', () => {
    render(<ComplexitySelector value={1} onChange={vi.fn()} />);
    expect(screen.getByText('implement')).toBeInTheDocument();
    expect(screen.getByText('delivery')).toBeInTheDocument();
  });

  it('shows all stages for critical complexity', () => {
    render(<ComplexitySelector value={3} onChange={vi.fn()} />);
    expect(screen.getByText('security-review')).toBeInTheDocument();
  });

  it('disables inputs when disabled prop is true', () => {
    render(<ComplexitySelector value={1} onChange={vi.fn()} disabled />);
    const radios = screen.getAllByRole('radio');
    radios.forEach((r) => expect(r).toBeDisabled());
  });
});