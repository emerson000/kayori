export interface INewsArticle {
    id: string;
    title: string;
    description: string;
    url: string;
    published: string;
    timestamp: number;
    author: string;
    categories: string[];
    service: string;
    service_id: string;
    checksum: string;
    job_id: string;
}

export class NewsArticle implements INewsArticle {
    id: string;
    title: string;
    description: string;
    url: string;
    published: string;
    timestamp: number;
    author: string;
    categories: string[];
    service: string;
    service_id: string;
    checksum: string;
    job_id: string;

    constructor({
        id,
        title,
        description,
        url,
        published,
        timestamp,
        author,
        categories,
        service,
        service_id,
        checksum,
        job_id
    }: INewsArticle) {
        this.id = id;
        this.title = title;
        this.description = description;
        this.url = url;
        this.published = published;
        this.timestamp = timestamp;
        this.author = author;
        this.categories = categories;
        this.service = service;
        this.service_id = service_id;
        this.checksum = checksum;
        this.job_id = job_id;
    }
}