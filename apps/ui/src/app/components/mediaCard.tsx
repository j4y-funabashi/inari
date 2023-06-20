import Link from "next/link"
import { Media } from "../apiClient"
import { format } from "date-fns"

interface MediaCardProps {
    m: Media
    handleDelete: () => void
}
export const MediaCard = function ({ m, handleDelete }: MediaCardProps) {
    const srcUrl = "/thumbnails/" + m.thumbnails.large

    const caption = (m.caption ? m.caption.trim() : "")
    const dat = format(new Date(m.date), "eee, do LLL y HH:mm:ss")

    const collections = m.collections.map(
        (m => {
            const collectionLink = "/collection/" + m.id
            return <li key={m.id}><Link href={collectionLink}>{m.title}</Link></li>
        })
    )
    const location = formatLocation(m)

    return (
        <div>
            <img src={srcUrl} />
            <p>{dat}</p>
            {caption !== "" && <p>{caption}</p>}
            {location !== "" && <p>{location}</p>}
            <ul>{collections}</ul>
            <div>
                <button className="bg-red text-white font-bold py-2 px-4 rounded" onClick={handleDelete}>
                    Delete
                </button>
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
