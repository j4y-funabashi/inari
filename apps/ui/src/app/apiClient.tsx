import { Fetcher } from "swr";

interface Media {
    id: string
    thumbnails: Thumbnails
}
interface Thumbnails {
    medium: string
    large: string
}

interface Collection {
    id: string;
    title: string;
    media_count: number;
    type: string;
}

interface CollectionDetail {
    collection_meta: Collection;
    media: Media[]
}

const getCollectionsByType = async function (type: string): Promise<Collection[]> {
    const res = await fetch("/api/timeline/months")
    console.log(res)

    return res.json()
}


const getCollectionDetail = async function (id: string): Promise<CollectionDetail> {
    const res = await fetch("/api/timeline/month/" + id)
    console.log(res)

    return res.json()
}

export const collectionListFetcher: Fetcher<Collection[], string> = (type) => getCollectionsByType(type)
export const collectionDetailFetcher: Fetcher<CollectionDetail, string> = (id) => getCollectionDetail(id)
