import { Card, CardValue, CardRow, CardDivider, Status } from '../ui/Card';

export interface LinkData {
  linkUp: boolean;
  speed: string;
  duplex: string;
  advertisedSpeeds: string[];
}

interface LinkCardProps {
  data: LinkData | null;
  loading?: boolean;
}

export function LinkCard({ data, loading }: LinkCardProps) {
  if (loading || !data) {
    return (
      <Card title="Link" status="loading">
        <CardValue value="..." size="lg" />
      </Card>
    );
  }

  const status: Status = data.linkUp ? 'success' : 'error';

  return (
    <Card title="Link" status={status}>
      <CardValue
        value={data.linkUp ? data.speed || 'Connected' : 'Down'}
        size="lg"
        status={status}
      />
      {data.linkUp && (
        <>
          <CardDivider />
          <CardRow label="Duplex" value={data.duplex || 'Unknown'} />
          {data.advertisedSpeeds && data.advertisedSpeeds.length > 0 && (
            <div className="mt-2">
              <p className="text-xs text-text-muted mb-1">Advertised</p>
              <div className="flex flex-wrap gap-1">
                {data.advertisedSpeeds.map((speed) => (
                  <span
                    key={speed}
                    className="text-xs px-2 py-0.5 bg-surface-hover rounded"
                  >
                    {speed}
                  </span>
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </Card>
  );
}
