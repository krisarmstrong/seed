import { Card, CardValue, CardRow, CardDivider, Status } from '../ui/Card';

export interface PingResult {
  time: number; // ms
  success: boolean;
}

export interface GatewayData {
  ip: string;
  pings: PingResult[];
  averageLatency: number;
  packetLoss: number; // percentage
}

interface GatewayCardProps {
  data: GatewayData | null;
  loading?: boolean;
  thresholds?: { warning: number; critical: number };
}

function getLatencyStatus(
  value: number,
  thresholds: { warning: number; critical: number }
): Status {
  if (value >= thresholds.critical) return 'error';
  if (value >= thresholds.warning) return 'warning';
  return 'success';
}

function formatTime(ms: number): string {
  return `${Math.round(ms)}ms`;
}

export function GatewayCard({ data, loading, thresholds }: GatewayCardProps) {
  const t = thresholds || { warning: 50, critical: 200 };

  if (loading) {
    return (
      <Card title="Gateway" status="loading">
        <CardValue value="Pinging..." size="lg" />
      </Card>
    );
  }

  if (!data) {
    return (
      <Card title="Gateway" status="unknown">
        <CardValue value="No gateway" size="md" />
      </Card>
    );
  }

  // Determine status based on packet loss and latency
  let status: Status = 'success';
  if (data.packetLoss === 100) {
    status = 'error';
  } else if (data.packetLoss > 0) {
    status = 'warning';
  } else {
    status = getLatencyStatus(data.averageLatency, t);
  }

  return (
    <Card title="Gateway" status={status}>
      <CardValue value={data.ip} size="lg" />
      <CardDivider />
      <div className="flex items-center gap-2 mb-2">
        {data.pings.map((ping, index) => (
          <div
            key={index}
            className={`flex-1 text-center py-1 rounded text-sm ${
              ping.success
                ? getLatencyStatus(ping.time, t) === 'success'
                  ? 'bg-status-success/20 text-status-success'
                  : getLatencyStatus(ping.time, t) === 'warning'
                  ? 'bg-status-warning/20 text-status-warning'
                  : 'bg-status-error/20 text-status-error'
                : 'bg-status-error/20 text-status-error'
            }`}
          >
            {ping.success ? formatTime(ping.time) : '✕'}
          </div>
        ))}
      </div>
      <CardRow
        label="Avg Latency"
        value={formatTime(data.averageLatency)}
        status={status}
      />
      {data.packetLoss > 0 && (
        <CardRow
          label="Packet Loss"
          value={`${data.packetLoss}%`}
          status={data.packetLoss === 100 ? 'error' : 'warning'}
        />
      )}
    </Card>
  );
}
