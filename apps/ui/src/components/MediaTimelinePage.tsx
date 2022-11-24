import React from 'react';
import { TimelineQuery, timelineResponse } from '../apiClient';
import MediaTimeline from './MediaTimeline';

interface MediaTimelinePageProps {
  fetchTimeline: TimelineQuery
}

const MediaTimelinePage: React.FunctionComponent<MediaTimelinePageProps> = (props: MediaTimelinePageProps) => {
  const [timelineData, setTimelineData] = React.useState<timelineResponse>({ months: [] });

  React.useEffect(() => {
    (async () => {
      const timelineResponse = await props.fetchTimeline()
      setTimelineData(timelineResponse)
    })()
  }, [setTimelineData, props])

  return (
    <MediaTimeline mediaTimeline={timelineData} />
  )
}

export default MediaTimelinePage
