'use client';

import { CollectionDetail, Media, NewFetchCollectionDetail, deleteMedia, updateMediaCaption } from "@/app/apiClient";
import { MediaCard, MediaCardDisplayType } from "@/app/components/mediaCard";
import { useState } from "react";
import useSWR from "swr";

interface CollectionDetailParams {
    params: {
        id: string
    }
}

interface MediaListModel {
    prev: Media[]
    current: Media
    next: Media[]
}

const getMediaList = (m: MediaListModel): Media[] => {
    if (!m.current) return [...m.prev, ...m.next]
    return [...m.prev, m.current, ...m.next]
}

const deleteFromMediaList = (m: MediaListModel, id: string): MediaListModel => {
    const ml = getMediaList(m).filter((m) => {
        return m.id !== id
    })

    const out = createMediaList(
        ml,
        ml[0] ? ml[0].id : "" // TODO this should be the next media?
    )
    return out
}

const updateMediaListItemCaption = (m: MediaListModel, id: string, caption: string): MediaListModel => {
    const ml = getMediaList(m).map((m) => {
        if (m.id === id) {
            m.caption = caption
        }
        return m
    })

    const out = createMediaList(
        ml,
        id,
    )
    return out
}

export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    const collectionID = params.id
    const collectionDetailFetcher = NewFetchCollectionDetail(process.env.NODE_ENV)

    const { data: collectionDetailData, error: collectionDetailError, isLoading: collectionDetailLoading } = useSWR<CollectionDetail>(collectionID, collectionDetailFetcher)
    if (!collectionDetailData) {
        return
    }


    return (
        <div>
            <MediaGallery data={collectionDetailData} />
        </div>
    )
}

interface MediaListProps {
    data: CollectionDetail
}

const createMediaList = (media: Media[], currentID: string): MediaListModel => {

    const current = media.filter(
        (m) => {
            return m.id === currentID
        }
    )
    const prev: Media[] = []
    const next: Media[] = []

    var pr = true
    media.forEach((m) => {
        if (m.id === currentID) {
            pr = false
            return
        }
        if (pr) {
            prev.push(m)
            return
        }
        next.push(m)
    })

    return {
        prev: prev,
        current: current[0],
        next: next
    }
}

const createGalleryModel = (data: CollectionDetail): MediaListModel => {
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

    const model = createMediaList(sortedMedia, sortedMedia[0].id)

    return model
}

const getCurrentMedia = (model: MediaListModel): Media => {
    return model.current
}


const MediaGallery = function ({ data }: MediaListProps) {


    const model = createGalleryModel(data)
    const [galleryModel, setGalleryModel] = useState<MediaListModel>(model)

    console.log(galleryModel)

    const handleDelete = async function (id: string) {
        const ml = deleteFromMediaList(galleryModel, id)
        setGalleryModel(ml)
        await deleteMedia(id)
    }

    const saveCaption = async (id: string, caption: string) => {
        const ml = updateMediaListItemCaption(galleryModel, id, caption)
        setGalleryModel(ml)
        await updateMediaCaption(id, caption)
    }

    const setCurrentMedia = async (id: string) => {
        const model = createMediaList(getMediaList(galleryModel), id)
        setGalleryModel(model)
    }

    const setNextMedia = async () => {
        if (!galleryModel.next.length) {
            return
        }

        console.log("setting next media")
        const model = createMediaList(getMediaList(galleryModel), galleryModel.next[0].id)
        setGalleryModel(model)
    }
    const setPrevMedia = async () => {
        if (!galleryModel.prev.length) {
            return
        }

        console.log("setting prev media")
        const model = createMediaList(getMediaList(galleryModel), galleryModel.prev[galleryModel.prev.length - 1].id)
        setGalleryModel(model)
    }

    const media = getMediaList(galleryModel)

    const mediaList = media.map(
        (m) => {
            return (
                <MediaCard
                    displayType={MediaCardDisplayType.list}
                    key={m.id}
                    m={m}
                    handleDelete={handleDelete}
                    saveCaption={saveCaption}
                    setCurrent={setCurrentMedia}
                    setNext={setNextMedia}
                    setPrev={setPrevMedia}
                />
            )
        }
    )
    const currentMedia = getCurrentMedia(galleryModel)

    if (!media.length) {
        return (
            <div>empty gallery</div>
        )
    }
    return (

        <div className="">
            <main className="">
                <MediaCard
                    displayType={MediaCardDisplayType.large}
                    key={currentMedia.id}
                    m={currentMedia}
                    handleDelete={handleDelete}
                    saveCaption={saveCaption}
                    setCurrent={setCurrentMedia}
                    setNext={setNextMedia}
                    setPrev={setPrevMedia}
                />
            </main>

            <aside className="">
                <h1 className="text-xl mt-4 mb-1">{data.collection_meta.title}</h1>

                <div className="grid gap-0.5 grid-cols-4">
                    {mediaList}
                </div>
            </aside>

        </div>

    )
}
