'use server'

import { redirect } from "next/navigation";

export async function searchNews(formData: FormData) {
    const search = formData.get('search') as string;
    redirect(`/analyze/news?search=${search}`);
}