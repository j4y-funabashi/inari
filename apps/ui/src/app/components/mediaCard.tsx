import Link from "next/link"
import { Media } from "../apiClient"
import { format } from "date-fns"
import { useState } from "react"

interface MediaCardProps {
    m: Media
    handleDelete: () => void
    saveCaption: (id: string, newCaption: string) => Promise<void>
}
export const MediaCard = function ({ m, handleDelete, saveCaption }: MediaCardProps) {
    const srcUrl = "/thumbnails/" + m.thumbnails.large

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
            <div
                className="my-8 rounded bg-gray-800">
                <figure>
                    <img src={srcUrl} className="rounded-t" alt={caption} />

                    <figcaption className="p-4">
                        <time dateTime={dat} className="text-blue text-xs">{fdat}</time>

                        {caption !== "" &&
                            <p className="leading-5 text-gray-500 dark:text-gray-400">
                                {caption}
                            </p>
                        }
                        <form onSubmit={handleCaptionSubmit}>
                            <label>
                                Caption:
                                <textarea
                                    name="newCaption"
                                    value={newCaption}
                                    onChange={e => setNewCaption(e.target.value)}
                                />
                            </label>
                            <input type="submit" value="save" />
                        </form>

                        {location !== "" &&
                            <p>{location}</p>
                        }
                    </figcaption>
                </figure>
                <ul>{collections}</ul>
                <div>
                    <button className="bg-red text-white font-bold py-1 px-2 rounded" onClick={handleDelete}>
                        Delete
                    </button>
                </div>
            </div>
        </div>
    )
}

const formatLocation = (m: Media): string => {
    const loc = []
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
