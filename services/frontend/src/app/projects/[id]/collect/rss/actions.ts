'use server'

import { redirect } from "next/navigation";
import { postTask } from "@/utils/shared";
import { createJob } from "@/services/jobService";
export async function createRssFeed(formData: FormData, id: string) {
    const title = formData.get('title') as string;
    const urls = (formData.getAll('urls[]') as string[]).filter((url) => url.length > 0);
    const schedule = {
        schedule: formData.get('schedule') === 'true',
        duration: parseInt(formData.get('duration') as string, 10),
        interval: formData.get('interval') as string,
    }
    const jobData = {
        title: title,
        service: 'rss',
        status: 'pending',
        category: 'collect',
        task: { urls },
        schedule: schedule,
        projects: [id],
    }
    const job = await createJob(id, jobData);
    if (!job) {
        throw new Error('Failed to create job');
    }
    redirect(`/projects/${id}/collect/${job.id}`);
}