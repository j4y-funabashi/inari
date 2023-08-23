import { Fetcher } from "swr";

export interface Media {
    id: string
    thumbnails: Thumbnails
    collections: Collection[]
    date: string
    location?: Location
    caption?: string
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

interface Location {
    country: Country
    region: string
    locality: string
    coordinates?: Coordinates
}
interface Country {
    short: string
    long: string
}
interface Coordinates {
    lat: number
    lng: number
}

export interface CollectionDetail {
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

export const deleteMedia = async function (id: string) {
    const requestOptions: RequestInit = {
        method: "DELETE",
        headers: { 'Content-Type': 'application/json' }
    }
    const res = await fetch("/api/media/" + id, requestOptions)
    console.log(res)
}

export const updateMediaCaption = async function (id: string, caption: string) {
    const requestOptions: RequestInit = {
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        body: caption
    }
    const res = await fetch("/api/media/" + id + "/caption", requestOptions)
    console.log(res)
}

export const collectionListFetcher: Fetcher<Collection[], string> = (type) => getCollectionsByType(type)
export const collectionDetailFetcher: Fetcher<CollectionDetail, string> = (id) => getCollectionDetail(id)
