import React from 'react';
import {fetchTimeline, timelineResponse} from '../apiClient';
import MediaTimeline from './MediaTimeline';

const MediaTimelinePage: React.FunctionComponent = () => {
  const [timelineData, setTimelineData] = React.useState<timelineResponse>({months: []});

  React.useEffect(() => {
    (async () => {
      const timelineResponse = await fetchTimeline()
      setTimelineData(timelineResponse)
    })()
  }, [setTimelineData])

  return (
    <MediaTimeline mediaTimeline={timelineData} />
  )
}

export default MediaTimelinePage
