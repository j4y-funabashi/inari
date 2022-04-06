import {render, screen} from '@testing-library/react';
import React from 'react';
import MediaTimeline from './MediaTimeline';
import {BrowserRouter} from 'react-router-dom';

test('renders component', async () => {
  const mediaTimeline = {
    months: [
      {ID: "123", date: "2009-01", media_count: 10}
    ]
  }

  render(<BrowserRouter><MediaTimeline mediaTimeline={mediaTimeline} /></BrowserRouter>)
  screen.debug()
});
