import TestController from 'testcafe';

const baseURL = "https://jay.funabashi.co.uk/"

fixture(`happy path`).page(baseURL)

test("happy", async (t) => {
	await login(t)
})

const login = async (t: TestController) => {
	await t.typeText("#login-user-name", "chicken")
}
