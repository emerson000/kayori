export interface IBaseModel {
    id: string;
    created_at: Date;
    updated_at: Date;
}

export class BaseModel implements IBaseModel {
    id: string;
    created_at: Date;
    updated_at: Date;

    constructor({ id, created_at, updated_at }: { id: string, created_at: Date, updated_at: Date }) {
        this.id = id;
        this.created_at = created_at;
        this.updated_at = updated_at;
    }
}
