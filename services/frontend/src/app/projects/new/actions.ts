'use server'

import { createProject } from "@/services/projectService"
import { redirect } from "next/navigation"
import { IProject } from "@/models/project"

export async function newProject(formData) {
    const data: Omit<IProject, 'getDocumentTitle' | 'id' | 'created_at' | 'updated_at'> = {
        title: formData.get('title'),
        description: '',
        number: formData.get('number'),
        status: formData.get('status'),
    };
    const project = await createProject(data)
    if(!project) {
        redirect('/projects/new')
    }
    redirect(`/projects/${project.id}`)
}