'use client';

import { CollectionDetail, collectionDetailFetcher, deleteMedia } from "@/app/apiClient";
import useSWR from "swr";

interface CollectionDetailParams {
    params: {
        id: string
    }
}

export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    const collectionID = params.id

    const { data, error, isLoading } = useSWR<CollectionDetail>(collectionID, collectionDetailFetcher)

    console.log(data, error, isLoading)

    if (!data) {
        return
    }


    return (
        <div>
            <h1>{collectionID}</h1>
            <MediaList data={data} />
        </div>
    )
}

interface MediaListProps {
    data: CollectionDetail
}

const MediaList = function ({ data }: MediaListProps) {
    const media = data.media.sort(

    )
    const mediaList = data?.media.map(
        (m) => {
            const srcUrl = "/thumbnails/" + m.thumbnails.medium
            const handleDeleteMedia = async function () {
                const newList = data.media.filter(
                    (nm) => {
                        return nm.id !== m.id
                    }
                )
                await deleteMedia(m.id)
                data.media = newList
            }

            return (
                <li key={m.id}>
                    <img src={srcUrl} />
                    <div>
                        <button className="bg-red text-white font-bold py-2 px-4 rounded" onClick={handleDeleteMedia}>
                            Delete
                        </button>
                    </div>
                </li>
            )
        }
    )

    return (
        <ul>{mediaList}</ul>
    )
}
