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
        const seenClusters = new Set();
        const news: INewsArticle[] = await Promise.all(data.map(async (article: any) => {
            const newsArticle = new NewsArticle(article);
            if (newsArticle.cluster_id) {
                if (seenClusters.has(newsArticle.cluster_id)) {
                    return null;
                }
                seenClusters.add(newsArticle.cluster_id);
                newsArticle.cluster_articles = await getClusterArticles(article.cluster_id);
            }
            return newsArticle;
        }));
        return news.filter((article: NewsArticle) => article &&
            ((article.cluster_id
                && article.cluster_articles ? article.cluster_articles.length > 0 : true)
                || !article.cluster_id));
    } catch (error) {
        console.error('Error fetching news articles:', error);
        throw error;
    }
};

function getClusterArticles(clusterId: string) {
    return fetch(`${API_URL}?cluster=${clusterId}`)
        .then((response) => response.json())
        .then((data) => data.map((article: any) => new NewsArticle(article)));
}
