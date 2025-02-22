'use server'

import { redirect } from "next/navigation";
import { postTask } from "../../../utils/shared";

export async function createRssFeed(formData: FormData) {
    const title = formData.get('title') as string;
    const urls = (formData.getAll('urls[]') as string[]).filter((url) => url.length > 0);
    const schedule = {
        schedule: formData.get('schedule') === 'true',
        duration: parseInt(formData.get('duration') as string, 10),
        interval: formData.get('interval') as string,
    }
    const jobId = await postTask('rss', title, { urls }, schedule);
    redirect(`/collect/${jobId}`);
}