'use client'

import Link from "next/link";
import { Collection } from "../apiClient";
import { use } from "react";

interface CollectionListProps {
    data: Promise<Collection[]>
}

export const CollectionList = ({ data }: CollectionListProps) => {

    const collections = use(data).map(
        (c) => {
            return (
                <MediaCollection c={c} key={c.id} />
            );
        }
    )

    return (
        <div>
            <ul>{collections}</ul>
        </div>
    )
}

interface MediaCollectionProps {
    c: Collection
}

const MediaCollection = ({ c }: MediaCollectionProps) => {
    const collectionLink = "/collection/" + c.id
    const mediaCount = c.media_count ?? 0;
    const exportedCount = c.exported_count ?? 0;

    return (
        <li key={c.id}>
            <Link className='p-2 block hover:bg-slate-900' href={collectionLink}>{c.title} ({exportedCount} / {mediaCount})</Link>
            <progress value={exportedCount} max={mediaCount} />
        </li>
    );
}

