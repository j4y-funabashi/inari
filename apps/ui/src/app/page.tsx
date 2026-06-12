import { Suspense } from "react";
import { CollectionList } from "./components/collectionList";
import Loading from "./components/Loading";
import { useCollections } from "./apiClient";

export default function Home() {

  const data = useCollections("inbox")

  return (
    <Suspense fallback={<Loading />}>
      <CollectionList data={data} />
    </Suspense>
  )
}

