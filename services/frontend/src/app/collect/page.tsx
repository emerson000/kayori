import Link from "next/link";
import { getJobs } from "../../services/jobService";
import { IJob } from "../../models/job";

export const dynamic = 'force-dynamic'

export default async function Page() {
    const jobs: IJob[] = await getJobs('collect');
    return <div>
        <div className="overflow-x-auto">
            <ul className="menu menu-horizontal bg-base-200 float-right">
                <li><Link href="/collect/new">New</Link></li>
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
                            <td><Link className="btn" href={`/collect/${job.id}`}>View</Link></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
}