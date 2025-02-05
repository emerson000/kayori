import Link from "next/link";

const sources = [
  {
    name: 'RSS',
    description: 'Ingest new articles from RSS feeds.',
    slug: 'rss'
  }
]

export default function Page() {
  return (
    <div className="flex flex-wrap">
      {sources.map((source, index) => (
        <div className="flex-1" key={index}>
          <div className="card card-border bg-base-100 w-96 hover:bg-base-200">
            <div className="card-body">
              <h2 className="card-title">{source.name}</h2>
              <p>{source.description}</p>
              <div className="card-actions justify-end">
                <Link className="btn btn-primary" href="/collect/rss">Start</Link>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}