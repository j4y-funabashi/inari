'use client';
import useSWR from 'swr';
import { collectionListFetcher } from './apiClient';
import Link from 'next/link';

export default function Home() {

  return (
    <CollectionList />
  )
}

const CollectionList = function () {

  const { data, error, isLoading } = useSWR('/api/timeline/months', collectionListFetcher)

  if (error) return <div>failed to load</div>
  if (isLoading) return <div>loading...</div>

  console.log(data, error, isLoading)

  const collections = data?.map(
    (c) => {
      const collectionLink = "/collection/" + c.id
      return <li key={c.id}><Link href={collectionLink}>{c.title}</Link> ({c.media_count})</li>
    }
  )

  return (
    <div>
      <h1>Hello!</h1>
      {collections}
    </div>
  )
}
