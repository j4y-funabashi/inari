import { API } from "aws-amplify";

export interface timelineResponse {
	months: mediaMonth[];
}
export interface mediaMonth {
	id: string;
	title: string;
	type: string;
	media_count: number;
}

interface mediaSrc {
	small: string;
	medium: string;
	large: string;
}
export interface media {
	id: string;
	media_src: mediaSrc;
	date: string;
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
			id: "1",
			title: "Collection!",
			type: "",
			media_count: 1,
		},
		media: [
			{
				id: "123",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-25T10:00:00",
			},
			{
				id: "456",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-26T10:00:00",
			},
			{
				id: "789",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-28T10:00:00",
			},
			{
				id: "101112",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-28T11:00:00",
			},
			{
				id: "101113",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-28T11:00:00",
			},
			{
				id: "101114",
				media_src: {
					small: "https://via.placeholder.com/320",
					medium: "https://via.placeholder.com/320",
					large: "https://via.placeholder.com/1080",
				},
				date: "1984-01-28T11:00:00",
			},
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
		media: {
			id: "101114",
			media_src: {
				small: "https://via.placeholder.com/320",
				medium: "https://via.placeholder.com/320",
				large: "https://via.placeholder.com/1080",
			},
			date: "1984-01-28T13:38:00",
		},
	};

	return res as mediaDetailResponse;
};
