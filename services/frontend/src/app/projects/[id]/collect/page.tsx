'use client'

import { useEffect, useState } from "react";
import { use } from "react";
import Link from "next/link";
import { getProject } from "@/services/projectService";
import { getJobs } from "@/services/jobService";
import { IJob } from "@/models/job";
import { Project } from "@/models/project";
import { notFound } from "next/navigation";
import ProjectHeader from "@/components/projects/projectHeader";
import CollectTable from "@/components/projects/collectTable";
import InfiniteScroll from "@/components/common/InfiniteScroll";

export default function Page({ params }: { params: Promise<{ id: string }> }) {
    const resolvedParams = use(params);
    const [project, setProject] = useState<any>(null);
    const [initialJobs, setInitialJobs] = useState<IJob[]>([]);

    useEffect(() => {
        const loadInitialData = async () => {
            const projectData = await getProject(resolvedParams.id, true);
            if (!projectData) {
                notFound();
            }
            setProject(projectData);
            const jobs = await getJobs(resolvedParams.id, 1, 10, true);
            setInitialJobs(jobs);
        };
        loadInitialData();
    }, [resolvedParams.id]);

    const loadMoreJobs = async (page: number) => {
        return getJobs(resolvedParams.id, page, 10, true);
    };

    if (!project) return <div>Loading...</div>;

    return <div>
        <ProjectHeader project={new Project(project)} currentPage="collect" />
        <div className="overflow-x-auto">
            <ul className="menu menu-horizontal bg-base-200 float-right">
                <li><Link href={`/projects/${resolvedParams.id}/collect/new`}>New</Link></li>
            </ul>
            <InfiniteScroll
                initialData={initialJobs}
                loadMore={loadMoreJobs}
            >
                {(jobs, loading) => (
                    <CollectTable jobs={jobs} id={resolvedParams.id} loading={loading} />
                )}
            </InfiniteScroll>
        </div>
    </div>
}