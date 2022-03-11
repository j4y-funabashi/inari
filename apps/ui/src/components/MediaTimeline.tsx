import React from 'react';
import {timelineResponse} from '../apiClient';

export interface MediaTimelineProps {
  mediaTimeline: timelineResponse
}

const MediaTimeline: React.FunctionComponent<MediaTimelineProps> = (props: MediaTimelineProps) => {
  const {mediaTimeline} = props

  const mediaMonths = mediaTimeline.months.map((m) => {
    return (
      <div>
        <h1>{m.date} ({m.media_count})</h1>
      </div>
    )
  })

  return (<div>{mediaMonths}</div>)
}

export default MediaTimeline
