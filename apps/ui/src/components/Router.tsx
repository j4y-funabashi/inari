import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";
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

const Router: React.FunctionComponent<RouterProps> = (props: RouterProps) => {
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

	// state
	const [mediaDetail, setMediaDetailData] = React.useState<mediaDetailResponse>(
		{
			media: mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
		},
	);

	return (
		<BrowserRouter>
			<Switch>
				<Route exact={true} path="/">
					<MediaTimelinePage fetchTimeline={timelineQuery} />
				</Route>
				<Route path="/time/month/:monthid">
					<MediaTimelineMonthPage fetchTimelineMonth={timelineMonthQuery} />
				</Route>
				<Route path="/media/:mediaid">
					<MediaDetailPage
						fetchMediaDetail={mediaDetailQuery}
						media={mediaDetail}
						setMediaDetailData={setMediaDetailData}
					/>
				</Route>
			</Switch>

			<NavBar />
		</BrowserRouter>
	);
};

export default Router;
