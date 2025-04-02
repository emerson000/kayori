'use server'

import { Job, IJob } from '../models/job';
import { getApiHostname } from '../utils/shared';
import { BaseService } from './baseService';

class JobService extends BaseService {
    public static instance: JobService;
    private projectEndpoints: Map<string, string> = new Map();

    protected constructor() {
        super('');
    }

    public static async getInstance(): Promise<JobService> {
        if (!JobService.instance) {
            JobService.instance = new JobService();
            const apiUrl = await getApiHostname();
            JobService.instance.apiUrl = `${apiUrl}/api`;
            JobService.instance.initialized = true;
        }
        return JobService.instance;
    }

    private getProjectEndpoint(projectId: string): string {
        let endpoint = this.projectEndpoints.get(projectId);
        if (!endpoint) {
            endpoint = `${this.apiUrl}/projects/${projectId}/jobs`;
            this.projectEndpoints.set(projectId, endpoint);
        }
        return endpoint;
    }

    async getJobs(projectId: string, page: number = 1, perPage: number = 10): Promise<IJob[]> {
        const data = await this.getAll<any[]>(page, perPage, this.getProjectEndpoint(projectId));
        return data.map(job => new Job(job));
    }

    async getJob(projectId: string, id: string): Promise<IJob | null> {
        const data = await this.get<any>(id, this.getProjectEndpoint(projectId));
        return data ? new Job(data) : null;
    }

    async createJob(projectId: string, projectData: Omit<IJob, 'id' | 'created_at' | 'updated_at'>): Promise<IJob | null> {
        const data = await this.post<any>(projectData, this.getProjectEndpoint(projectId));
        return data ? new Job(data) : null;
    }

    async updateJob(projectId: string, id: string, jobData: Partial<Omit<IJob, 'id' | 'created_at' | 'updated_at'>>): Promise<IJob | null> {
        const data = await this.put<any>(id, jobData, this.getProjectEndpoint(projectId));
        return data ? new Job(data) : null;
    }

    async deleteJob(projectId: string, id: string): Promise<boolean> {
        return await this.delete(id, this.getProjectEndpoint(projectId));
    }

    async getJobArtifacts(projectId: string, id: string, page = 1, limit = 10) : Promise<object[] | null> {
        const url = `${this.getProjectEndpoint(projectId)}/${id}/artifacts`
        const data = await this.getAll<any[]>(page, limit, url);
        return data;
    }

}

// Update the singleton creation
let jobService: JobService;
export const getJobService = async () => {
    if (!jobService) {
        jobService = await JobService.getInstance();
    }
    return jobService;
};

// Update the exported methods to include projectId
export const getJobs = async (projectId: string, page: number = 1, perPage: number = 10, plainObjects: boolean = false) => {
    const service = await getJobService();
    const jobs = await service.getJobs(projectId, page, perPage);
    return plainObjects ? jobs.map(job => ({ ...job })) : jobs;
};

export const getJob = async (projectId: string, id: string) => {
    const service = await getJobService();
    return service.getJob(projectId, id);
};

export const createJob = async (projectId: string, jobData: Omit<IJob, 'id' | 'created_at' | 'updated_at'>) => {
    const service = await getJobService();
    return service.createJob(projectId, jobData);
};

export const updateJob = async (projectId: string, id: string, jobData: Partial<Omit<IJob, 'id' | 'created_at' | 'updated_at'>>) => {
    const service = await getJobService();
    return service.updateJob(projectId, id, jobData);
};

export const deleteJob = async (projectId: string, id: string) => {
    const service = await getJobService();
    return service.deleteJob(projectId, id);
};

export const getJobArtifacts = async (projectId: string, id: string) => {
    const service = await getJobService()
    return service.getJobArtifacts(projectId, id)
}