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
        
        const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
        const containerBottom = containerRef.current.getBoundingClientRect().bottom;
        const windowHeight = window.innerHeight;
        
        if (containerBottom <= windowHeight * threshold) {
            loadMoreItems();
        }
    };

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, [page, loading, hasMore]);

    return (
        <div 
            ref={containerRef}
            className="w-full h-full"
        >
            {children(items, loading)}
        </div>
    );
} 