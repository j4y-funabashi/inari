'use client'

import { type Collection, type Media } from "~/apiClient";

interface AddToCollectionProps {
  m: Media
  collections: Collection[]
}

export const AddToCollection = ({ m, collections }: AddToCollectionProps) => {

  const collectionsList = collections.map(
    (c) => {
      return (
        <MediaCollection c={c} m={m} key={c.id} />
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
}

const MediaCollection = ({ c, m }: MediaCollectionProps) => {
  return (
    <li key={c.id}>
      <button className='p-2 block hover:bg-slate-900' onClick={() => { }}>{c.title}</button>
    </li>
  );
}


