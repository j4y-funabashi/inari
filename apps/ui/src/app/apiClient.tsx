import { Fetcher } from "swr";
import {
  getMockCollection,
  getMockCollectionDetail,
  getMockMedia,
} from "./__fixtures__";

export interface Media {
  id: string;
  thumbnails: Thumbnails;
  collections: Collection[];
  date: string;
  location?: Location;
  caption?: string;
}
interface Thumbnails {
  medium: string;
  large: string;
}

export interface Collection {
  id: string;
  title: string;
  media_count: number;
  type: CollectionType;
}

interface Location {
  country: Country;
  region: string;
  locality: string;
  coordinates?: Coordinates;
}
interface Country {
  short: string;
  long: string;
}
interface Coordinates {
  lat: number;
  lng: number;
}

export interface CollectionDetail {
  collection_meta: Collection;
  media: Media[];
}

export enum CollectionType {
  CollectionTypeInbox = "inbox",
  CollectionTypeCamera = "camera",
  CollectionTypeTimelineMonth = "timeline_month",
  CollectionTypeTimelineDay = "timeline_day",
  CollectionTypePlacesCountry = "places_country",
  CollectionTypePlacesRegion = "places_region",
  CollectionTypeHashTag = "hashtag",
}

const getCollectionsByType = async function (
  type: string,
): Promise<Collection[]> {
  const res = await fetch("/api/timeline/months");
  console.log(res);

  return res.json();
};
const mockGetCollectionsByType = function (type: string): Collection[] {
  return [
    getMockCollection("inbox Apr 2023"),
    getMockCollection("inbox Mar 2023"),
    getMockCollection("inbox Feb 2023"),
    getMockCollection("inbox Jan 2023"),
    getMockCollection("inbox Dec 2022"),
  ];
};

const getCollectionDetail = async function (
  id: string,
): Promise<CollectionDetail> {
  const res = await fetch("/api/timeline/month/" + id);
  console.log(res);

  return res.json();
};

const mockGetCollectionDetail = function (id: string): CollectionDetail {
  return getMockCollectionDetail();
};

export const deleteMedia = async function (id: string) {
  const requestOptions: RequestInit = {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
  };
  const res = await fetch("/api/media/" + id, requestOptions);
  console.log(res);
};

export const updateMediaCaption = async function (id: string, caption: string) {
  const requestOptions: RequestInit = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: caption,
  };
  const res = await fetch("/api/media/" + id + "/caption", requestOptions);
  console.log(res);
};

export const updateMediaHashtag = async function (id: string, hashtag: string) {
  const requestOptions: RequestInit = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: hashtag,
  };
  const res = await fetch("/api/media/" + id + "/tag", requestOptions);
  console.log(res);
};

export const NewCollectionLister = (
  env: string,
): Fetcher<Collection[], string> => {
  switch (env) {
    case "production":
      const collectionListFetcher: Fetcher<Collection[], string> = (type) =>
        getCollectionsByType(type);
      return collectionListFetcher;

    default:
      return mockGetCollectionsByType;
  }
};

export const NewFetchCollectionDetail = (
  env: string,
): Fetcher<CollectionDetail, string> => {
  switch (env) {
    case "production":
      const fetchCollectionDetail: Fetcher<CollectionDetail, string> = (id) =>
        getCollectionDetail(id);
      return fetchCollectionDetail;

    default:
      return mockGetCollectionDetail;
  }
};
