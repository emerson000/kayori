'use server'

import { redirect } from "next/navigation";
import { postTask } from "../../../utils/shared";

export async function createRssFeed(formData: FormData) {
    const title = formData.get('title') as string;
    const urls = formData.getAll('urls[]') as string[];
    const jobId = await postTask('rss', title, { urls });
    redirect(`/collect/${jobId}`);
}