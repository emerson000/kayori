export interface IJob {
    id: string;
    title: string;
    service: string;
    status: string;
    projects: string[];
    task: any;
}

export class Job implements IJob {
    id: string;
    title: string;
    service: string;
    status: string;
    projects: string[];
    task: any;

    constructor({ id, title, service, status, task, projects }: { id: string; title: string; service: string; status: string, task: any, projects: string[] }) {
        this.id = id;
        this.title = title;
        this.service = service;
        this.status = status;
        this.projects = projects
        this.task = task;
    }
}

