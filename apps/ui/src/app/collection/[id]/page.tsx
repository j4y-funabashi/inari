'use client'

import { useCollectionDetail, useCollections } from "@/app/apiClient";
import { MediaGallery } from "@/app/components/mediaGallery";

interface CollectionDetailParams {
    params: {
        id: string
    }
}


export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    const collectionsData = useCollections("hashtag")

    const collectionID = params.id
    const collectionDetailData = useCollectionDetail(collectionID)

    return (
        <div>
            <MediaGallery collectionDetail={collectionDetailData} collections={collectionsData} />
        </div>
    )
}
