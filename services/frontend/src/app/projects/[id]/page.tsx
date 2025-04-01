import { Metadata } from "next";
import ProjectHeader from "@/components/projects/projectHeader";
import { getProject } from "@/services/projectService";
import { IProject } from "@/models/project";
import { notFound } from "next/navigation";
import ProjectStats from "@/components/projects/projectStats";

export default async function ProjectPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return <div>
        <ProjectHeader project={project} currentPage="" />
    </div>
}