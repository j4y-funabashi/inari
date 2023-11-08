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
    setCurrent: (id: string) => Promise<void>
    displayType: MediaCardDisplayType
}
export const MediaCard = function ({ m, displayType, handleDelete, saveCaption, setCurrent }: MediaCardProps) {

    const srcPrefix = process.env.NODE_ENV === "production" ? "/thumbnails/" : ""
    const srcUrl = displayType === MediaCardDisplayType.large ? `${srcPrefix}${m.thumbnails.large}`
        : `${srcPrefix}${m.thumbnails.medium}`

    const caption = (m.caption ? m.caption.trim() : "")
    const dat = m.date
    const fdat = format(new Date(m.date), "eee, do LLL y HH:mm:ss")

    const collections = m.collections.map(
        (m => {
            const collectionLink = "/collection/" + m.id
            return <li key={m.id}><Link className="text-white bg-gray text-xs" href={collectionLink}>{m.title}</Link></li>
        })
    )
    const location = formatLocation(m)


    const [newCaption, setNewCaption] = useState(caption);

    const handleCaptionSubmit = async function (event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault()
        await saveCaption(m.id, newCaption)
    }

    return (
        <div>
            <a href="#" onClick={() => { setCurrent(m.id) }}>
                <img src={srcUrl} className="rounded-t" alt={caption} />
            </a>

            {displayType === MediaCardDisplayType.large &&
                <div>
                    <time dateTime={dat} className="text-blue text-xs">{fdat}</time>
                    {caption !== "" &&
                        <p className="leading-5 text-gray-500 dark:text-gray-400">
                            {caption}
                        </p>
                    }
                    <form onSubmit={handleCaptionSubmit}>
                        <input
                            type="text"
                            value={newCaption}
                            onChange={e => setNewCaption(e.target.value)}
                        />
                        <input type="submit" value="save" />
                    </form>
                    {location !== "" &&
                        <p>{location}</p>
                    }
                    <ul>{collections}</ul>
                    <div>
                        <button className="bg-red text-white font-bold py-1 px-2 rounded" onClick={() => { handleDelete(m.id) }}>
                            Delete
                        </button>
                    </div>
                </div>
            }
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
