import React from 'react';
import type { Theme } from '../hooks/useTheme';

interface ThemeToggleProps {
  theme: Theme;
  toggleTheme: () => void;
}

export const ThemeToggle: React.FC<ThemeToggleProps> = ({ theme, toggleTheme }) => (
  <button
    className="theme-toggle"
    onClick={toggleTheme}
    title="Переключить тему"
  >
    {theme === 'light' ? '🌙' : '☀️'}
  </button>
);