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
const mockGetCollectionsByType = function (type: string): Collection[] {
    return [
        getMockCollection()
    ]
}


const getCollectionDetail = async function (id: string): Promise<CollectionDetail> {
    const res = await fetch("/api/timeline/month/" + id)
    console.log(res)

    return res.json()
}
const mockGetCollectionDetail = function (id: string): CollectionDetail {
    return {
        collection_meta: getMockCollection(),
        media: [
            getMockMedia(),
            getMockMedia(),
            getMockMedia(),
            getMockMedia(),
            getMockMedia(),
        ]
    }
}

const getMockMedia = (): Media => {
    const uuid = crypto.randomUUID()
    return {
        id: `testid-${uuid}`,
        thumbnails: {
            medium: "https://placekitten.com/420/420",
            large: "https://placekitten.com/1080/600",
        },
        date: "2022-01-28T10:01:02Z",
        collections: [
            getMockCollection(),
            getMockCollection(),
            getMockCollection(),
        ],
        location: {
            country: { long: "Country", short: "c" },
            region: "Region",
            locality: "Locality",
        },
        caption: "This is the caption",
    }
}

const getMockCollection = (): Collection => {
    const uuid = crypto.randomUUID()
    return {
        id: `test-1-${uuid}`,
        title: `c ${uuid}`,
        media_count: 5,
        type: "inbox"
    }
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

export const NewCollectionLister = (env: string): Fetcher<Collection[], string> => {
    switch (env) {
        case "production":
            const collectionListFetcher: Fetcher<Collection[], string> = (type) => getCollectionsByType(type)
            return collectionListFetcher

        default:
            return mockGetCollectionsByType
    }
}

export const NewFetchCollectionDetail = (env: string): Fetcher<CollectionDetail, string> => {
    switch (env) {
        case "production":
            const fetchCollectionDetail: Fetcher<CollectionDetail, string> = (id) => getCollectionDetail(id)
            return fetchCollectionDetail

        default:
            return mockGetCollectionDetail
    }

}

