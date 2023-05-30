'use client';

import { collectionDetailFetcher } from "@/app/apiClient";
import useSWR from "swr";

interface CollectionDetailParams {
    params: {
        id: string
    }
}

export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    const collectionID = params.id

    const { data, error, isLoading } = useSWR(collectionID, collectionDetailFetcher)

    console.log(data, error, isLoading)

    const mediaList = data?.media.map(
        (m) => {
            const srcUrl = "/thumbnails/" + m.thumbnails.medium
            return (
                <li key={m.id}>
                    <img src={srcUrl} />
                    <div>
                        <button className="bg-red text-white font-bold py-2 px-4 rounded">
                            Delete
                        </button>
                    </div>
                </li>
            )
        }
    )

    return (
        <div>
            <h1>{collectionID}</h1>
            <ul>{mediaList}</ul>
        </div>
    )
}
