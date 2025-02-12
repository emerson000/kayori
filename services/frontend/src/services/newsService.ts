import { NewsArticle, INewsArticle } from '../models/newsArticle';
import { getApiHostname } from '../utils/shared';
const API_URL = getApiHostname() + '/api/entities/news_articles';

export const getNews = async () => {
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
