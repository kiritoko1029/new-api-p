/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

const DEFAULT_CLAUDE_MAX_SESSIONS = 0;
const DEFAULT_CLAUDE_SESSION_TTL_MINUTES = 30;

const parseInteger = (value) => {
  if (value === '' || value === null || value === undefined) {
    return null;
  }

  const normalized =
    typeof value === 'number' ? value : Number.parseInt(String(value), 10);

  return Number.isFinite(normalized) ? normalized : null;
};

export const getClaudeSessionFormValues = (settings = {}) => {
  const maxSessions = parseInteger(settings.claude_max_sessions);
  const sessionTtlMinutes = parseInteger(settings.claude_session_ttl_minutes);

  return {
    claude_max_sessions:
      maxSessions !== null && maxSessions >= 0
        ? maxSessions
        : DEFAULT_CLAUDE_MAX_SESSIONS,
    claude_session_ttl_minutes:
      sessionTtlMinutes !== null && sessionTtlMinutes > 0
        ? sessionTtlMinutes
        : DEFAULT_CLAUDE_SESSION_TTL_MINUTES,
  };
};

export const applyClaudeSessionSettings = (settings = {}, formValues = {}) => {
  const nextSettings = { ...settings };
  const maxSessions = parseInteger(formValues.claude_max_sessions);
  const sessionTtlMinutes = parseInteger(formValues.claude_session_ttl_minutes);

  if (maxSessions !== null && maxSessions > 0) {
    nextSettings.claude_max_sessions = maxSessions;
  } else {
    delete nextSettings.claude_max_sessions;
  }

  if (sessionTtlMinutes !== null && sessionTtlMinutes > 0) {
    nextSettings.claude_session_ttl_minutes = sessionTtlMinutes;
  } else {
    delete nextSettings.claude_session_ttl_minutes;
  }

  return nextSettings;
};
