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
            id: generateUUID(),
            service: service,
            task: task
        })
    });
    return id;
}

export function getApiHostname() {
    return process.env.BACKEND_URL || 'http://localhost:3001';
}