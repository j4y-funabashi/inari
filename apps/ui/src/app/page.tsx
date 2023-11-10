'use client';
import useSWR from 'swr';
import { NewCollectionLister } from './apiClient';
import Link from 'next/link';

export default function Home() {

  return (
    <CollectionList />
  )
}

const CollectionList = function () {

  const collectionLister = NewCollectionLister(process.env.NODE_ENV)

  const { data, error, isLoading } = useSWR('/api/timeline/months', collectionLister)

  if (error) return <div>failed to load</div>
  if (isLoading) return <div>loading...</div>

  console.log(data, error, isLoading)

  const collections = data?.map(
    (c) => {
      const collectionLink = "/collection/" + c.id
      return <li key={c.id}><Link className='p-2 block hover:bg-slate-900' href={collectionLink}>{c.title} ({c.media_count})</Link></li>
    }
  )

  return (
    <div>
      <h1>Hello!</h1>
      <ul>{collections}</ul>
    </div>
  )
}
