'use server'

import { redirect } from "next/navigation";


export async function createRssFeed(formData: FormData) {
    const title = formData.get('title') as string;
    const urls = formData.getAll('urls[]') as string[];
    
    console.log({ title, urls });
    fetch('/api/rss');

    redirect('/collect/rss');
}