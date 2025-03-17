import { getProjects } from "@/services/projectService";
import Link from "next/link";

export default async function ProjectsPage() {
    const projects = await getProjects();
    return <div className="p-4">
        <div className="flex justify-between">
            <h1 className="text-2xl font-bold">Projects</h1>
            <a href="/projects/new" className="btn">New</a>
        </div>
        <div className="overflow-x-auto">
            <table className="table table-zebra">
                <thead>
                    <tr>
                        <th>Title</th>
                        <th>Created</th>
                        <th>Updated</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {projects.length === 0 ? (
                        <tr>
                            <td colSpan={4} className="text-center">No projects found</td>
                        </tr>
                    ) : (
                        projects.map((project) => (
                            <tr key={project.id}>
                                <td>{project.getDocumentTitle()}</td>
                                <td>{new Date(project.created_at).toLocaleDateString()}</td>
                                <td>{new Date(project.updated_at).toLocaleDateString()}</td>
                                <td>
                                    <Link href={`/projects/${project.id}`} className="btn btn-primary">View</Link>
                                </td>
                            </tr>
                        ))
                    )}
                </tbody>
            </table>
        </div>
    </div>;
}