import { http, HttpResponse } from "msw";

export const getTemp2 = http.get("https://api.example.com/product", () => {
	return HttpResponse.json({
		id: "abc-123",
		firstName: "John",
		lastName: "Maverick",
	});
});

