'use client';

import { useEffect, useState } from "react";
import { use } from "react";
import MessageCard from "@/components/MessageCard";
import SearchBar from "./searchBar";
import { NewsArticle } from "@/models/newsArticle";
import { getProject } from "@/services/projectService";
import ProjectHeader from "@/components/projects/projectHeader";
import { notFound } from "next/navigation";
import { getArtifacts } from "@/services/artifactService";
import { Project } from "@/models/project";
import InfiniteScroll from "@/components/common/InfiniteScroll";

export default function Page({ params, searchParams }: { params: Promise<{ id: string }>, searchParams: Promise<{ [key: string]: string | string[] | undefined }> }) {
    const { id } = use(params);
    const { search } = use(searchParams);
    const [project, setProject] = useState<Project | null>(null);
    const [initialArticles, setInitialArticles] = useState<NewsArticle[]>([]);

    useEffect(() => {
        const loadInitialData = async () => {
            const projectData = await getProject(id, true);
            if (!projectData) {
                notFound();
            }
            setProject(new Project(projectData));
            const articles = await getArtifacts(id, "news_article", 1, 10, true, search as string);
            setInitialArticles(articles as unknown as NewsArticle[]);
        };
        loadInitialData();
    }, [id, search]);

    const loadMoreArticles = async (page: number) => {
        const articles = await getArtifacts(id, "news_article", page, 10, true, search as string);
        return articles as unknown as NewsArticle[];
    };

    if (!project) {
        return <div>Loading...</div>;
    }

    return (
        <div>
            <ProjectHeader
                project={project}
                actions={<SearchBar search={search as string} className="mt-2 w-full lg:w-1/4" id={id} />}
                currentPage="analyze"
            />
            <h1 className="text-2xl font-bold sticky top-0 bg-base-100 z-10">News</h1>
            <InfiniteScroll
                initialData={initialArticles}
                loadMore={loadMoreArticles}
                threshold={1.2}
            >
                {(articles, loading) => (
                    <div>
                        {articles.length > 0 ? (
                            <div className="clear-both mt-15">
                                {articles.map((article: NewsArticle, index) => (
                                    <MessageCard message={article} key={index} />
                                ))}
                            </div>
                        ) : (
                            <div className="text-center clear-both py-10 bg-base-200 mt-15 mb-5">
                                No results found
                            </div>
                        )}
                        {loading && (
                            <div className="text-center py-4">
                                Loading more articles...
                            </div>
                        )}
                    </div>
                )}
            </InfiniteScroll>
        </div>
    );
}