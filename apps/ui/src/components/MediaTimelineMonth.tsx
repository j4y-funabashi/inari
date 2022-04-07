import React from 'react';
import {timelineMonthResponse} from '../apiClient';
import {Link} from 'react-router-dom';

export interface MediaTimelineMonthProps {
  mediaTimeline: timelineMonthResponse
}

const MediaTimelineMonth: React.FunctionComponent<MediaTimelineMonthProps> = (props: MediaTimelineMonthProps) => {
  const {mediaTimeline} = props

  const header = (
    <div>
      <h1>{ mediaTimeline.collection_meta.date }</h1>
      <small>{ mediaTimeline.collection_meta.media_count }</small>
    </div>
  )
  const media = mediaTimeline.media.map((m) => {
    return (
      <li key={m.id}><Link to={`/media/${m.id}`}>{m.media_src}</Link> <small>({m.date})</small></li>
    )
  })

  return (
    <div>
      {header}
      <ol>{media}</ol>
    </div>
  )
}

export default MediaTimelineMonth
