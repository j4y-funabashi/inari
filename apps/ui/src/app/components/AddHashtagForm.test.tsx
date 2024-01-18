import { cleanup, render, screen } from "@testing-library/react";
import { AddHashtagForm } from "./AddHashtagForm";
import { userEvent } from "@testing-library/user-event";

describe("AddHashtagForm", () => {
  const handleSaveHashtag = jest.fn();

  afterEach(() => {
    cleanup();
  });

  it("removes special chars", async () => {
    // arrange
    const user = userEvent.setup();
    render(<AddHashtagForm handleSaveHashtag={handleSaveHashtag} />);

    // act
    const textbox = screen.getByRole("textbox", { name: "Add Hashtag" });
    await user.type(textbox, "h@e# ()l#l");

    // assert
    expect(textbox).toHaveValue("hell");
  });

  it("should call saveHashtag on click", async () => {
    // arrange
    const user = userEvent.setup();
    render(<AddHashtagForm handleSaveHashtag={handleSaveHashtag} />);

    // act
    const textbox = screen.getByRole("textbox", { name: "Add Hashtag" });
    const button = screen.getByRole("button", { name: "Save" });
    await user.type(textbox, "h@e# ()l#l");
    await user.click(button);

    // assert
    expect(handleSaveHashtag).toBeCalledWith("hell");
  });

  it("should not call saveHashtag when textbox is empty", async () => {
    // arrange
    const user = userEvent.setup();
    render(<AddHashtagForm handleSaveHashtag={handleSaveHashtag} />);

    // act
    const button = screen.getByRole("button", { name: "Save" });
    const textbox = screen.getByRole("textbox", { name: "Add Hashtag" });
    await user.click(button);

    // assert
    expect(textbox).toHaveValue("");
    expect(handleSaveHashtag).toBeCalledTimes(0);
  });
});
