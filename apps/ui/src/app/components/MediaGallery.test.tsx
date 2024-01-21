import { cleanup, render, screen } from "@testing-library/react";
import { AddHashtagForm } from "./AddHashtagForm";
import { userEvent } from "@testing-library/user-event";
import { MediaGallery } from "./MediaGallery";
import { CollectionDetail } from "../apiClient";
import { getMockCollectionDetail } from "../__fixtures__";

describe("AddHashtagForm", () => {
  afterEach(() => {
    cleanup();
  });

  it("works", async () => {
    // arrange
    const collectionDetail = getMockCollectionDetail();
    render(<MediaGallery data={collectionDetail} />);

    // act
    const heading = screen.getByRole("heading", {
      name: collectionDetail.collection_meta.title,
    });

    const media = screen.getByRole("imaage", {});

    // assert
    expect(heading).toBeVisible();
  });
});
