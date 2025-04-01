import { getProject } from "@/services/projectService"
import { Metadata } from "next"

type Params = Promise<{ id: string }>

export async function generateMetadata({ params }: { params: Params }): Promise<Metadata> {
    const { id } = await params
    const project = await getProject(id)
    return {
        title: project ? `Project - ${project.number}` : 'Projects'
    }
}

export default async function ProjectLayout({ children }: { children: React.ReactNode }) {
    return <div className="container mx-auto px-4">
        {children}
    </div>
}