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
} from "../apiClient";

import MediaDetailPage from "./MediaDetailPage";
import MediaTimelineMonthPage from "./MediaTimelineMonthPage";
import MediaTimelinePage from "./MediaTimelinePage";
import NavBar from "./NavBar";

interface RouterProps {
	isDevMode: boolean;
}

const Router: React.FunctionComponent<RouterProps> = (props: RouterProps) => {
	const timelineQuery: TimelineQuery = props.isDevMode
		? mockFetchTimeline
		: fetchTimeline;
	const timelineMonthQuery: TimelineMonthQuery = props.isDevMode
		? mockFetchTimelineMonth
		: fetchTimelineMonth;
	const mediaDetailQuery: MediaDetailQuery = props.isDevMode
		? mockFetchMediaDetail
		: fetchMediaDetail;

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
					<MediaDetailPage fetchMediaDetail={mediaDetailQuery} />
				</Route>
			</Switch>

			<NavBar />
		</BrowserRouter>
	);
};

export default Router;
