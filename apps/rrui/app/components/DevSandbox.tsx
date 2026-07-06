import type { Collection } from "~/apiClient";
import { CollectionList } from "./collectionList";

export default function DevSandbox() {
  const collectionData: Collection[] = [
    {
      id: "123",
      title: "title",
      media_count: 11,
      type: "t",
      exported_count: 5
    }
  ]

  return (
    <div>
      <h1>CollectionList</h1>

      <CollectionList data={collectionData} />
    </div>
  )
}
