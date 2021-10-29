import {render, screen} from '@testing-library/react';
import React from 'react';
import MediaTimeline from './MediaTimeline';

test('renders component', async () => {
  const mediaTimeline = {
    days: []
  }

  render(<MediaTimeline mediaTimeline={mediaTimeline} />)
  screen.debug()
});
