import { IBaseModel, BaseModel } from "./baseModel";

export const ProjectStatus = {
    pending_start: { label: 'Pending Start', className: 'ghost' },
    in_progress: { label: 'In Progress', className: 'success' },
    closed: { label: 'Closed', className: 'ghost', }
}

export interface IProject extends IBaseModel {
    title: string;
    description: string;
    number: string;
    status: string;
    getDocumentTitle(): string;
}

export class Project extends BaseModel implements IProject {
    title: string;
    description: string;
    number: string;
    status: string;
    constructor({ id,
        created_at,
        updated_at,
        title,
        description,
        number,
        status,
    }: {
        id: string,
        created_at: Date,
        updated_at: Date,
        title: string,
        description: string,
        number: string,
        status: string
    }) {
        super({ id, created_at, updated_at });
        this.title = title;
        this.description = description;
        this.number = number;
        this.status = status;
    }

    getDocumentTitle(): string {
        return `${this.number}: ${this.title}`;
    }
}