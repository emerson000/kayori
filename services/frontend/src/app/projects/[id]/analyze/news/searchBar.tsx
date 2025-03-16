'use client'

import Form from "next/form";
import { useEffect, useRef } from "react";
import { useRouter } from "next/navigation";

interface SearchBarProps {
    search: string;
    className?: string;
    id: string;
}

export default function SearchBar({ search, className, id }: SearchBarProps) {
    const router = useRouter();
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === '/' && document.activeElement?.tagName !== 'INPUT') {
                event.preventDefault();
                inputRef.current?.focus();
                const length = inputRef.current?.value.length || 0;
                inputRef.current?.setSelectionRange(0, length);
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, []);

    const handleSubmit = async (formData: FormData, id: string) => {
        const search = formData.get('search') as string;
        router.push(`/projects/${id}/analyze/news?search=${search}`);
        inputRef.current?.blur();
    }

    return <Form action={(formData) => handleSubmit(formData, id)} className={className}>
        <input ref={inputRef} type="text" placeholder="Search" name="search" className="input w-full" defaultValue={search} />
    </Form>
}