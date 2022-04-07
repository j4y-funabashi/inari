import React from 'react';
import {fetchTimelineMonth, timelineMonthResponse} from '../apiClient';
import {useParams} from 'react-router-dom';
import MediaTimelineMonth from './MediaTimelineMonth'

type urlParams = {
  monthid: string
}

const MediaTimelineMonthPage: React.FunctionComponent = () => {
  const [timelineData, setTimelineData] = React.useState<timelineMonthResponse>({media: [],collection_meta: {date:"", ID:"", media_count:0}});

  const {monthid} = useParams<urlParams>()
  console.log(monthid)

  React.useEffect(() => {
    (async () => {
      const timelineResponse = await fetchTimelineMonth(monthid)
      setTimelineData(timelineResponse)
    })()
  }, [setTimelineData])

  console.log(timelineData)

  return (<MediaTimelineMonth mediaTimeline={timelineData} />)
}

export default MediaTimelineMonthPage
