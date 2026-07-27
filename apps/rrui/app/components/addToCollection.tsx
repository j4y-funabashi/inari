'use client'

import { useState } from "react";
import { type Collection, type Media } from "~/apiClient";

interface AddToCollectionProps {
  m: Media
  collections: Collection[]
  saveHashtag: (id: string, newHashtag: string) => Promise<void>
}

export const AddToCollection = ({ m, collections, saveHashtag }: AddToCollectionProps) => {

  const [newHashtag, setNewHashtag] = useState("");
  const handleHashtagSave = async function() {
    await saveHashtag(m.id, newHashtag)
  }


  const collectionsList = collections.map(
    (c) => {
      return (
        <MediaCollection c={c} m={m} key={c.id} saveHashtag={saveHashtag} />
      );
    }
  )

  return (
    <div>
      <div className="grid grid-cols-6">
        <input
          type="text"
          value={newHashtag}
          onChange={e => setNewHashtag(e.target.value)}
          placeholder="add a collection"
          className="col-span-5 text-black bg-white"
        />
        <button
          className="bg-lime-600 text-white font-bold py-1 px-2 col-span-1"
          onClick={() => { handleHashtagSave() }}>Save</button>
      </div>


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


