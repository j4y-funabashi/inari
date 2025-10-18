'use client';
import Link from "next/link"

import { CollectionDetail, Media, NewFetchCollectionDetail, exportMedia, deleteMedia, updateMediaCaption, updateMediaHashtag } from "@/app/apiClient";
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

const updateMediaListItemExportedStatus = (m: MediaListModel, id: string): MediaListModel => {
    const ml = getMediaList(m).map((m) => {
        if (m.id === id) {
            m.is_exported = true
        }
        return m
    })

    const out = createMediaList(
        ml,
        id,
    )
    return out
}

const deleteFromMediaList = (m: MediaListModel, id: string): MediaListModel => {

    const nextMediaId = m.next.length > 0 ? m.next[0].id
        : m.prev.length > 0 ? m.prev[0].id
            : "";

    const ml = getMediaList(m).filter((m) => {
        return m.id !== id
    })

    const out = createMediaList(
        ml,
        nextMediaId
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

enum GalleryViewMode {
    grid,
    single,
    confirmDelete
}

interface MediaListProps {
    data: CollectionDetail
}

const MediaGallery = function({ data }: MediaListProps) {

    const model = createGalleryModel(data)
    const [galleryModel, setGalleryModel] = useState<MediaListModel>(model)
    const [viewMode, setViewMode] = useState<GalleryViewMode>(GalleryViewMode.grid)

    console.log(galleryModel)

    const handleConfirmDelete = async function(id: string) {
        console.log("confirm delete " + id);
        setViewMode(GalleryViewMode.confirmDelete)
    }

    const handleDelete = async function(id: string) {
        const ml = deleteFromMediaList(galleryModel, id)
        setGalleryModel(ml)
        await deleteMedia(id)
        setViewMode(GalleryViewMode.single)
    }

    const handleExport = async (id: string) => {
        const ml = updateMediaListItemExportedStatus(galleryModel, id)
        setGalleryModel(ml)
        await exportMedia(id);
        setViewMode(GalleryViewMode.grid)
    }

    const saveCaption = async (id: string, caption: string) => {
        const ml = updateMediaListItemCaption(galleryModel, id, caption)
        setGalleryModel(ml)
        await updateMediaCaption(id, caption)
    }

    const saveHashtag = async (id: string, hashtag: string) => {
        const ml = updateMediaListItemCaption(galleryModel, id, hashtag)
        setGalleryModel(ml)
        await updateMediaHashtag(id, hashtag)
    }

    const setCurrentMedia = async (id: string) => {
        const model = createMediaList(getMediaList(galleryModel), id)
        setGalleryModel(model)
        setViewMode(GalleryViewMode.single)
    }

    const handleCloseMedia = async () => {
        setViewMode(GalleryViewMode.grid)
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
                    handleDelete={handleConfirmDelete}
                    exportMedia={handleExport}
                    saveCaption={saveCaption}
                    saveHashtag={saveHashtag}
                    setCurrent={setCurrentMedia}
                    setNext={setNextMedia}
                    setPrev={setPrevMedia}
                    setBack={handleCloseMedia}
                    showNav={false}
                    showMeta={true}
                    showEditButtons={false}
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
            <nav>
                <Link href="/" className="bg-purple text-white font-bold py-1 px-2 block w-full">Collections</Link>
            </nav>
            {viewMode === GalleryViewMode.single && <main className="">
                <MediaCard
                    displayType={MediaCardDisplayType.large}
                    key={currentMedia.id}
                    m={currentMedia}
                    handleDelete={handleConfirmDelete}
                    exportMedia={handleExport}
                    saveCaption={saveCaption}
                    saveHashtag={saveHashtag}
                    setCurrent={setCurrentMedia}
                    setNext={setNextMedia}
                    setPrev={setPrevMedia}
                    setBack={handleCloseMedia}
                    showNav={true}
                    showMeta={true}
                    showEditButtons={true}
                />
            </main>
            }

            {viewMode === GalleryViewMode.grid &&
                <aside className="">
                    <h1 className="text-xl mt-4 mb-1">{data.collection_meta.title}</h1>

                    <div className="grid gap-0.5 grid-cols-3">
                        {mediaList}
                    </div>
                </aside>
            }

            {viewMode === GalleryViewMode.confirmDelete && <main className="">
                <button className="bg-red text-white font-bold py-1 px-2 rounded" onClick={() => { handleDelete(currentMedia.id) }}>
                    Confirm Delete
                </button>
            </main>
            }

        </div>

    )
}
