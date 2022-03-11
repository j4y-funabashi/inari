import {API} from 'aws-amplify';

export interface timelineResponse {
    months: mediaMonth[]
}
interface mediaMonth {
    date: string;
    ID: string;
}

export type TimelineQuery = () => Promise<timelineResponse>

export const fetchTimeline:TimelineQuery = async ():Promise<timelineResponse> => {
    const res = await API.get("photosAPIdev", "/timeline/", {})
    console.log(res)

    return res as timelineResponse
}
