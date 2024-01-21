import {
  Collection,
  CollectionDetail,
  CollectionType,
  Media,
} from "../apiClient";
import { randomUUID } from "crypto";

export const getMockCollection = (title: string): Collection => {
  const uuid = randomUUID();
  return {
    id: `test-1-${uuid}`,
    title: `${title}`,
    media_count: 5,
    type: CollectionType.CollectionTypeInbox,
  };
};

export const getMockMedia = (): Media => {
  const urlPrefix = "https://picsum.photos";
  const uuid = randomUUID();
  return {
    id: `testid-${uuid}`,
    thumbnails: {
      medium: `${urlPrefix}/420/420`,
      large: `${urlPrefix}/1080/600`,
    },
    date: "2022-01-28T10:01:02Z",
    collections: [
      getMockCollection("inbox Jan 2022"),
      getMockCollection("January 2022"),
      getMockCollection("West Yorkshire, United Kingdom"),
    ],
    location: {
      country: { long: "United Kingdom", short: "c" },
      region: "West Yorkshire",
      locality: "Meanwood",
    },
    caption: "This is the caption",
  };
};

export const getMockCollectionDetail = (): CollectionDetail => {
  return {
    collection_meta: getMockCollection("inbox Jan 2023"),
    media: getMockMediaList(3),
  };
};
const getMockMediaList = (count: number): Media[] => {
  const out: Media[] = [];

  for (let index = 0; index < count; index++) {
    out.push(getMockMedia());
  }

  return out;
};
