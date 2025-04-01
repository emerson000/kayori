import { IBaseModel, BaseModel } from "./baseModel";

export interface INewsArticle extends IBaseModel {
    title: string;
    description: string;
    entity_type: string;
    url: string;
    published: string;
    timestamp: number;
    author: string;
    categories: string[];
    service: string;
    service_id: string;
    checksum: string;
    job_id: string;
    cluster_id?: string;
    cluster_articles?: INewsArticle[];
}

export class NewsArticle extends BaseModel implements INewsArticle {
    title: string;
    description: string;
    url: string;
    entity_type: string;
    published: string;
    timestamp: number;
    author: string;
    categories: string[];
    service: string;
    service_id: string;
    checksum: string;
    job_id: string;
    cluster_id?: string;
    cluster_articles?: INewsArticle[];

    constructor({
        id,
        created_at,
        updated_at,
        title,
        description,
        entity_type,
        url,
        published,
        timestamp,
        author,
        categories,
        service,
        service_id,
        checksum,
        job_id,
        cluster_id,
        cluster_articles
    }: INewsArticle) {
        super({ id, created_at, updated_at });
        this.title = title;
        this.description = description;
        this.entity_type = entity_type;
        this.url = url;
        this.published = published;
        this.timestamp = timestamp;
        this.author = author;
        this.categories = categories;
        this.service = service;
        this.service_id = service_id;
        this.checksum = checksum;
        this.job_id = job_id;
        this.cluster_id = cluster_id;
        this.cluster_articles = cluster_articles;
    }

    getSecondLevelDomain(): string {
        const url = new URL(this.url);
        const host = url.hostname;
        const parts = host.split('.');
        if (parts.length > 2) {
            return parts.slice(parts.length - 2).join('.');
        }
        return host;
    }
}