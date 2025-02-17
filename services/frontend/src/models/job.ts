export interface IJob {
    id: string;
    title: string;
    service: string;
    status: string;
    task: any;
}

export class Job implements IJob {
    id: string;
    title: string;
    service: string;
    status: string;
    task: any;

    constructor({ id, title, service, status, task }: { id: string; title: string; service: string; status: string, task: any }) {
        this.id = id;
        this.title = title;
        this.service = service;
        this.status = status;
        this.task = task;
    }
}

