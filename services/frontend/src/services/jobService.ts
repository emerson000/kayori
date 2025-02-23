'use server'

import { Job, IJob } from '../models/job';
import { getApiHostname } from '../utils/shared';
const API_URL = await getApiHostname() + '/api/jobs';

export const getJobs = async () => {
    if (process.env.SKIP_API_CALL == 'true') {
        return [];
    }
    try {
        const response = await fetch(API_URL);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        const jobs: IJob[] = data.map((job: any) => new Job(job));
        return jobs;
    } catch (error) {
        console.error('Error fetching jobs:', error);
        throw error;
    }
};
