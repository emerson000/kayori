'use server'

import { NewsArticle, INewsArticle } from '../models/newsArticle';
import { getApiHostname } from '../utils/shared';
const API_URL = await getApiHostname() + '/api/entities/news_articles';

export const getNews = async (page: number, search: string) => {
    if (process.env.SKIP_API_CALL == 'true') {
        return [];
    }
    try {
        let path = `${API_URL}?page=${page}`;
        if (search) {
            path += `&search=${search}`;
        }
        const response = await fetch(path);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        const news: INewsArticle[] = data.map((article: any) => new NewsArticle(article));
        return news;
    } catch (error) {
        console.error('Error fetching news articles:', error);
        throw error;
    }
};
