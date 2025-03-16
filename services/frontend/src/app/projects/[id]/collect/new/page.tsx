import Link from "next/link";
import { getProject } from "@/services/projectService";
import { notFound } from "next/navigation";
import ProjectHeader from "@/components/projects/projectHeader";
const sources = [
  {
    name: 'RSS',
    description: 'Ingest new articles from RSS feeds.',
    slug: 'rss'
  }
]

export default async function Page({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const project = await getProject(id);
  if (!project) {
    notFound();
  }
  return (<div>
    <ProjectHeader project={project} currentPage="collect" />
    <div className="flex flex-wrap">
      {sources.map((source, index) => (
        <div className="flex-1" key={index}>
          <div className="card card-border bg-base-100 w-96 hover:bg-base-200">
            <div className="card-body">
              <h2 className="card-title">{source.name}</h2>
              <p>{source.description}</p>
              <div className="card-actions justify-end">
                <Link className="btn btn-primary" href={`/projects/${id}/collect/${source.slug}`}>Start</Link>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  </div>
  );
}