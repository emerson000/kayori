'use server'

import { redirect } from "next/navigation";
import { postTask } from "../../../utils/shared";

export async function createDeduplicateTask(formData: FormData) {
    const field = formData.get('field') as string;
    const jobs = (formData.getAll('jobs[]') as string[]).filter((job) => job.length > 0);
    await postTask('process', 'deduplicate', "Deduplicate", { jobs, field }, { schedule: false });
    redirect(`/collect`);
}