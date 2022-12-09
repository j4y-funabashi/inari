import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import {
	mockFetchTimeline,
	fetchTimeline,
	TimelineQuery,
	TimelineMonthQuery,
	mockFetchTimelineMonth,
	fetchTimelineMonth,
	MediaDetailQuery,
	mockFetchMediaDetail,
	fetchMediaDetail,
	mediaDetailResponse,
	mockMedia,
} from "../apiClient";

import MediaDetailPage from "./MediaDetailPage";
import MediaTimelineMonthPage from "./MediaTimelineMonthPage";
import MediaTimelinePage from "./MediaTimelinePage";
import NavBar from "./NavBar";

interface RouterProps {
	isDevMode: boolean;
}

const Router: React.FunctionComponent<React.PropsWithChildren<RouterProps>> = (
	props: RouterProps,
) => {
	// API calls
	const timelineQuery: TimelineQuery = props.isDevMode
		? mockFetchTimeline
		: fetchTimeline;
	const timelineMonthQuery: TimelineMonthQuery = props.isDevMode
		? mockFetchTimelineMonth
		: fetchTimelineMonth;
	const mediaDetailQuery: MediaDetailQuery = props.isDevMode
		? mockFetchMediaDetail
		: fetchMediaDetail;

	// timeline state
	const [mediaDetail, setMediaDetailData] = React.useState<mediaDetailResponse>(
		{
			media: mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
		},
	);

	return (
		<BrowserRouter>
			<Routes>
				<Route
					path="/"
					element={<MediaTimelinePage fetchTimeline={timelineQuery} />}
				/>
				<Route
					path="/time/month/:monthid"
					element={
						<MediaTimelineMonthPage fetchTimelineMonth={timelineMonthQuery} />
					}
				/>
				<Route
					path="/media/:mediaid"
					element={
						<MediaDetailPage
							fetchMediaDetail={mediaDetailQuery}
							media={mediaDetail}
							setMediaDetailData={setMediaDetailData}
						/>
					}
				/>
			</Routes>
			<NavBar />
		</BrowserRouter>
	);
};

export default Router;
