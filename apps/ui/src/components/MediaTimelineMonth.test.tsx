import { render, screen } from '@testing-library/react';
import React from 'react';
import MediaTimelineMonth from './MediaTimelineMonth';
import { media } from '../apiClient'
import { BrowserRouter } from 'react-router-dom';

test('renders component', async () => {
  const testMedia: media = {
    id: "test-id-1",
    media_src: {
      small: "img-sm.jpg",
      medium: "img-sm.jpg",
      large: "img-sm.jpg",
    },
    date: "media-dt-123",
  }

  const mediaTimeline = {
    collection_meta: {
      id: "123",
      title: "Jan 2006",
      type: "timeline_month",
      media_count: 200
    },
    media: [testMedia]
  }

  render(<BrowserRouter><MediaTimelineMonth mediaTimeline={mediaTimeline} /></BrowserRouter>)
  screen.debug()
});

