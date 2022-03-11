import React from 'react';
import {timelineResponse} from '../apiClient';

export interface MediaTimelineProps {
  mediaTimeline: timelineResponse
}

const MediaTimeline: React.FunctionComponent<MediaTimelineProps> = (props: MediaTimelineProps) => {
  const {mediaTimeline} = props

  const mediaMonths = mediaTimeline.months.map((m) => {
    return (
        <li>{m.date} <small>({m.media_count})</small></li>
    )
  })

  return (<ol>{mediaMonths}</ol>)
}

export default MediaTimeline
