import { v7 as generateUUID } from 'uuid';

export async function postTask(service: string, title: string, task: object) {
    const path = getApiHostname() + '/api/task';
    console.log(path);
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
}

export function getApiHostname() {
    return process.env.BACKEND_URL || 'http://localhost:3001';
}