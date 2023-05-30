'use client';
import useSWR, { Fetcher } from 'swr';

interface Collection {
  id: string
  title: string
  media_count: number
  type: string
}

export default function Home() {

  return (
    <CollectionList />
  )
}

const collectionListFetcher: Fetcher<Collection[], string> = (type) => getCollectionsByType(type)

const CollectionList = function () {

  const { data, error, isLoading } = useSWR('/api/timeline/months', collectionListFetcher)

  if (error) return <div>failed to load</div>
  if (isLoading) return <div>loading...</div>

  console.log(data, error, isLoading)

  const collections = data?.map(
    (c) => {
      const collectionLink = "/collection/" + c.id
      return <li key={c.id}><a href={collectionLink}>{c.title}</a> ({c.media_count})</li>
    }
  )

  return (
    <div>
      <h1>Hello!</h1>
      {collections}
    </div>
  )
}

const getCollectionsByType = async function (type: string): Promise<Collection[]> {
  const res = await fetch("/api/timeline/months")
  console.log(res)

  return res.json()
}
