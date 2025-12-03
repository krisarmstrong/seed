import { Card, CardValue, CardRow, CardDivider, Status } from '../ui/Card';

export interface DHCPTiming {
  discover: number; // ms
  offer: number;
  request: number;
  ack: number;
  total: number;
}

export interface DHCPData {
  mode: 'dhcp' | 'static';
  ip: string | null;
  subnet: string | null;
  gateway: string | null;
  dns: string[];
  server: string | null;
  leaseTime: number | null;
  timing: DHCPTiming | null;
}

interface DHCPCardProps {
  data: DHCPData | null;
  loading?: boolean;
  thresholds?: {
    total: { warning: number; critical: number };
    perPhase: { warning: number; critical: number };
  };
}

function getTimingStatus(
  value: number,
  thresholds: { warning: number; critical: number }
): Status {
  if (value >= thresholds.critical) return 'error';
  if (value >= thresholds.warning) return 'warning';
  return 'success';
}

function formatTime(ms: number): string {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(1)}s`;
  }
  return `${Math.round(ms)}ms`;
}

export function DHCPCard({ data, loading, thresholds }: DHCPCardProps) {
  const defaultThresholds = {
    total: { warning: 500, critical: 2000 },
    perPhase: { warning: 200, critical: 1000 },
  };
  const t = thresholds || defaultThresholds;

  if (loading) {
    return (
      <Card title="DHCP" status="loading">
        <CardValue value="Acquiring..." size="lg" />
      </Card>
    );
  }

  if (!data) {
    return (
      <Card title="DHCP" status="unknown">
        <CardValue value="No data" size="md" />
      </Card>
    );
  }

  if (data.mode === 'static') {
    return (
      <Card title="IP (Static)" status="success">
        <CardValue value={data.ip || 'Not configured'} size="lg" />
        {data.ip && (
          <>
            <CardDivider />
            {data.subnet && <CardRow label="Subnet" value={data.subnet} />}
            {data.gateway && <CardRow label="Gateway" value={data.gateway} />}
            {data.dns.length > 0 && (
              <CardRow label="DNS" value={data.dns.join(', ')} />
            )}
          </>
        )}
      </Card>
    );
  }

  const totalStatus = data.timing
    ? getTimingStatus(data.timing.total, t.total)
    : 'unknown';

  return (
    <Card title="DHCP" status={data.ip ? totalStatus : 'error'}>
      <CardValue
        value={data.ip || 'Failed'}
        size="lg"
        status={data.ip ? undefined : 'error'}
      />
      {data.ip && (
        <>
          <CardDivider />
          {data.server && <CardRow label="Server" value={data.server} />}
          {data.leaseTime && (
            <CardRow label="Lease" value={`${data.leaseTime}s`} />
          )}
          {data.timing && (
            <>
              <CardDivider />
              <p className="text-xs text-text-muted mb-2">Timing</p>
              <div className="space-y-1">
                <CardRow
                  label="Discover → Offer"
                  value={formatTime(data.timing.discover)}
                  status={getTimingStatus(data.timing.discover, t.perPhase)}
                />
                <CardRow
                  label="Offer → Request"
                  value={formatTime(data.timing.offer)}
                  status={getTimingStatus(data.timing.offer, t.perPhase)}
                />
                <CardRow
                  label="Request → Ack"
                  value={formatTime(data.timing.request)}
                  status={getTimingStatus(data.timing.request, t.perPhase)}
                />
                <div className="pt-1 border-t border-surface-border">
                  <CardRow
                    label="Total"
                    value={formatTime(data.timing.total)}
                    status={totalStatus}
                  />
                </div>
              </div>
            </>
          )}
        </>
      )}
    </Card>
  );
}
