export default function Page() {
    return (
        <div>
            <div className="card bg-base-200 w-96 shadow-sm">
                <div className="card-body">
                    <h2 className="card-title">Deduplication</h2>
                    <p>Deduplicate collected artifacts.</p>
                    <div className="card-actions justify-end">
                        <a href="/process/deduplicate" className="btn btn-primary">Start</a>
                    </div>
                </div>
            </div>
        </div>
    );
}