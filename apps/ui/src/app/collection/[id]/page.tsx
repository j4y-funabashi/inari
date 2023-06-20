'use client';

import { CollectionDetail, Media, collectionDetailFetcher, deleteMedia } from "@/app/apiClient";
import { MediaCard } from "@/app/components/mediaCard";
import { useState } from "react";
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

    const sortedMedia = data.media.sort(
        (a, b) => {
            if (a.date === b.date) {
                return 0
            }
            if (a.date < b.date) {
                return -1
            }
            return 1
        }
    )

    const [media, setMedia] = useState<Media[]>(sortedMedia)

    const mediaList = media.map(
        (m) => {

            const handleDelete = async function () {
                const newList = media.filter(
                    (nm) => {
                        return nm.id !== m.id
                    }
                )
                await deleteMedia(m.id)
                setMedia(newList)
            }

            return (
                <li key={m.id}>
                    <MediaCard m={m} handleDelete={handleDelete} />
                </li>
            )
        }
    )

    return (
        <ul>{mediaList}</ul>
    )
}
