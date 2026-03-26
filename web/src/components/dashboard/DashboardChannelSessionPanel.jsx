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
import { Card, Empty, Space, Tag, Typography, Spin } from '@douyinfe/semi-ui';
import { timestamp2string } from '../../helpers';
import { useClaudeChannelSessionStream } from '../../hooks/channels/useClaudeChannelSessionStream';
import { getChannelIcon } from '../../helpers';

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

const getChannelStatusColor = (state) => {
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

const DashboardChannelSessionPanel = ({ CARD_PROPS, t }) => {
  // Get all Claude channels (no specific channelIds means all Claude channels)
  const { sessionStateMap, sessionState } = useClaudeChannelSessionStream({
    channelIds: [],
    details: true,
    enabled: true,
  });

  const channels = Object.values(sessionStateMap || {});
  const hasChannels = channels.length > 0;

  if (!hasChannels) {
    return null;
  }

  return (
    <Card
      {...CARD_PROPS}
      className='w-full'
      title={
        <Space spacing={6} align='center'>
          <span>{t('渠道会话状态')}</span>
          <Text type='tertiary' size='small'>
            {t('Claude 渠道实时会话信息')}
          </Text>
        </Space>
      }
    >
      <div className='space-y-4'>
        {channels.map((channelState) => {
          const activeSessions = Number(channelState?.active_sessions) || 0;
          const maxSessions = Number(channelState?.max_sessions) || 0;
          const ttlMinutes = Number(channelState?.ttl_minutes) || 30;
          const sessions = Array.isArray(channelState?.sessions)
            ? channelState.sessions
            : [];
          const channelType = Number(channelState?.channel_type);

          return (
            <div
              key={channelState.channel_id}
              className='border border-[var(--semi-color-border)] rounded-lg p-4'
            >
              {/* Channel Header */}
              <div className='flex items-center justify-between mb-3'>
                <Space spacing={8} align='center'>
                  <span className='text-lg font-medium'>
                    {channelState.channel_name || `Channel ${channelState.channel_id}`}
                  </span>
                  {getChannelIcon(channelType)}
                </Space>
                <Space spacing={8} wrap>
                  <Tag color={getChannelStatusColor(channelState)} shape='circle'>
                    {maxSessions > 0
                      ? t('活跃 {{active}}/{{max}}', {
                          active: activeSessions,
                          max: maxSessions,
                        })
                      : t('活跃 {{active}}/不限', { active: activeSessions })}
                  </Tag>
                  <Tag color='white' type='ghost' shape='circle'>
                    {t('超时 {{ttl}} 分钟', { ttl: ttlMinutes })}
                  </Tag>
                </Space>
              </div>

              {/* Sessions Table */}
              {sessions.length === 0 ? (
                <Empty
                  description={t('当前暂无活跃会话')}
                  style={{ padding: '16px 0 8px' }}
                />
              ) : (
                <div className='border border-[var(--semi-color-border)] rounded-lg overflow-hidden'>
                  <div className='grid grid-cols-5 gap-4 px-4 py-3 bg-[var(--semi-color-fill-0)] text-xs font-medium text-[var(--semi-color-text-2)]'>
                    <div>{t('用户')}</div>
                    <div>{t('会话')}</div>
                    <div>{t('最近活跃')}</div>
                    <div>{t('过期时间')}</div>
                    <div>{t('剩余')}</div>
                  </div>
                  {sessions.map((session, idx) => (
                    <div
                      key={`${session.session_id_masked}-${session.last_active_at}-${idx}`}
                      className='grid grid-cols-5 gap-4 px-4 py-3 text-sm border-t border-[var(--semi-color-border)]'
                    >
                      <div className='font-medium text-[var(--semi-color-text-1)]'>
                        {session.masked_username || '-'}
                      </div>
                      <div className='font-mono text-xs text-[var(--semi-color-text-2)]'>
                        {session.session_id_masked}
                      </div>
                      <div className='text-[var(--semi-color-text-2)]'>
                        {timestamp2string(session.last_active_at)}
                      </div>
                      <div className='text-[var(--semi-color-text-2)]'>
                        {timestamp2string(session.expires_at)}
                      </div>
                      <div className='text-[var(--semi-color-text-2)]'>
                        {formatRemaining(session.remaining_seconds, t)}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </Card>
  );
};

export default DashboardChannelSessionPanel;
