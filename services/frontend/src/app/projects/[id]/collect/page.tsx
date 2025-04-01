import Link from "next/link";
import { getProject } from "@/services/projectService";
import { getJobs } from "@/services/jobService";
import { IJob } from "@/models/job";
import { notFound } from "next/navigation";
import ProjectHeader from "@/components/projects/projectHeader";

export const dynamic = 'force-dynamic'

export default async function Page({ params }: { params: Promise<{id: string}> }) {
    const { id } = await params;
    const jobs: IJob[] = await getJobs(id);
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return <div>
        <ProjectHeader project={project} currentPage="collect" />
        <div className="overflow-x-auto">
            <ul className="menu menu-horizontal bg-base-200 float-right">
                <li><Link href={`/projects/${id}/collect/new`}>New</Link></li>
            </ul>
            <table className="table">
                <thead>
                    <tr>
                        <th>Title</th>
                        <th>Status</th>
                        <th>Service</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {jobs.map(job => (
                        <tr key={job.id}>
                            <td>{job.title}</td>
                            <td>
                                <div className={`badge ${job.status === 'pending' ? 'badge-info' : job.status === 'done' ? 'badge-success' : ''}`}>
                                    {job.status ? job.status.toUpperCase() : 'UNKNOWN'}
                                </div>
                            </td>
                            <td>{job.service}</td>
                            <td><Link className="btn" href={`/projects/${id}/collect/${job.id}`}>View</Link></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
}