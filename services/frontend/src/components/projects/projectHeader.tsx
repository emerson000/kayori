import { IProject } from "@/models/project";

const menuItems = [
    {
        label: "Overview",
        href: ''
    },
    {
        label: "Collect",
        href: 'collect'
    },
    {
        label: "Process",
        href: 'process'
    },
    {
        label: "Analyze",
        href: 'analyze'
    },
    {
        label: "Disseminate",
        href: 'disseminate'
    }
]

export default function ProjectHeader({
    project,
    actions,
    currentPage }: {
        project: IProject,
        actions?: React.ReactNode,
        currentPage: string
    }) {

    return <div>
        <div className="flex justify-between items-center outline outline-base-300 rounded-box p-2">
            <h1 className="text-2xl font-bold">
                {project.getDocumentTitle()}
                <label className="badge badge-outline badge-success ml-5">
                    {project.status}
                </label>
            </h1>
            <div className="flex gap-2">
                <a href={`/projects/${project.id}/edit`} className="btn btn-sm btn-base-200">Edit</a>
            </div>
        </div>
        <div className="flex flex-col lg:flex-row justify-between">
            <ul className="menu menu-horizontal bg-base-200 rounded-box my-2 overflow-x-auto w-full lg:w-auto">
                {menuItems.map((item) => (
                    <li key={item.href}>
                        <a href={`/projects/${project.id}/${item.href}`}
                            className={item.href === currentPage ? "menu-active" : ""}
                        >
                            {item.label}
                        </a>
                    </li>
                ))}
            </ul>
            {actions}
        </div>
    </div>
}