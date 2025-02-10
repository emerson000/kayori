export interface IJob {
    id: string;
    title: string;
    service: string;
    task: any;
}

export class Job implements IJob {
    id: string;
    title: string;
    service: string;
    task: any;

    constructor({ id, title, service, task }: { id: string; title: string; service: string; task: any }) {
        this.id = id;
        this.title = title;
        this.service = service;
        this.task = task;
    }
}

