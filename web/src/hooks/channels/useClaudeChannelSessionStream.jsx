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

import { useEffect, useMemo, useRef, useState } from 'react';
import { SSE } from 'sse.js';
import { API, authHeader, getUserIdFromLocalStorage } from '../../helpers';

const toUniquePositiveIds = (ids = []) => {
  const seen = new Set();
  return ids
    .map((id) => Number(id))
    .filter((id) => Number.isInteger(id) && id > 0)
    .filter((id) => {
      if (seen.has(id)) {
        return false;
      }
      seen.add(id);
      return true;
    })
    .sort((a, b) => a - b);
};

const buildStreamUrl = ({ channelId, channelIds, details }) => {
  const params = new URLSearchParams();
  if (channelId) {
    params.set('channel_id', String(channelId));
  } else if (channelIds && channelIds.length > 0) {
    params.set('channel_ids', channelIds.join(','));
  }
  // If channelIds is undefined or empty, fetch ALL channels (no channel_ids param)
  if (details) {
    params.set('details', 'true');
  }

  const baseUrl = API.defaults.baseURL || '';
  const query = params.toString();
  return `${baseUrl}/api/channel/session_limits/stream${query ? `?${query}` : ''}`;
};

const toSessionStateMap = (states = []) => {
  const nextMap = {};
  states.forEach((state) => {
    const channelId = Number(state?.channel_id);
    if (!Number.isInteger(channelId) || channelId <= 0) {
      return;
    }
    nextMap[channelId] = state;
  });
  return nextMap;
};

export const useClaudeChannelSessionStream = ({
  channelId,
  channelIds = [],
  details = false,
  enabled = true,
} = {}) => {
  const sourceRef = useRef(null);
  const [sessionStateMap, setSessionStateMap] = useState({});
  const normalizedChannelId = Number(channelId) || 0;
  const normalizedChannelIds = useMemo(
    () => toUniquePositiveIds(channelIds),
    [channelIds],
  );
  const channelIdsKey = normalizedChannelIds.join(',');

  useEffect(() => {
    if (!enabled) {
      setSessionStateMap({});
      return undefined;
    }
    // Only skip if channelIds is explicitly undefined (not passed at all)
    // An empty channelIds array [] means "fetch all channels" for the dashboard
    if (!normalizedChannelId && channelIds === undefined) {
      setSessionStateMap({});
      return undefined;
    }

    const source = new SSE(
      buildStreamUrl({
        channelId: normalizedChannelId,
        channelIds: normalizedChannelIds,
        details,
      }),
      {
        headers: {
          'New-API-User': getUserIdFromLocalStorage(),
          ...authHeader(),
        },
      },
    );
    sourceRef.current = source;

    source.addEventListener('message', (event) => {
      if (!event?.data || event.data === '[DONE]') {
        return;
      }
      try {
        const payload = JSON.parse(event.data);
        setSessionStateMap(toSessionStateMap(payload?.channels || []));
      } catch (error) {
        console.error('Failed to parse channel session SSE payload:', error);
      }
    });

    source.addEventListener('error', (event) => {
      if (source.readyState === 2) {
        return;
      }
      if (event?.data) {
        console.error('Channel session stream error:', event.data);
      }
    });

    try {
      source.stream();
    } catch (error) {
      console.error('Failed to start channel session stream:', error);
    }

    return () => {
      source.close();
      if (sourceRef.current === source) {
        sourceRef.current = null;
      }
    };
  }, [channelIdsKey, details, enabled, normalizedChannelId]);

  return {
    sessionStateMap,
    sessionState:
      normalizedChannelId > 0
        ? sessionStateMap[normalizedChannelId] || null
        : null,
  };
};
