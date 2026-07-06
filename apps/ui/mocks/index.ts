async function initMocks() {
	if (typeof window === "undefined") {
		console.log("init mocks - server")
		const { server } = await import("./server");
		server.listen();
	} else {
		console.log("init mocks - browser")
		const { worker } = await import("./browser");
		worker.start();
	}
}

initMocks();

export { };
