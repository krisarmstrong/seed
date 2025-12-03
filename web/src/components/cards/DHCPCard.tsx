import { Card, CardValue, CardRow, CardDivider, Status } from '../ui/Card';

export interface DHCPTiming {
  discover: number; // ms
  offer: number;
  request: number;
  ack: number;
  total: number;
}

export interface IPv4Info {
  address: string;
  subnet: string;
  gateway: string | null;
  dhcpServer: string | null;
  leaseTime: number | null;
}

export interface IPv6Info {
  address: string;
  prefix: number;
  scope: 'global' | 'link-local' | 'unique-local';
  source: 'slaac' | 'dhcpv6' | 'static' | 'temporary';
}

export interface DHCPData {
  mac: string;
  mode: 'dhcp' | 'static' | 'auto';
  ipv4: IPv4Info | null;
  ipv6: IPv6Info[];
  dns: string[];
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

function formatLeaseTime(seconds: number): string {
  if (seconds >= 86400) {
    const days = Math.floor(seconds / 86400);
    return `${days}d`;
  }
  if (seconds >= 3600) {
    const hours = Math.floor(seconds / 3600);
    return `${hours}h`;
  }
  if (seconds >= 60) {
    const mins = Math.floor(seconds / 60);
    return `${mins}m`;
  }
  return `${seconds}s`;
}

function getScopeLabel(scope: IPv6Info['scope']): string {
  switch (scope) {
    case 'global': return 'Global';
    case 'link-local': return 'Link-Local';
    case 'unique-local': return 'ULA';
    default: return scope;
  }
}

function getSourceLabel(source: IPv6Info['source']): string {
  switch (source) {
    case 'slaac': return 'SLAAC';
    case 'dhcpv6': return 'DHCPv6';
    case 'static': return 'Static';
    case 'temporary': return 'Temporary';
    default: return source;
  }
}

export function DHCPCard({ data, loading, thresholds }: DHCPCardProps) {
  const defaultThresholds = {
    total: { warning: 500, critical: 2000 },
    perPhase: { warning: 200, critical: 1000 },
  };
  const t = thresholds || defaultThresholds;

  if (loading) {
    return (
      <Card title="IP Config" status="loading">
        <CardValue value="Loading..." size="lg" />
      </Card>
    );
  }

  if (!data) {
    return (
      <Card title="IP Config" status="unknown">
        <CardValue value="No data" size="md" />
      </Card>
    );
  }

  const hasIPv4 = data.ipv4 !== null;
  const hasIPv6 = data.ipv6.length > 0;
  const globalIPv6 = data.ipv6.filter(ip => ip.scope === 'global');

  // Determine overall status
  const status: Status = hasIPv4 || globalIPv6.length > 0 ? 'success' : 'warning';

  // Primary display value
  const primaryIP = data.ipv4?.address || globalIPv6[0]?.address || 'No IP';

  return (
    <Card title="IP Config" status={status}>
      <CardValue value={primaryIP} size="lg" />

      <CardDivider />

      {/* MAC Address */}
      <CardRow label="MAC" value={data.mac} />
      <CardRow label="Mode" value={data.mode.toUpperCase()} />

      {/* IPv4 Section */}
      {hasIPv4 && data.ipv4 && (
        <>
          <CardDivider />
          <p className="text-xs text-text-muted mb-1 font-medium">IPv4</p>
          <CardRow label="Address" value={`${data.ipv4.address}/${data.ipv4.subnet}`} />
          {data.ipv4.gateway && <CardRow label="Gateway" value={data.ipv4.gateway} />}
          {data.ipv4.dhcpServer && <CardRow label="DHCP Server" value={data.ipv4.dhcpServer} />}
          {data.ipv4.leaseTime && (
            <CardRow label="Lease" value={formatLeaseTime(data.ipv4.leaseTime)} />
          )}
        </>
      )}

      {/* IPv6 Section */}
      {hasIPv6 && (
        <>
          <CardDivider />
          <p className="text-xs text-text-muted mb-1 font-medium">IPv6</p>
          <div className="space-y-2">
            {data.ipv6.map((ip, idx) => (
              <div key={idx} className="text-xs">
                <div className="font-mono text-text-primary break-all">
                  {ip.address}/{ip.prefix}
                </div>
                <div className="flex gap-2 mt-0.5 text-text-muted">
                  <span className={ip.scope === 'global' ? 'text-status-success' : ''}>
                    {getScopeLabel(ip.scope)}
                  </span>
                  <span>•</span>
                  <span>{getSourceLabel(ip.source)}</span>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {/* DNS Section */}
      {data.dns.length > 0 && (
        <>
          <CardDivider />
          <p className="text-xs text-text-muted mb-1 font-medium">DNS Servers</p>
          <div className="flex flex-col gap-0.5">
            {data.dns.map((server, idx) => (
              <span key={idx} className="text-xs font-mono text-text-secondary">
                {server}
              </span>
            ))}
          </div>
        </>
      )}

      {/* DHCP Timing (if available) */}
      {data.timing && (
        <>
          <CardDivider />
          <p className="text-xs text-text-muted mb-2 font-medium">DHCP Timing</p>
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
                status={getTimingStatus(data.timing.total, t.total)}
              />
            </div>
          </div>
        </>
      )}
    </Card>
  );
}
