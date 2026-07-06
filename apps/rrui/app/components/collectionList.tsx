import { type Collection } from "~/apiClient";

interface CollectionListProps {
    data: Collection[]
}

export const CollectionList = ({ data }: CollectionListProps) => {

    const collections = data.map(
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
            <a className='p-2 block hover:bg-slate-900' href={collectionLink}>{c.title} ({exportedCount} / {mediaCount})</a>
            <progress value={exportedCount} max={mediaCount} />
        </li>
    );
}

