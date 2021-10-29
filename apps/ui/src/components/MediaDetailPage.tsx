import React from 'react';
import {useParams} from 'react-router-dom';

type MediaDetailURLParams = {
  mediaid: string
}

const MediaDetailPage: React.FunctionComponent = () => {
  // const {mediaid} = useParams<MediaDetailURLParams>()

  return (

    <article className="vh-100 dt w-100 bg-dark-gray">
      <div className="dtc v-mid tc">
        <img src="https://photos-dev.funabashi.co.uk/thmnb/lg_20211016_143550_5deb3260c820dc1adc1b29282ad4d3d6.JPG" alt="" />
      </div>
    </article>

  )
}

export default MediaDetailPage
