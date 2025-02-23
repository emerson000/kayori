'use server'

import { NewsArticle, INewsArticle } from '../models/newsArticle';
import { getApiHostname } from '../utils/shared';
const API_URL = await getApiHostname() + '/api/entities/news_articles';

export const getNews = async () => {
    if (process.env.SKIP_API_CALL == 'true') {
        return [];
    }
    try {
        const response = await fetch(API_URL);
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
