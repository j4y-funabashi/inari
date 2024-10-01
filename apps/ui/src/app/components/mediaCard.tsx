import Link from "next/link"
import { Media } from "../apiClient"
import { format } from "date-fns"
import { useState } from "react"

export enum MediaCardDisplayType {
    list,
    grid,
    large,
}

interface MediaCardProps {
    m: Media
    handleDelete: (id: string) => Promise<void>
    saveCaption: (id: string, newCaption: string) => Promise<void>
    saveHashtag: (id: string, newHashtag: string) => Promise<void>
    setCurrent: (id: string) => Promise<void>
    setNext: () => Promise<void>
    setPrev: () => Promise<void>
    setBack: () => Promise<void>
    displayType: MediaCardDisplayType
}
export const MediaCard = function ({ m, displayType, handleDelete, saveCaption, saveHashtag, setCurrent, setNext, setPrev, setBack }: MediaCardProps) {

    const srcPrefix = process.env.NODE_ENV === "production" ? "/thumbnails/" : ""
    const srcUrl = displayType === MediaCardDisplayType.large ? `${srcPrefix}${m.thumbnails.large}`
        : `${srcPrefix}${m.thumbnails.medium}`

    const caption = (m.caption ? m.caption.trim() : "")
    const dat = m.date
    const fdat = format(new Date(m.date), "eee, do LLL y HH:mm:ss")

    const location = formatLocation(m)


    const [newCaption, setNewCaption] = useState(caption);
    const [newHashtag, setNewHashtag] = useState("");

    const handleCaptionSave = async function () {
        await saveCaption(m.id, newCaption)
    }

    const handleHashtagSave = async function () {
        await saveHashtag(m.id, newHashtag)
    }

    return (
        <div>
            {displayType === MediaCardDisplayType.large &&
                <nav className="grid grid-cols-3">
                    <button className="bg-black text-white font-bold py-1 px-2 block w-full" onClick={() => { setBack() }}>
                        Back
                    </button>

                    <button className="bg-green text-white font-bold py-1 px-2 block w-full" onClick={() => { setPrev() }}>
                        Prev
                    </button>

                    <button className="bg-green text-white font-bold py-1 px-2 block w-full" onClick={() => { setNext() }}>
                        Next
                    </button>
                </nav>
            }
            <a href="#" onClick={() => { setCurrent(m.id) }}>
                <img src={srcUrl} className="" alt={caption} />
            </a>

            {displayType === MediaCardDisplayType.large &&
                <div>
                    <time dateTime={dat} className="text-blue text-xs">{fdat}</time>
                    {caption !== "" &&
                        <p className="leading-5 text-gray-500 dark:text-gray-400">
                            {caption}
                        </p>
                    }
                    <div className="grid grid-cols-6">
                        <input
                            type="text"
                            value={newCaption}
                            onChange={e => setNewCaption(e.target.value)}
                            className="col-span-5 text-black"
                        />
                        <button
                            className="bg-green text-white font-bold py-1 px-2 col-span-1"
                            onClick={() => { handleCaptionSave() }}>Save</button>
                    </div>

                    <div className="grid grid-cols-6">
                        <input
                            type="text"
                            value={newHashtag}
                            onChange={e => setNewHashtag(e.target.value)}
                            placeholder="add a hashtag"
                            className="col-span-5 text-black"
                        />
                        <button
                            className="bg-green text-white font-bold py-1 px-2 col-span-1"
                            onClick={() => { handleHashtagSave() }}>Save</button>
                    </div>

                    {/* LocationList */}
                    {location !== "" &&
                        <p>{location}</p>
                    }

                    <CollectionList m={m} />

                    <div className="my-4">
                        <button className="bg-red text-white font-bold py-1 px-2 rounded" onClick={() => { handleDelete(m.id) }}>
                            Delete
                        </button>
                    </div>

                </div>
            }
        </div>
    )
}

interface CollectionListProps {
    m: Media
}
const CollectionList = ({ m }: CollectionListProps) => {
    const collections = m.collections.map(
        (m => {
            const collectionLink = "/collection/" + m.id
            return <Link key={m.id} className="inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-indigo-100 bg-indigo-700 rounded mx-1" href={collectionLink}>{m.title}</Link>
        })
    )

    return (
        <div className="my-2">
            {collections}
        </div>
    )
}

const formatLocation = (m: Media): string => {
    const loc = []
    if (m.location?.coordinates?.lat) {
        const latlng = `(${m.location.coordinates.lat}, ${m.location.coordinates.lng})`
        loc.push(latlng)

    }
    if (m.location?.locality && m.location.locality != "") {
        loc.push(m.location.locality)
    }
    if (m.location?.region && m.location.region != "" && (m.location.region !== m.location.locality)) {
        loc.push(m.location.region)
    }
    if (m.location?.country.long && m.location.country.long != "") {
        loc.push(m.location.country.long)
    }

    return loc.join(", ")
}
