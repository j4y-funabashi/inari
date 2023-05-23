import axios from "axios";
import { formatISO } from "date-fns";

const httpClient = axios.create({
	baseURL: "/api/",
	timeout: 1000,
});

export type collectionsResponse = collection[]

export interface collection {
	id: string;
	title: string;
	type: string;
	media_count: number;
}

export interface media {
	id: string;
	thumbnails: thumbnails;
	media_metadata: MediaMetadata;
	date: string;
	location: location;
	caption: string;
	collections: collection[]
}

interface MediaMetadata {
	date: string;
}

interface thumbnails {
	small: string;
	medium: string;
	large: string;
}

interface location {
	region: string;
	locality: string;
	country: country;
	cooridnates: coordinates;
}
interface country {
	short: string;
	long: string;
}
interface coordinates {
	lat: number;
	lng: number;
}

export interface timelineMonthResponse {
	collection_meta: collection;
	media: media[];
}
export interface mediaDetailResponse {
	media: media;
}

export type TimelineQuery = () => Promise<collectionsResponse>;
export type TimelineMonthQuery = (
	monthID: string,
) => Promise<timelineMonthResponse>;
export type MediaDetailQuery = (
	mediaID: string,
) => Promise<mediaDetailResponse>;

export const fetchTimeline: TimelineQuery =
	async (): Promise<collectionsResponse> => {
		const { data: res } = await httpClient.get<collectionsResponse>("/timeline/months");
		console.log(res);

		return res;
	};

export const mockFetchTimeline: TimelineQuery =
	async (): Promise<collectionsResponse> => {
		const mockRes: collectionsResponse = [
			{
				id: "2020-01",
				title: "2020 Jan",
				type: "test-type",
				media_count: 1,
			},
		]

		return mockRes;
	};

export const fetchTimelineMonth: TimelineMonthQuery = async (
	monthID: string,
): Promise<timelineMonthResponse> => {
	const { data: res } = await httpClient.get<timelineMonthResponse>(
		`/timeline/month/${monthID}`,
	);
	console.log(res);

	return res;
};
export const mockFetchTimelineMonth: TimelineMonthQuery = async (
	monthID: string,
): Promise<timelineMonthResponse> => {
	const mRes: timelineMonthResponse = {
		collection_meta: {
			id: crypto.randomUUID(),
			title: "Collection!",
			type: "",
			media_count: 1,
		},
		media: [
			mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 25, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 25, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 25, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 2, 19, 0, 52)),
			mockMedia(new Date(1984, 0, 2, 19, 0, 52)),
		],
	};
	return mRes as timelineMonthResponse;
};

export const fetchMediaDetail: MediaDetailQuery = async (
	mediaID: string,
): Promise<mediaDetailResponse> => {
	const { data: res } = await httpClient.get<mediaDetailResponse>(
		`/media/${mediaID}`,
	);
	console.log(res);

	return res;
};

export const mockFetchMediaDetail: MediaDetailQuery = async (
	mediaID: string,
): Promise<mediaDetailResponse> => {
	const res = {
		media: mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
	};

	return res as mediaDetailResponse;
};

export const mockMedia = (dat: Date): media => {
	return {
		id: crypto.randomUUID(),
		thumbnails: {
			small: "https://via.placeholder.com/92",
			medium: "https://via.placeholder.com/420",
			large: "https://via.placeholder.com/1080",
		},
		media_metadata: {
			date: formatISO(dat),
		},
		date: formatISO(dat),
		location: mockLocation(),
		caption: "hello this is a good media!",
		collections: []
	};
};
const mockLocation = (): location => {
	return {
		region: "West Yorkshire",
		locality: "Leeds",
		country: {
			short: "GB",
			long: "United Kingdom",
		},
		cooridnates: {
			lat: 53.8700189722222,
			lng: -1.561703,
		},
	};
};
