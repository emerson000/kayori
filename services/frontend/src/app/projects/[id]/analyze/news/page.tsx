import Link from "next/link";
import MessageCard from "@/components/MessageCard";
import SearchBar from "./searchBar";
import { NewsArticle } from "@/models/newsArticle";
import { getProject } from "@/services/projectService";
import ProjectHeader from "@/components/projects/projectHeader";
import { notFound } from "next/navigation";
import { getArtifacts } from "@/services/artifactService";
export default async function Page({ params, searchParams }: { params: Promise<{ id: string }>, searchParams: Promise<{ [key: string]: string | string[] | undefined }> }) {
    const { id } = await params;
    const project = await getProject(id);
    const { page = '1', search } = await searchParams;
    const pageNumber = parseInt(page as string);
    const articles = await getArtifacts(id, "news_article", pageNumber, 10);
    const searchEncoded = encodeURIComponent(search as string);
    if (!project) {
        notFound();
    }
    return (
        <div>
            <ProjectHeader
                project={project}
                actions={<SearchBar search={search as string} className="mt-2 w-full lg:w-1/4" id={id} />}
                currentPage="analyze"
            />
            <h1 className="text-2xl font-bold">News</h1>
            <div className="clear-both">
                {pageNumber > 1 && <Link className="btn float-left" href={`?page=${pageNumber - 1}${search ? `&search=${searchEncoded}` : ''}`}>Previous</Link>}
                <Link className="btn float-right" href={`?page=${pageNumber + 1}${search ? `&search=${searchEncoded}` : ''}`}>Next</Link>
            </div>
            {articles && articles.length > 0 && <div className="clear-both mt-15">
                {articles && articles.map((article: NewsArticle, index) => <MessageCard message={article} key={index} />)}
            </div>}
            {articles && articles.length == 0 && <div className="text-center clear-both py-10 bg-base-200 mt-15 mb-5">No results found</div>}
            <div className="clear-both mb-20">
                {pageNumber > 1 && <Link className="btn float-left" href={`?page=${pageNumber - 1}${search ? `&search=${searchEncoded}` : ''}`}>Previous</Link>}
                <Link className="btn float-right" href={`?page=${pageNumber + 1}${search ? `&search=${searchEncoded}` : ''}`}>Next</Link>
            </div>
        </div>
    );
}