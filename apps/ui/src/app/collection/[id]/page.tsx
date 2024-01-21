"use client";

import { CollectionDetail, NewFetchCollectionDetail } from "@/app/apiClient";
import { MediaGallery } from "@/app/components/MediaGallery";
import useSWR from "swr";

interface CollectionDetailParams {
  params: {
    id: string;
  };
}

export default function CollectionDetailPage({
  params,
}: CollectionDetailParams) {
  const collectionID = params.id;
  const collectionDetailFetcher = NewFetchCollectionDetail(
    process.env.NODE_ENV,
  );

  const {
    data: collectionDetailData,
    error: collectionDetailError,
    isLoading: collectionDetailLoading,
  } = useSWR<CollectionDetail>(collectionID, collectionDetailFetcher);
  if (!collectionDetailData) {
    return;
  }

  return (
    <div>
      <MediaGallery data={collectionDetailData} />
    </div>
  );
}
