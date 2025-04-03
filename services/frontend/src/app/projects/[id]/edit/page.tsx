import { getProject } from "@/services/projectService";
import ProjectHeader from "@/components/projects/projectHeader";
import { notFound } from "next/navigation";
import ProjectForm from "@/components/projects/projectForm";
import { editProject } from "./actions";
import { IProject } from "@/models/project";
export default async function ProjectEditPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return <div>
        <ProjectHeader project={project as IProject} currentPage="edit" />
        <div className="container">
            <h1 className="text-2xl font-bold mb-4">Edit Project</h1>
            <ProjectForm action={editProject} project={project as IProject} />
        </div>
    </div>
}