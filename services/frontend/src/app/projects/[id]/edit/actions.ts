'use server'
import { getProject, updateProject } from "@/services/projectService";
import { redirect } from "next/navigation";

export async function editProject(formData: FormData) {
    const project = await getProject(formData.get('id') as string);
    if (!project || !project.id) {
        throw new Error('Project not found');
    }
    const number = formData.get('number') as string;
    const title = formData.get('title') as string;
    const status = formData.get('status') as string;
    const projectData = {
        number,
        title,
        status
    }
    await updateProject(project.id, projectData);
    redirect(`/projects/${project.id}`)
}