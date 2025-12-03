import { Card, CardValue, CardDivider, Status } from '../ui/Card';

export interface DNSData {
  server: string;
  testHostname: string;
  forward: {
    result: string;
    time: number; // ms
    status: Status;
  } | null;
  reverse: {
    result: string;
    time: number;
    status: Status;
  } | null;
}

interface DNSCardProps {
  data: DNSData | null;
  loading?: boolean;
}

function formatTime(ms: number): string {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(1)}s`;
  }
  return `${Math.round(ms)}ms`;
}

export function DNSCard({ data, loading }: DNSCardProps) {

  if (loading) {
    return (
      <Card title="DNS" status="loading">
        <CardValue value="Testing..." size="lg" />
      </Card>
    );
  }

  if (!data) {
    return (
      <Card title="DNS" status="unknown">
        <CardValue value="No data" size="md" />
      </Card>
    );
  }

  // Determine overall status based on forward/reverse lookups
  let overallStatus: Status = 'success';
  if (data.forward?.status === 'error' || data.reverse?.status === 'error') {
    overallStatus = 'error';
  } else if (
    data.forward?.status === 'warning' ||
    data.reverse?.status === 'warning'
  ) {
    overallStatus = 'warning';
  }

  return (
    <Card title="DNS" status={overallStatus}>
      <CardValue value={data.server} size="md" />
      <p className="text-xs text-text-muted">Testing: {data.testHostname}</p>
      <CardDivider />
      {data.forward && (
        <div className="mb-2">
          <div className="flex items-center justify-between">
            <span className="text-xs text-text-muted">Forward</span>
            <span
              className={`text-xs font-medium ${
                data.forward.status === 'success'
                  ? 'text-status-success'
                  : data.forward.status === 'warning'
                  ? 'text-status-warning'
                  : 'text-status-error'
              }`}
            >
              {formatTime(data.forward.time)}
            </span>
          </div>
          <p className="text-sm truncate">{data.forward.result}</p>
        </div>
      )}
      {data.reverse && (
        <div>
          <div className="flex items-center justify-between">
            <span className="text-xs text-text-muted">Reverse</span>
            <span
              className={`text-xs font-medium ${
                data.reverse.status === 'success'
                  ? 'text-status-success'
                  : data.reverse.status === 'warning'
                  ? 'text-status-warning'
                  : 'text-status-error'
              }`}
            >
              {formatTime(data.reverse.time)}
            </span>
          </div>
          <p className="text-sm truncate">{data.reverse.result}</p>
        </div>
      )}
    </Card>
  );
}
