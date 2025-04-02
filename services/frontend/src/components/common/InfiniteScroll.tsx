import { ReactNode, useEffect, useRef, useState } from 'react';

interface InfiniteScrollProps<T> {
    initialData: T[];
    loadMore: (page: number) => Promise<T[]>;
    children: (items: T[], loading: boolean) => ReactNode;
    threshold?: number;
}

export default function InfiniteScroll<T>({ 
    initialData, 
    loadMore, 
    children, 
    threshold = 1.2 
}: InfiniteScrollProps<T>) {
    const [items, setItems] = useState<T[]>(initialData);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setItems(initialData);
        setPage(1);
        setHasMore(true);
    }, [initialData]);

    const loadMoreItems = async () => {
        if (loading || !hasMore) return;
        
        setLoading(true);
        const nextPage = page + 1;
        const newItems = await loadMore(nextPage);
        
        setItems(prev => [...prev, ...newItems]);
        setPage(nextPage);
        setHasMore(newItems.length > 0);
        setLoading(false);
    };

    const handleScroll = () => {
        if (!containerRef.current) return;
        
        const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
        if (scrollHeight - scrollTop <= clientHeight * threshold) {
            loadMoreItems();
        }
    };

    return (
        <div 
            ref={containerRef}
            onScroll={handleScroll}
            className="w-full"
            style={{ maxHeight: 'calc(100vh - 200px)', overflowY: 'auto' }}
        >
            {children(items, loading)}
        </div>
    );
} 