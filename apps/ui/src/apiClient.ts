import {API} from 'aws-amplify';

export interface timelineResponse {
    months: mediaMonth[]
}
interface mediaMonth {
    date: string;
    ID: string;
    media_count: number;
}
interface media {
    id: string;
    media_src: string;
}
export interface timelineMonthResponse {
    collection_meta: mediaMonth;
    media: media[]
}

export type TimelineQuery = () => Promise<timelineResponse>
export type TimelineMonthQuery = (monthID:string) => Promise<timelineMonthResponse>

export const fetchTimeline:TimelineQuery = async ():Promise<timelineResponse> => {
    const res = await API.get("photosAPIdev", "/time/", {})
    console.log(res)

    return res as timelineResponse
}

export const fetchTimelineMonth:TimelineMonthQuery = async (monthID:string):Promise<timelineMonthResponse> => {
    const res = await API.get("photosAPIdev", "/time/month" + monthID, {})
    console.log(res)

    return res as timelineMonthResponse
}
