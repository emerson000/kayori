import { getProject } from "@/services/projectService";
import { notFound } from "next/navigation";
import ProjectHeader from "@/components/projects/projectHeader";
import { IProject } from "@/models/project";
export default async function Page({ params }: { params: Promise<{id: string}> }) {
    const { id } = await params;
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return (
        <div>
            <ProjectHeader project={project as IProject} currentPage="analyze" />
            <div className="card bg-base-200 w-96 shadow-sm">
                <div className="card-body">
                    <h2 className="card-title">News Articles</h2>
                    <p>Analyze news articles collected from various sources.</p>
                    <div className="card-actions justify-end">
                        <a href={`/projects/${id}/analyze/news`} className="btn btn-primary">Analyze</a>
                    </div>
                </div>
            </div>
        </div>
    );
}