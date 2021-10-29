import React from 'react';
import {Link} from 'react-router-dom';
import {timelineResponse} from '../apiClient';

export interface MediaTimelineProps {
  mediaTimeline: timelineResponse
}

const MediaTimeline: React.FunctionComponent<MediaTimelineProps> = (props: MediaTimelineProps) => {
  const {mediaTimeline} = props

  const mediaDays = mediaTimeline.days.map((mediaDay) => {
    const media = mediaDay.media.map((mediaItem) => {
      const srcString = "/thmnb/sqsm_"+mediaItem.media_src
      const linkString = "/media/"+mediaItem.id
      return (
        <Link to={linkString}><img src={srcString} alt="" /></Link>
      )
    })
    return (
      <div>
        <h1>{mediaDay.date}</h1>
        <div>
          {media}
        </div>
      </div>
    )
  })

  return (<div>{mediaDays}</div>)
}

export default MediaTimeline
