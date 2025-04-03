'use server'

import { Artifact, IArtifact } from '../models/artifact';
import { NewsArticle, INewsArticle } from '../models/newsArticle';
import { getApiHostname } from '../utils/shared';
import { BaseService } from './baseService';
import { getJobService } from './jobService';

class ArtifactService extends BaseService {
    public static instance: ArtifactService;
    private projectEndpoints: Map<string, string> = new Map();

    protected constructor() {
        super('');
    }

    public static async getInstance(): Promise<ArtifactService> {
        if (!ArtifactService.instance) {
            ArtifactService.instance = new ArtifactService();
            const apiUrl = await getApiHostname();
            ArtifactService.instance.apiUrl = `${apiUrl}/api`;
            ArtifactService.instance.initialized = true;
        }
        return ArtifactService.instance;
    }

    private getProjectEndpoint(projectId: string): string {
        let endpoint = this.projectEndpoints.get(projectId);
        if (!endpoint) {
            endpoint = `${this.apiUrl}/projects/${projectId}/artifacts`;
            this.projectEndpoints.set(projectId, endpoint);
        }
        return endpoint;
    }

    async getArtifacts(projectId: string, entityType: string, page: number = 1, perPage: number = 10): Promise<IArtifact[] | INewsArticle[]> {
        const endpoint = this.getProjectEndpoint(projectId);
        const data = await this.getAll<any[]>(page, perPage, endpoint, { type: entityType });
        switch (entityType) {
            case 'news_article':
                if (data && data.length > 0) {
                    return data.map(artifact => new NewsArticle(artifact));
                }
                return [];
            default:
                if (data && data.length > 0) {
                    return data.map(artifact => new Artifact(artifact));
                }
                return [];
        }
    }

}

// Update the singleton creation
let artifactService: ArtifactService;
export const getArtifactService = async () => {
    if (!artifactService) {
        artifactService = await ArtifactService.getInstance();
    }
    return artifactService;
};

// Update the exported methods to include projectId
export const getArtifacts = async (projectId: string, entityType: string, page: number = 1, perPage: number = 10, plainObjects: boolean = false) => {
    const service = await getArtifactService();
    if (plainObjects) {
        return service.getArtifacts(projectId, entityType, page, perPage).then(artifacts => artifacts.map((artifact) => ({ ...artifact })));
    }
    return service.getArtifacts(projectId, entityType, page, perPage);
};
