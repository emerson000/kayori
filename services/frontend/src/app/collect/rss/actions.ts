'use server'

import { redirect } from "next/navigation";


export async function createRssFeed(formData: FormData) {
    const title = formData.get('title') as string;
    const urls = formData.getAll('urls[]') as string[];
    
    console.log({ title, urls });
    await fetch('/api/task', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: "12345",
            service: "rss",
            task: {
                urls: urls
            }
        })
    });

    redirect('/collect/rss');
}