import { API } from "aws-amplify";
import { formatISO } from "date-fns";

export interface timelineResponse {
	months: mediaMonth[];
}
export interface mediaMonth {
	id: string;
	title: string;
	type: string;
	media_count: number;
}

export interface media {
	id: string;
	media_src: mediaSrc;
	date: string;
	location: location;
}

interface mediaSrc {
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
	collection_meta: mediaMonth;
	media: media[];
}
export interface mediaDetailResponse {
	media: media;
}

export type TimelineQuery = () => Promise<timelineResponse>;
export type TimelineMonthQuery = (
	monthID: string,
) => Promise<timelineMonthResponse>;
export type MediaDetailQuery = (
	mediaID: string,
) => Promise<mediaDetailResponse>;

export const fetchTimeline: TimelineQuery =
	async (): Promise<timelineResponse> => {
		const res = await API.get("photosAPIdev", "/months", {});
		console.log(res);

		return res as timelineResponse;
	};

export const mockFetchTimeline: TimelineQuery =
	async (): Promise<timelineResponse> => {
		const mockRes: timelineResponse = {
			months: [
				{
					id: "2020-01",
					title: "2020 Jan",
					type: "test-type",
					media_count: 1,
				},
			],
		};

		return mockRes;
	};

export const fetchTimelineMonth: TimelineMonthQuery = async (
	monthID: string,
): Promise<timelineMonthResponse> => {
	const res = await API.get("photosAPIdev", `/month/${monthID}`, {});
	console.log(res);

	return res as timelineMonthResponse;
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
	const res = await API.get("photosAPIdev", `/media/${mediaID}`, {});
	console.log(res);

	return res as mediaDetailResponse;
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
		media_src: {
			small: "https://via.placeholder.com/320",
			medium: "https://via.placeholder.com/320",
			large: "https://via.placeholder.com/1080",
		},
		date: formatISO(dat),
		location: mockLocation(),
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
