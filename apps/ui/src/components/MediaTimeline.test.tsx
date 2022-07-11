import { render, screen } from '@testing-library/react';
import React from 'react';
import MediaTimeline from './MediaTimeline';
import { BrowserRouter } from 'react-router-dom';

test('renders component', async () => {
  const mediaTimeline = {
    months: [
      { id: "123", title: "2009-01", media_count: 10, type: "timeline_month" }
    ]
  }

  render(<BrowserRouter><MediaTimeline mediaTimeline={mediaTimeline} /></BrowserRouter>)
  screen.debug()
});
