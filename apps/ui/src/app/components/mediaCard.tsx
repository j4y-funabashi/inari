import Link from "next/link";
import { CollectionType, Media } from "../apiClient";
import { format } from "date-fns";
import { useState } from "react";
import { AddHashtagForm } from "./AddHashtagForm";

export enum MediaCardDisplayType {
  list,
  grid,
  large,
}

interface MediaCardProps {
  m: Media;
  handleDelete: (id: string) => Promise<void>;
  saveCaption: (id: string, newCaption: string) => Promise<void>;
  saveTag: (id: string, newTag: string) => Promise<void>;
  setCurrent: (id: string) => Promise<void>;
  setNext: () => Promise<void>;
  setPrev: () => Promise<void>;
  setBack: () => Promise<void>;
  displayType: MediaCardDisplayType;
}
export const MediaCard = function ({
  m,
  displayType,
  handleDelete,
  saveCaption,
  saveTag,
  setCurrent,
  setNext,
  setPrev,
  setBack,
}: MediaCardProps) {
  const srcPrefix = process.env.NODE_ENV === "production" ? "/thumbnails/" : "";
  const srcUrl =
    displayType === MediaCardDisplayType.large
      ? `${srcPrefix}${m.thumbnails.large}`
      : `${srcPrefix}${m.thumbnails.medium}`;

  const caption = m.caption ? m.caption.trim() : "";
  const dat = m.date;
  const fdat = format(new Date(m.date), "eee, do LLL y HH:mm:ss");

  const location = formatLocation(m);

  const [newCaption, setNewCaption] = useState(caption);

  const handleCaptionSave = async function () {
    await saveCaption(m.id, newCaption);
  };

  const handleSaveHashtag = async (hashtag: string) => {
    await saveTag(m.id, hashtag);
  };

  return (
    <div>
      {displayType === MediaCardDisplayType.large && (
        <nav className="grid grid-cols-3">
          <button
            className="bg-black text-white font-bold py-1 px-2 block w-full"
            onClick={() => {
              setBack();
            }}
          >
            Back
          </button>

          <button
            className="bg-green text-white font-bold py-1 px-2 block w-full"
            onClick={() => {
              setPrev();
            }}
          >
            Prev
          </button>

          <button
            className="bg-green text-white font-bold py-1 px-2 block w-full"
            onClick={() => {
              setNext();
            }}
          >
            Next
          </button>
        </nav>
      )}
      <a
        href="#"
        onClick={() => {
          setCurrent(m.id);
        }}
      >
        <img src={srcUrl} className="" alt={caption} />
      </a>

      {displayType === MediaCardDisplayType.large && (
        <div>
          <time dateTime={dat} className="text-blue text-xs">
            {fdat}
          </time>
          {caption !== "" && (
            <p className="leading-5 text-gray-500 dark:text-gray-400">
              {caption}
            </p>
          )}
          <div className="grid grid-cols-6">
            <input
              type="text"
              value={newCaption}
              onChange={(e) => setNewCaption(e.target.value)}
              className="col-span-5 text-black"
            />
            <button
              className="bg-green text-white font-bold py-1 px-2 col-span-1"
              onClick={() => {
                handleCaptionSave();
              }}
            >
              Save
            </button>
          </div>

          {/* LocationList */}
          {location !== "" && <p>{location}</p>}

          <CollectionList m={m} />

          {/* AddATag */}
          <AddHashtagForm handleSaveHashtag={handleSaveHashtag} />

          <div className="my-4">
            <button
              className="bg-red text-white font-bold py-1 px-2 rounded"
              onClick={() => {
                handleDelete(m.id);
              }}
            >
              Delete
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

interface CollectionListProps {
  m: Media;
}
const CollectionList = ({ m }: CollectionListProps) => {
  const CollectionTypeColors = new Map<CollectionType, string>();
  CollectionTypeColors.set(CollectionType.CollectionTypeInbox, "bg-orange-700");
  CollectionTypeColors.set(CollectionType.CollectionTypeCamera, "indigo-700");
  CollectionTypeColors.set(
    CollectionType.CollectionTypeTimelineMonth,
    "indigo-700",
  );
  CollectionTypeColors.set(
    CollectionType.CollectionTypeTimelineDay,
    "indigo-700",
  );
  CollectionTypeColors.set(
    CollectionType.CollectionTypePlacesCountry,
    "indigo-700",
  );
  CollectionTypeColors.set(
    CollectionType.CollectionTypePlacesRegion,
    "indigo-700",
  );
  CollectionTypeColors.set(CollectionType.CollectionTypeHashTag, "indigo-700");

  const collections = m.collections.map((m) => {
    const bgColor = CollectionTypeColors.get(m.type);
    const className = `inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-indigo-100 ${bgColor} rounded mx-1`;
    const collectionLink = "/collection/" + m.id;
    console.log(bgColor, m.type);
    return (
      <Link key={m.id} className={className} href={collectionLink}>
        {m.title}
      </Link>
    );
  });

  return <div className="my-2">{collections}</div>;
};

const formatLocation = (m: Media): string => {
  const loc = [];
  if (m.location?.coordinates?.lat) {
    const latlng = `(${m.location.coordinates.lat}, ${m.location.coordinates.lng})`;
    loc.push(latlng);
  }
  if (m.location?.locality && m.location.locality != "") {
    loc.push(m.location.locality);
  }
  if (
    m.location?.region &&
    m.location.region != "" &&
    m.location.region !== m.location.locality
  ) {
    loc.push(m.location.region);
  }
  if (m.location?.country.long && m.location.country.long != "") {
    loc.push(m.location.country.long);
  }

  return loc.join(", ");
};
