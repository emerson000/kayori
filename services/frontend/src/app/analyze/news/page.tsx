import Link from "next/link";
import MessageCard from "../../../components/MessageCard";
import { getNews } from "../../../services/newsService";
import SearchBar from "./searchBar";

export default async function Page({ searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }) {
    const { page = '1', search } = await searchParams;
    const pageNumber = parseInt(page as string);
    const articles = await getNews(pageNumber, search as string);
    return (
        <div>
            <h1 className="text-2xl font-bold">News</h1>
            <SearchBar search={search as string} />
            {articles && articles.map((article, index) => <MessageCard message={article} key={index} />)}
            {pageNumber > 1 && <Link className="btn float-left" href={`?page=${pageNumber - 1}`}>Previous</Link>}
            <Link className="btn float-right" href={`?page=${pageNumber + 1}`}>Next</Link>
        </div>
    );
}