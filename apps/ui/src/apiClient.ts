import {API} from 'aws-amplify';

export interface timelineResponse {
    days: mediaDay[]
}
interface mediaDay {
    date: string;
    media: mediaItem[];
}
interface mediaItem {
    id: string;
    mime_type: string;
    date: string;
    media_src: string;
}

export type TimelineQuery = () => Promise<timelineResponse>

export const fetchTimeline:TimelineQuery = async ():Promise<timelineResponse> => {
    const res = await API.get("photosAPIdev", "/timeline/", {})
    console.log(res)

    return res as timelineResponse
}
