import ProjectHeader from "@/components/projects/projectHeader";
import { getProject } from "@/services/projectService";
import { notFound } from "next/navigation";
import { IProject } from "@/models/project";

export default async function ProjectPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return <div>
        <ProjectHeader project={project as IProject} currentPage="" />
    </div>
}