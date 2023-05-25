import React from 'react';
import { TimelineQuery, collectionsResponse } from '../apiClient';
import MediaTimeline from './MediaTimeline';

interface MediaTimelinePageProps {
  fetchTimeline: TimelineQuery
}

const MediaTimelinePage: React.FunctionComponent<React.PropsWithChildren<MediaTimelinePageProps>> = (props: MediaTimelinePageProps) => {
  const [timelineData, setTimelineData] = React.useState<collectionsResponse>([]);

  const collection_type = "inbox"

  React.useEffect(() => {
    (async () => {
      const timelineResponse = await props.fetchTimeline()
      setTimelineData(timelineResponse)
    })()
  }, [setTimelineData, props])

  return (
    <MediaTimeline mediaTimeline={timelineData} collection_type={collection_type} />
  )
}

export default MediaTimelinePage
