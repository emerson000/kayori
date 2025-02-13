import MessageCard from "../../../components/MessageCard";
import { getNews } from "../../../services/newsService";

export async function getServiceSideProps() {
    const articles = await getNews();
    return { props: { articles }};
}

export default async function Page({ articles }) {
    return (
        <div>
            <h1 className="text-2xl font-bold">News</h1>
            {articles.map((article, index) => <MessageCard message={article} key={index} />)} 
        </div>
    );
}