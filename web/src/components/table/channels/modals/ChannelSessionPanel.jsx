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

import React from 'react';
import { useTranslation } from 'react-i18next';
import { Card, Empty, Space, Tag, Typography } from '@douyinfe/semi-ui';
import { timestamp2string } from '../../../../helpers';
import { useClaudeChannelSessionStream } from '../../../../hooks/channels/useClaudeChannelSessionStream';

const { Text } = Typography;

const formatRemaining = (seconds, t) => {
  const remaining = Math.max(0, Number(seconds) || 0);
  if (remaining < 60) {
    return t('{{count}} 秒', { count: remaining });
  }
  if (remaining < 3600) {
    return t('{{count}} 分钟', { count: Math.ceil(remaining / 60) });
  }
  return t('{{count}} 小时', { count: Math.ceil(remaining / 3600) });
};

const getSummaryColor = (state) => {
  if (!state) {
    return 'grey';
  }
  const active = Number(state.active_sessions) || 0;
  const max = Number(state.max_sessions) || 0;
  if (max > 0 && active >= max) {
    return 'red';
  }
  if (max > 0 && active >= Math.max(1, max - 1)) {
    return 'yellow';
  }
  return max > 0 ? 'green' : 'blue';
};

const ChannelSessionPanel = ({
  channelId,
  visible,
  fallbackMaxSessions = 0,
  fallbackTtlMinutes = 30,
}) => {
  const { t } = useTranslation();
  const { sessionState } = useClaudeChannelSessionStream({
    channelId,
    details: true,
    enabled: visible && Number(channelId) > 0,
  });

  if (!visible || !Number(channelId)) {
    return null;
  }

  const activeSessions = Number(sessionState?.active_sessions) || 0;
  const maxSessions = sessionState
    ? Number(sessionState?.max_sessions) || 0
    : Number(fallbackMaxSessions) || 0;
  const ttlMinutes = sessionState
    ? Number(sessionState?.ttl_minutes) || 30
    : Number(fallbackTtlMinutes) || 30;
  const sessions = Array.isArray(sessionState?.sessions)
    ? sessionState.sessions
    : [];

  return (
    <Card
      bordered
      bodyStyle={{ padding: 16 }}
      className='mt-4'
      title={t('实时会话')}
    >
      <Space spacing={8} wrap>
        <Tag color={getSummaryColor(sessionState)} shape='circle'>
          {maxSessions > 0
            ? t('活跃会话 {{active}}/{{max}}', {
                active: activeSessions,
                max: maxSessions,
              })
            : t('活跃会话 {{active}}/不限', { active: activeSessions })}
        </Tag>
        <Tag color='white' type='ghost' shape='circle'>
          {t('超时 {{ttl}} 分钟', { ttl: ttlMinutes })}
        </Tag>
        {sessionState?.updated_at ? (
          <Text type='tertiary' size='small'>
            {t('最近刷新')}: {timestamp2string(sessionState.updated_at)}
          </Text>
        ) : null}
      </Space>

      {sessions.length === 0 ? (
        <Empty
          description={t('当前暂无活跃会话')}
          image={null}
          style={{ padding: '24px 0 8px' }}
        />
      ) : (
        <div className='mt-4 border border-[var(--semi-color-border)] rounded-lg overflow-hidden'>
          <div className='grid grid-cols-4 gap-4 px-4 py-3 bg-[var(--semi-color-fill-0)] text-xs font-medium text-[var(--semi-color-text-2)]'>
            <div>{t('会话')}</div>
            <div>{t('最近活跃')}</div>
            <div>{t('过期时间')}</div>
            <div>{t('剩余')}</div>
          </div>
          {sessions.map((session) => (
            <div
              key={`${session.session_id_masked}-${session.last_active_at}`}
              className='grid grid-cols-4 gap-4 px-4 py-3 text-sm border-t border-[var(--semi-color-border)]'
            >
              <div className='font-mono text-xs sm:text-sm'>
                {session.session_id_masked}
              </div>
              <div>{timestamp2string(session.last_active_at)}</div>
              <div>{timestamp2string(session.expires_at)}</div>
              <div>{formatRemaining(session.remaining_seconds, t)}</div>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
};

export default ChannelSessionPanel;
