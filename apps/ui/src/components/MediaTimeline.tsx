import React from 'react';
import {timelineResponse} from '../apiClient';
import {Link} from 'react-router-dom';

export interface MediaTimelineProps {
  mediaTimeline: timelineResponse
}

const MediaTimeline: React.FunctionComponent<MediaTimelineProps> = (props: MediaTimelineProps) => {
  const {mediaTimeline} = props

  const mediaMonths = mediaTimeline.months.map((m) => {
    return (
        <li key={m.ID}><Link to={`/time/month/${m.ID}`}>{m.date}</Link> <small>({m.media_count})</small></li>
    )
  })

  return (<ol>{mediaMonths}</ol>)
}

export default MediaTimeline
