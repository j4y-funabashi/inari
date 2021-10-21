import {API} from 'aws-amplify';

const fetchTimeline = async () => {
    const res = await API.get("photosAPIdev", "/timeline/", {})
    console.log(res)
}
