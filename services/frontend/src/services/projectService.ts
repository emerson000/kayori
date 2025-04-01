'use server'

import { Project, IProject } from '../models/project';
import { getApiHostname } from '../utils/shared';
import { BaseService } from './baseService';

class ProjectService extends BaseService {
    public static instance: ProjectService;

    protected constructor() {
        super('');
    }

    public static async getInstance(): Promise<ProjectService> {
        if (!ProjectService.instance) {
            ProjectService.instance = new ProjectService();
            const apiUrl = await getApiHostname();
            ProjectService.instance.apiUrl = `${apiUrl}/api/projects`;
            ProjectService.instance.initialized = true;
        }
        return ProjectService.instance;
    }

    async getProjects(page: number = 1, perPage: number = 10): Promise<IProject[]> {
        const data = await this.getAll<any[]>(page, perPage);
        return data.map(project => new Project(project));
    }

    async getProject(id: string): Promise<IProject | null> {
        const data = await this.get<any>(id);
        return data ? new Project(data) : null;
    }

    async createProject(projectData: Omit<IProject, 'id' | 'created_at' | 'updated_at' | 'getDocumentTitle'>): Promise<IProject | null> {
        const data = await this.post<any>(projectData);
        return data ? new Project(data) : null;
    }

    async updateProject(id: string, projectData: Partial<Omit<IProject, 'id' | 'created_at' | 'updated_at' | 'getDocumentTitle'>>): Promise<IProject | null> {
        const data = await this.put<any>(id, projectData);
        return data ? new Project(data) : null;
    }

    async deleteProject(id: string): Promise<boolean> {
        return await this.delete(id);
    }
}

// Update the singleton creation
let projectService: ProjectService;
export const getProjectService = async () => {
    if (!projectService) {
        projectService = await ProjectService.getInstance();
    }
    return projectService;
};

// Update the exported methods to be async
export const getProjects = async () => {
    const service = await getProjectService();
    return service.getProjects();
};

export const getProject = async (id: string) => {
    const service = await getProjectService();
    return service.getProject(id);
};

export const createProject = async (projectData: Omit<IProject, 'id' | 'created_at' | 'updated_at' | 'getDocumentTitle'>) => {
    const service = await getProjectService();
    return service.createProject(projectData);
};

export const updateProject = async (id: string, projectData: Partial<Omit<IProject, 'id' | 'created_at' | 'updated_at' | 'getDocumentTitle'>>) => {
    const service = await getProjectService();
    return service.updateProject(id, projectData);
};

export const deleteProject = async (id: string) => {
    const service = await getProjectService();
    return service.deleteProject(id);
};
