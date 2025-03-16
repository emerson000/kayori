export abstract class BaseService {
    protected apiUrl: string;

    constructor(apiUrl: string) {
        this.apiUrl = apiUrl;
    }

    protected async get<T>(id?: string): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const url = id ? `${this.apiUrl}/${id}` : this.apiUrl;
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

    protected async getAll<T>(page?: number, perPage?: number): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        const url = `${this.apiUrl}?page=${page}&per_page=${perPage}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return await response.json();
    }

    protected async post<T>(data: any): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const response = await fetch(this.apiUrl, {
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

    protected async put<T>(id: string, data: any): Promise<T> {
        if (process.env.SKIP_API_CALL == 'true') {
            return null as T;
        }
        try {
            const response = await fetch(`${this.apiUrl}/${id}`, {
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

    protected async delete(id: string): Promise<boolean> {
        if (process.env.SKIP_API_CALL == 'true') {
            return true;
        }
        try {
            const response = await fetch(`${this.apiUrl}/${id}`, {
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