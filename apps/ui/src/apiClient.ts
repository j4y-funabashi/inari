import { API } from 'aws-amplify';

export interface timelineResponse {
    months: mediaMonth[]
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
    media: media[]
}

export type TimelineQuery = () => Promise<timelineResponse>
export type TimelineMonthQuery = (monthID: string) => Promise<timelineMonthResponse>

export const fetchTimeline: TimelineQuery = async (): Promise<timelineResponse> => {
    const res = await API.get("photosAPIdev", "/months", {})
    console.log(res)

    return res as timelineResponse
}

export const fetchTimelineMonth: TimelineMonthQuery = async (monthID: string): Promise<timelineMonthResponse> => {
    const res = await API.get("photosAPIdev", "/month/" + monthID, {})
    console.log(res)

    return res as timelineMonthResponse
}
