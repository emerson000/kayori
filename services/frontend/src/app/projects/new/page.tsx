import { newProject } from "./actions";
import ProjectForm from "@/components/projects/projectForm";

export default async function ProjectNewPage() {
    return <div className="container mx-auto max-w-4xl px-5">
        <h1 className="text-2xl font-bold mb-4">New Project</h1>
        <ProjectForm action={newProject}/>
    </div>
}