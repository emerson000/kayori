import Link from "next/link";

export default function CollectTable({ jobs, id, loading }: { jobs: any[], id: string, loading: boolean }) {
    return <table className="table table-pin-rows">
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
            {loading && (
                <tr>
                    <td colSpan={4} className="text-center">
                        <div className="loading loading-spinner loading-md"></div>
                    </td>
                </tr>
            )}
        </tbody>
    </table>
}