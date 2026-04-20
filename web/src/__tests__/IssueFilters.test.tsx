import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { IssueFilters } from '../components/IssueFilters';

describe('IssueFilters', () => {
  const issues = [
    { flagged_by: 'review', status: 'open' },
    { flagged_by: 'implement', status: 'open' },
    { flagged_by: 'review', status: 'resolved' },
  ];

  it('renders All and Open filter buttons', () => {
    render(
      <IssueFilters
        issues={issues}
        statusFilter="open"
        roleFilter=""
        sortOrder="newest"
        onStatusFilterChange={() => {}}
        onRoleFilterChange={() => {}}
        onSortOrderChange={() => {}}
      />
    );
    expect(screen.getByText('All')).toBeInTheDocument();
    expect(screen.getByText('Open')).toBeInTheDocument();
  });

  it('renders flagged-by filter buttons from issues', () => {
    render(
      <IssueFilters
        issues={issues}
        statusFilter="open"
        roleFilter=""
        sortOrder="newest"
        onStatusFilterChange={() => {}}
        onRoleFilterChange={() => {}}
        onSortOrderChange={() => {}}
      />
    );
    expect(screen.getByText('review')).toBeInTheDocument();
    expect(screen.getByText('implement')).toBeInTheDocument();
  });

  it('renders sort dropdown', () => {
    render(
      <IssueFilters
        issues={issues}
        statusFilter="open"
        roleFilter=""
        sortOrder="newest"
        onStatusFilterChange={() => {}}
        onRoleFilterChange={() => {}}
        onSortOrderChange={() => {}}
      />
    );
    expect(screen.getByText('Sort:')).toBeInTheDocument();
  });

  it('shows active state for selected filter', () => {
    render(
      <IssueFilters
        issues={issues}
        statusFilter="open"
        roleFilter=""
        sortOrder="newest"
        onStatusFilterChange={() => {}}
        onRoleFilterChange={() => {}}
        onSortOrderChange={() => {}}
      />
    );
    const openBtn = screen.getByText('Open');
    expect(openBtn.className).toContain('border-cistern-accent');
  });
});