import { IProject } from "@/models/project";

export default function ProjectHeader({ project, actions }: { project: IProject, actions?: React.ReactNode }) {
    return <div>
        <h1 className="text-2xl font-bold outline outline-base-300 rounded-box p-2">
            {project.getDocumentTitle()}
            <label className="badge badge-outline badge-success ml-5">
                {project.status}
            </label>
        </h1>
        <div className="flex flex-col lg:flex-row justify-between">
            <ul className="menu menu-horizontal bg-base-200 rounded-box my-2 overflow-x-auto w-full lg:w-auto">
                <li><a href={`/projects/${project.id}`}>Overview</a></li>
                <li><a href={`/projects/${project.id}/collect`}>Collect</a></li>
                <li><a href={`/projects/${project.id}/process`}>Process</a></li>
                <li><a href={`/projects/${project.id}/analyze`}>Analyze</a></li>
                <li><a href={`/projects/${project.id}/disseminate`}>Disseminate</a></li>
            </ul>
            {actions}
        </div>
    </div>
}