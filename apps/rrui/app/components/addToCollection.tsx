'use client'

import { type Collection, type Media } from "~/apiClient";

interface AddToCollectionProps {
  m: Media
  collections: Collection[]
  saveHashtag: (id: string, newHashtag: string) => Promise<void>
}

export const AddToCollection = ({ m, collections, saveHashtag }: AddToCollectionProps) => {

  const collectionsList = collections.map(
    (c) => {
      return (
        <MediaCollection c={c} m={m} key={c.id} saveHashtag={saveHashtag} />
      );
    }
  )

  return (
    <div>
      <ul>{collectionsList}</ul>
    </div>
  )
}

interface MediaCollectionProps {
  c: Collection
  m: Media
  saveHashtag: (id: string, newHashtag: string) => Promise<void>
}

const MediaCollection = ({ c, m, saveHashtag }: MediaCollectionProps) => {
  return (
    <li key={c.id}>
      <button className='p-2 block hover:bg-slate-900' onClick={() => { saveHashtag(m.id, c.title) }}>{c.title}</button>
    </li>
  );
}


