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
    return [...m.prev, m.current, ...m.next]
}

export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    const collectionID = params.id
    const collectionDetailFetcher = NewFetchCollectionDetail(process.env.NODE_ENV)

    const { data, error, isLoading } = useSWR<CollectionDetail>(collectionID, collectionDetailFetcher)

    console.log(data, error, isLoading)

    if (!data) {
        return
    }


    return (
        <div>
            <MediaGallery data={data} />
        </div>
    )
}

interface MediaListProps {
    data: CollectionDetail
}

const createMediaList = (media: Media[], currentID: string): MediaListModel => {
    const model: MediaListModel = {
        prev: [],
        current: media[0],
        next: [],
    }
    var prev = true
    media.forEach((m) => {
        if (m.id === currentID) {
            model.current = m
            prev = false
            return
        }
        if (prev) {
            model.prev.push(m)
            return
        }
        model.next.push(m)
    })

    return model
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

    const handleDelete = async function () {
    }

    const saveCaption = async (id: string, caption: string) => {
    }

    const setCurrentMedia = async (id: string) => {
        const model = createMediaList(getMediaList(galleryModel), id)
        setGalleryModel(model)
    }

    const media = getMediaList(galleryModel)
    const mediaList = media.map(
        (m) => {

            return (
                <MediaCard displayType={MediaCardDisplayType.list} key={m.id} m={m} handleDelete={handleDelete} saveCaption={saveCaption} setCurrent={setCurrentMedia} />
            )
        }
    )
    const currentMedia = getCurrentMedia(galleryModel)

    return (

        <div className="grid grid-cols-7">
            <aside className="col-span-2 overflow-scroll h-screen">
                <h1 className="">{data.collection_meta.title}</h1>

                <div className="grid grid-cols-3">
                    {mediaList}
                </div>
            </aside>

            <main className="col-span-5">
                <MediaCard displayType={MediaCardDisplayType.large} key={currentMedia.id} m={currentMedia} handleDelete={handleDelete} saveCaption={saveCaption} setCurrent={setCurrentMedia} />
            </main>
        </div>

    )
}
