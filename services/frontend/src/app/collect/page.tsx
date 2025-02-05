import Link from "next/link";

export default function Page() {
    return <div>
        <div className="overflow-x-auto">
            <ul className="menu menu-horizontal bg-base-200 float-right">
                <li><Link href="/collect/new">New</Link></li>
            </ul>
            <table className="table">
                <thead>
                    <tr>
                        <th>Ordered</th>
                        <th>Title</th>
                        <th>Case</th>
                        <th>Status</th>
                        <th>Source</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>02/04/25 13:00</td>
                        <td>General news collection</td>
                        <td>Operation PAPI</td>
                        <td><div className="badge badge-info">In Progress</div></td>
                        <td>RSS</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
}