import { http, HttpResponse } from "msw";

export const getTemp1 = http.get("http://localhost/api/timeline/months", () => {
	return HttpResponse.json(
		[
			{
				id: "abc-123",
				title: "John",
				media_count: 1,
				type: "",
				exported_count: 1
			}
		]
	);
});

