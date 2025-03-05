'use client'

import Form from "next/form";
import { searchNews } from "./actions";
import { useEffect, useRef } from "react";

interface SearchBarProps {
    search: string;
}

export default function SearchBar({ search }: SearchBarProps) {
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === '/' && document.activeElement?.tagName !== 'INPUT') {
                event.preventDefault();
                inputRef.current?.focus();
                const length = inputRef.current?.value.length || 0;
                inputRef.current?.setSelectionRange(length, length);
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, []);

    return <Form action={searchNews}>
        <input ref={inputRef} type="text" placeholder="Search" name="search" className="input float-right mb-4" defaultValue={search} />
    </Form>
}