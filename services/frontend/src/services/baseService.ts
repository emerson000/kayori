export abstract class BaseService {
    protected apiUrl: string;
    protected initialized: boolean = false;

    constructor(apiUrl: string) {
        this.apiUrl = apiUrl;
    }

    protected async get<T>(id?: string, urlOverride?: string): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const baseUrl = urlOverride || this.apiUrl;
            const url = id ? `${baseUrl}/${id}` : baseUrl;
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching data:', error);
            throw error;
        }
    }

    protected async getAll<T>(page?: number, perPage?: number, urlOverride?: string): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        const baseUrl = urlOverride || this.apiUrl;
        const url = `${baseUrl}?page=${page}&per_page=${perPage}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return await response.json();
    }

    protected async post<T>(data: any, urlOverride?: string): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const url = urlOverride || this.apiUrl;
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating resource:', error);
            throw error;
        }
    }

    protected async put<T>(id: string, data: any, urlOverride?: string): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const baseUrl = urlOverride || this.apiUrl;
            console.log({baseUrl, id, data})
            const response = await fetch(`${baseUrl}/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return await response.json();
        } catch (error) {
            console.error('Error updating resource:', error);
            throw error;
        }
    }

    protected async delete(id: string, urlOverride?: string): Promise<boolean> {
        if (process.env.SKIP_API_CALL == 'true') {
            return true;
        }
        try {
            const baseUrl = urlOverride || this.apiUrl;
            const response = await fetch(`${baseUrl}/${id}`, {
                method: 'DELETE',
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return true;
        } catch (error) {
            console.error('Error deleting resource:', error);
            throw error;
        }
    }
}