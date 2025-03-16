import { Metadata } from "next";
import ProjectHeader from "@/components/projects/projectHeader";
import { getProject } from "@/services/projectService";
import { IProject } from "@/models/project";
import { notFound } from "next/navigation";

export async function generateMetadata({ params }: { params: { id: string } }): Promise<Metadata> {
    const { id } = await params;
    return {
        title: `Case #${id}`
    }
}

export default async function ProjectPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const project = await getProject(id);
    if (!project) {
        notFound();
    }
    return <div>
        <ProjectHeader project={project} />
        <div className="w-full overflow-x-auto pt-2 pb-1 px-1">
            <div className="stats shadow outline outline-base-300 flex flex-row">
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">Critical Alerts</div>
                    <div className="stat-value text-error">1</div>
                </a>
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">Alerts</div>
                    <div className="stat-value text-warning">15</div>
                </a>
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">Saved Artifacts</div>
                    <div className="stat-value">20</div>
                </a>
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">People</div>
                    <div className="stat-value">48</div>
                </a>
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">Organizations</div>
                    <div className="stat-value">12</div>
                </a>
                <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                    <div className="stat-title">News Articles</div>
                    <div className="stat-value">80,100</div>
                </a>
            </div>
        </div>

        <div className="w-full outline outline-base-300 rounded-box mx-1 my-4 p-2">
            <h2 className="text-xl font-bold">Feed</h2>
            <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                    <div className="flex items-start space-x-4">
                        <div className="avatar">
                            <div className="w-12 h-12 rounded-full">
                                <img src="https://i.pravatar.cc/300?u=me" alt="Your avatar" />
                            </div>
                        </div>
                        <div className="flex-1">
                            <textarea
                                className="textarea textarea-bordered w-full"
                                placeholder="Write a new post..."
                                rows={3}
                            ></textarea>
                            <div className="flex justify-between mt-3">
                                <div className="flex space-x-2">
                                    <button className="btn btn-sm">📎 Attach</button>
                                    <button className="btn btn-sm">🏷️ Tag</button>
                                </div>
                                <button className="btn btn-sm btn-primary">Post</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="card bg-base-100 shadow-xl my-4">
                <div className="card-body">
                    <div className="flex items-start space-x-4">
                        <div className="avatar">
                            <div className="w-12 h-12 rounded-full">
                                <img src="https://i.pravatar.cc/300?u=1" alt="User avatar" />
                            </div>
                        </div>
                        <div className="flex-1">
                            <div className="flex justify-between items-center">
                                <h3 className="font-bold">John Doe</h3>
                                <span className="text-sm text-gray-500">2 hours ago</span>
                            </div>
                            <p className="text-sm mt-1">
                                Added a new financial document to the case. The bank statement shows suspicious transactions from offshore accounts.
                            </p>
                            <div className="card bg-base-200 p-3 mt-3">
                                <p className="text-sm font-semibold">📄 Bank_Statement_March2025.pdf</p>
                            </div>
                            <div className="flex justify-between items-center mt-4">
                                <div className="flex space-x-2">
                                    <button className="btn btn-sm btn-ghost">👍 Like</button>
                                    <button className="btn btn-sm btn-ghost">💬 Comment</button>
                                    <button className="btn btn-sm btn-ghost">🔗 Share</button>
                                </div>
                                <div className="badge badge-neutral">Financial Document</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
}