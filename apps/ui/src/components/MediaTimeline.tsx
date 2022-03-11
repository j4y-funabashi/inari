import React from 'react';
import {Link} from 'react-router-dom';
import {timelineResponse} from '../apiClient';

export interface MediaTimelineProps {
  mediaTimeline: timelineResponse
}

const MediaTimeline: React.FunctionComponent<MediaTimelineProps> = (props: MediaTimelineProps) => {
  const {mediaTimeline} = props

  const mediaMonths = mediaTimeline.months.map((mediaDay) => {
    return (
      <div>
        <h1>{mediaDay.date}</h1>
      </div>
    )
  })

  return (<div>{mediaMonths}</div>)
}

export default MediaTimeline
