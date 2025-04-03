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
            <ProjectHeader project={project as IProject} currentPage="process" />
            <div className="card bg-base-200 w-96 shadow-sm">
                <div className="card-body">
                    <h2 className="card-title">Deduplication</h2>
                    <p>Deduplicate collected artifacts.</p>
                    <div className="card-actions justify-end">
                        <a href="/process/deduplicate" className="btn btn-primary">Start</a>
                    </div>
                </div>
            </div>
        </div>
    );
}