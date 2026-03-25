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

import { describe, expect, test } from 'bun:test';

import {
  applyClaudeSessionSettings,
  getClaudeSessionFormValues,
} from './claudeSessionSettings';

describe('claudeSessionSettings', () => {
  test('keeps configured values when loading Claude session settings', () => {
    expect(
      getClaudeSessionFormValues({
        claude_max_sessions: 7,
        claude_session_ttl_minutes: 3,
      }),
    ).toEqual({
      claude_max_sessions: 7,
      claude_session_ttl_minutes: 3,
    });
  });

  test('falls back to unlimited sessions and default ttl for missing or invalid values', () => {
    expect(
      getClaudeSessionFormValues({
        claude_max_sessions: 0,
        claude_session_ttl_minutes: 0,
      }),
    ).toEqual({
      claude_max_sessions: 0,
      claude_session_ttl_minutes: 30,
    });
  });

  test('writes numeric Claude session settings back into channel settings', () => {
    expect(
      applyClaudeSessionSettings(
        { allow_inference_geo: true },
        {
          claude_max_sessions: '3',
          claude_session_ttl_minutes: '45',
        },
      ),
    ).toEqual({
      allow_inference_geo: true,
      claude_max_sessions: 3,
      claude_session_ttl_minutes: 45,
    });
  });

  test('clears Claude session settings when using unlimited/default values', () => {
    expect(
      applyClaudeSessionSettings(
        {
          claude_max_sessions: 9,
          claude_session_ttl_minutes: 60,
          allow_inference_geo: true,
        },
        {
          claude_max_sessions: 0,
          claude_session_ttl_minutes: '',
        },
      ),
    ).toEqual({
      allow_inference_geo: true,
    });
  });
});
