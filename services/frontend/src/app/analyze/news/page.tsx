import Link from "next/link";
import MessageCard from "../../../components/MessageCard";
import { getNews } from "../../../services/newsService";
import SearchBar from "./searchBar";
import { NewsArticle } from "../../../models/newsArticle";

export default async function Page({ searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }) {
    const { page = '1', search } = await searchParams;
    const pageNumber = parseInt(page as string);
    const articles = await getNews(pageNumber, search as string);
    const searchEncoded = encodeURIComponent(search as string);
    return (
        <div>
            <h1 className="text-2xl font-bold">News</h1>
            <SearchBar search={search as string} />
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