
export default function Page() {
    return (
        <div>
            <div className="card bg-base-200 w-96 shadow-sm">
                <div className="card-body">
                    <h2 className="card-title">News Articles</h2>
                    <p>Analyze news articles collected from various sources.</p>
                    <div className="card-actions justify-end">
                        <a href="/analyze/news" className="btn btn-primary">Analyze</a>
                    </div>
                </div>
            </div>
        </div>
    );
}