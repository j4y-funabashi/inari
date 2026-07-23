import MediaGallery from "~/components/mediaGallery";
import type { Route } from "./+types/collection";
import { type CollectionDetail, type Collection, mockListCollections, mockGetCollectionDetail, getCollectionDetail, listCollections } from "~/apiClient";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Inari" },
    { name: "description", content: "" },
  ];
}

interface MediaGalleryResponse {
  collectionDetail: CollectionDetail
  collections: Collection[]
}

export async function clientLoader({ params }: Route.ClientLoaderArgs): Promise<MediaGalleryResponse> {
  if (process.env.NODE_ENV == "development") {
    return {
      collectionDetail: await mockGetCollectionDetail(params.collectionid),
      collections: await mockListCollections("hashtag")
    }
  }
  return {
    collectionDetail: await getCollectionDetail(params.collectionid),
    collections: await listCollections("hashtag")
  }
}

export default function CollectionDetailPage({ loaderData }: Route.ComponentProps) {
  return <MediaGallery collections={loaderData.collections} collectionDetail={loaderData.collectionDetail} />;
}
