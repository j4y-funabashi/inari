import { render, screen } from '@testing-library/react';
import React from 'react';
import MediaTimelineMonth from './MediaTimelineMonth';
import { media } from '../apiClient'
import { BrowserRouter } from 'react-router-dom';

test('renders component', async () => {
  const testMedia1: media = {
    id: "test-id-1",
    media_src: {
      small: "img-sm.jpg",
      medium: "img-sm.jpg",
      large: "img-sm.jpg",
    },
    date: "2022-06-09T22:19:14Z",
  }
  const testMedia2: media = {
    id: "test-id-2",
    media_src: {
      small: "img-sm.jpg",
      medium: "img-sm.jpg",
      large: "img-sm.jpg",
    },
    date: "2022-06-09T23:19:14Z",
  }
  const testMedia3: media = {
    id: "test-id-3",
    media_src: {
      small: "img-sm.jpg",
      medium: "img-sm.jpg",
      large: "img-sm.jpg",
    },
    date: "2022-06-19T23:19:14Z",
  }

  const mediaTimeline = {
    collection_meta: {
      id: "123",
      title: "Jan 2006",
      type: "timeline_month",
      media_count: 200
    },
    media: [testMedia2, testMedia1, testMedia3]
  }

  render(<BrowserRouter><MediaTimelineMonth mediaTimeline={mediaTimeline} /></BrowserRouter>)
  screen.debug()
});

