'use server'

export async function postTask(service: string, title: string, task: object, schedule: object): Promise<string> {
    const path = await getApiHostname() + '/api/jobs';
    const taskBody = {
        title: title,
        service: service,
        status: "pending",
        task: task
    };
    if (schedule['schedule'] === true) {
        taskBody['schedule'] = schedule;
    }
    const result = await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(taskBody)
    });
    if (!result.ok) {
        const errorText = await result.text();
        throw new Error(errorText);
    }
    const resultJson = await result.json();
    return resultJson.id;
}

export async function getJobArtifacts(jobId: string): Promise<object> {
    const path = await getApiHostname() + '/api/jobs/' + jobId + '/artifacts';
    const response = await fetch(path);
    return await response.json();
}

export async function getApiHostname() {
    return process.env.BACKEND_URL || 'http://localhost:3001';
}