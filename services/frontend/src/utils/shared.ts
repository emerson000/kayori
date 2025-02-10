import { v7 as generateUUID } from 'uuid';

export async function postTask(service: string, title: string, task: object): Promise<string> {
    const path = getApiHostname() + '/api/task';
    const id = generateUUID();
    await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: id,
            title: title,
            service: service,
            task: task
        })
    });
    return id;
}

export async function getJobArtifacts(jobId: string): Promise<object> {
    const path = getApiHostname() + '/api/jobs/'+jobId +'/artifacts';
    const response = await fetch(path);
    return await response.json();
}

export function getApiHostname() {
    return process.env.BACKEND_URL || 'http://localhost:3001';
}