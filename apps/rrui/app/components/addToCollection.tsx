'use client'

import { use } from "react";
import { Collection, Media, useCollections } from "../apiClient";

interface AddToCollectionProps {
  m: Media
  collections: Promise<Collection[]>
}

export const AddToCollection = ({ m, collections }: AddToCollectionProps) => {

  const collectionsList = use(collections).map(
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
      <button onClick={() => { }}>{c.title}</button>
    </li>
  );
}


