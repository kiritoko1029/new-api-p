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

import React, { useState } from 'react';
import {
  Modal,
  Button,
  Descriptions,
  Tag,
  Spin,
  Space,
  Typography,
} from '@douyinfe/semi-ui';
import { API } from '../../../../helpers/api';

const { Title, Text } = Typography;

const CodingPlanUsageView = ({ t, record, initialUsage }) => {
  const [loading, setLoading] = useState(false);
  const [usage, setUsage] = useState(initialUsage);

  const parseUsage = () => {
    if (usage) return usage;
    if (record?.coding_plan_usage) {
      try {
        return JSON.parse(record.coding_plan_usage);
      } catch {
        return null;
      }
    }
    return null;
  };

  const handleRefresh = async () => {
    setLoading(true);
    try {
      const res = await API.post(
        `/api/channel/${record.id}/coding_plan/refresh`,
      );
      const { success, data, message } = res.data;
      if (success) {
        setUsage(
          typeof data === 'string' ? JSON.parse(data) : data,
        );
      } else {
        Modal.error({ title: t('错误'), content: message });
      }
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  const data = parseUsage();

  return (
    <div>
      {loading && (
        <Spin style={{ display: 'block', textAlign: 'center', margin: '20px 0' }} />
      )}
      {data && (
        <Space vertical style={{ width: '100%' }}>
          {data.error ? (
            <Tag color='red' size='large'>
              {data.error}
            </Tag>
          ) : (
            <>
              <Descriptions row>
                <Descriptions.Item itemKey={t('套餐等级')}>
                  <Tag color='blue' size='large'>
                    {data.level || '-'}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item itemKey={t('查询时间')}>
                  {data.query_time
                    ? new Date(data.query_time).toLocaleString()
                    : '-'}
                </Descriptions.Item>
              </Descriptions>
              {data.limits?.map((limit, idx) => (
                <div
                  key={idx}
                  style={{
                    background: 'var(--semi-color-fill-0)',
                    borderRadius: 8,
                    padding: 16,
                  }}
                >
                  <Title heading={6}>
                    {limit.type === 'TIME_LIMIT'
                      ? t('时间限制')
                      : limit.type}
                  </Title>
                  <Descriptions row>
                    <Descriptions.Item itemKey={t('限额')}>
                      {limit.usage} {t('次')}
                    </Descriptions.Item>
                    <Descriptions.Item itemKey={t('已使用')}>
                      {limit.current_value} {t('次')}
                    </Descriptions.Item>
                    <Descriptions.Item itemKey={t('剩余')}>
                      {limit.remaining} {t('次')}
                    </Descriptions.Item>
                    <Descriptions.Item itemKey={t('使用率')}>
                      <Tag
                        color={
                          limit.percentage > 80
                            ? 'red'
                            : limit.percentage > 50
                              ? 'orange'
                              : 'green'
                        }
                      >
                        {limit.percentage}%
                      </Tag>
                    </Descriptions.Item>
                    <Descriptions.Item itemKey={t('重置时间')}>
                      {limit.next_reset_time
                        ? new Date(limit.next_reset_time).toLocaleString()
                        : '-'}
                    </Descriptions.Item>
                  </Descriptions>
                  {limit.usage_details?.length > 0 && (
                    <div style={{ marginTop: 8 }}>
                      <Text type='tertiary' size='small'>
                        {t('使用详情')}：
                      </Text>
                      <div
                        style={{
                          marginTop: 4,
                          display: 'flex',
                          gap: 8,
                          flexWrap: 'wrap',
                        }}
                      >
                        {limit.usage_details.map((d, i) => (
                          <Tag key={i}>
                            {d.model_code}: {d.usage}
                          </Tag>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </>
          )}
        </Space>
      )}
      <div style={{ marginTop: 16, textAlign: 'right' }}>
        <Button onClick={handleRefresh} loading={loading}>
          {t('刷新用量')}
        </Button>
      </div>
    </div>
  );
};

export const openCodingPlanUsageModal = ({ t, record }) => {
  const tt = typeof t === 'function' ? t : (v) => v;

  Modal.info({
    title:
      (record?.name || '') +
      ' - ' +
      (record?.type === 58
        ? tt('智谱编程套餐')
        : tt('智谱编程套餐(国际版)')),
    centered: true,
    width: 700,
    style: { maxWidth: '95vw' },
    content: <CodingPlanUsageView t={tt} record={record} />,
    footer: (
      <div className='flex justify-end gap-2'>
        <Button type='primary' theme='solid' onClick={() => Modal.destroyAll()}>
          {tt('关闭')}
        </Button>
      </div>
    ),
  });
};
