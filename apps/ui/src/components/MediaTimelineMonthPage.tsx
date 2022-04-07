import React from 'react';
import {fetchTimelineMonth, timelineMonthResponse} from '../apiClient';

const MediaTimelineMonthPage: React.FunctionComponent = () => {
  const [timelineData, setTimelineData] = React.useState<timelineMonthResponse>({media: [],collection_meta: {date:"", ID:"", media_count:0}});

  const monthID = "2018-05"

  React.useEffect(() => {
    (async () => {
      const timelineResponse = await fetchTimelineMonth(monthID)
      console.log(timelineResponse)
      setTimelineData(timelineResponse)
    })()
  }, [setTimelineData])

  console.log(timelineData)

  return (<div></div>)
}

export default MediaTimelineMonthPage
